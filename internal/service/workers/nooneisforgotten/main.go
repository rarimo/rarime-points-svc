package nooneisforgotten

import (
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/referralid"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func Run(cfg config.Config, sig chan struct{}) {
	db := cfg.DB().Clone()

	evType := cfg.EventTypes().Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType != nil {
		if err := updatePassportScanEvents(db); err != nil {
			panic(err)
		}

		if evType.AutoClaim {
			err := claimPassportScanEvents(cfg)
			if err != nil {
				panic(err)
			}
		}
	}

	evType = cfg.EventTypes().Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evType != nil {
		if err := updateReferralUserEvents(db); err != nil {
			panic(err)
		}
	}
	sig <- struct{}{}
}

func updatePassportScanEvents(db *pgdb.DB) error {
	balances, err := pg.NewBalances(db).WithoutPassportEvent()
	if err != nil {
		return fmt.Errorf("failed to select balances without points for passport scan: %w", err)
	}

	toUpdate := make([]string, 0, len(balances))
	for _, balance := range balances {
		if balance.EventID != nil {
			toUpdate = append(toUpdate, *balance.EventID)
			continue
		}
	}

	if len(toUpdate) != 0 {
		_, err = pg.NewEvents(db).
			FilterByID(toUpdate...).
			Update(data.EventFulfilled, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to update passport scan events: %w", err)
		}
	}

	return nil
}

func claimPassportScanEvents(cfg config.Config) error {
	evType := cfg.EventTypes().Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType == nil {
		return nil
	}

	db := cfg.DB().Clone()
	events, err := pg.NewEvents(db).FilterByType(evtypes.TypePassportScan).FilterByStatus(data.EventFulfilled).Select()
	if err != nil {
		return fmt.Errorf("failed to select passport scan events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	eventsMap := make(map[string]data.Event, len(events))
	nullifiers := make([]string, len(events))
	for i, event := range events {
		nullifiers[i] = event.Nullifier
		eventsMap[event.Nullifier] = event
	}

	balances, err := pg.NewBalances(db).FilterByNullifier(nullifiers...).FilterDisabled().Select()
	if err != nil {
		return fmt.Errorf("failed to select balances for claim passport scan event: %w", err)
	}

	if len(balances) == 0 {
		return nil
	}

	countryCodesMap := make(map[string]int64, len(balances))
	for _, balance := range balances {
		// normally should never happen
		if balance.Country == nil {
			return fmt.Errorf("balance have fulfilled passport scan event, but have no country")
		}
		if _, ok := countryCodesMap[*balance.Country]; !ok {
			countryCodesMap[*balance.Country] = 0
		}
	}

	countryCodesSlice := make([]string, 0, len(countryCodesMap))
	for k := range countryCodesMap {
		countryCodesSlice = append(countryCodesSlice, k)
	}

	countries, err := pg.NewCountries(db).FilterByCodes(countryCodesSlice...).Select()
	if err != nil {
		return fmt.Errorf("failed to select countries for claim passport scan events: %w", err)
	}

	for _, country := range countries {
		if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
			delete(countryCodesMap, country.Code)
			continue
		}
		countryCodesMap[country.Code] = country.ReserveLimit - country.Reserved
	}

	for _, balance := range balances {
		// country should exists because of previous validation
		if _, ok := countryCodesMap[*balance.Country]; !ok {
			continue
		}

		countryCodesMap[*balance.Country] -= evType.Reward
		if countryCodesMap[*balance.Country] <= 0 {
			delete(countryCodesMap, *balance.Country)
		}

		_, err = claimEventWithPoints(cfg, eventsMap[balance.Nullifier], evType.Reward, &balance)
		if err != nil {
			return fmt.Errorf("failed to claim passport scan event: %w", err)
		}
	}

	return nil
}

func updateReferralUserEvents(db *pgdb.DB) error {
	refPairs, err := pg.NewBalances(db).WithoutReferralEvent()
	if err != nil {
		return fmt.Errorf("failed to select balances without points for referred users: %w", err)
	}

	toInsert := make([]data.Event, 0, len(refPairs))
	for _, ref := range refPairs {
		toInsert = append(toInsert, data.Event{
			Nullifier: ref.Referrer,
			Type:      evtypes.TypeReferralSpecific,
			Status:    data.EventFulfilled,
			Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, ref.Referred)),
		})
	}

	if len(toInsert) != 0 {
		err = pg.NewEvents(db).Insert(toInsert...)
		if err != nil {
			return fmt.Errorf("failed to insert referred user events: %w", err)
		}
	}

	return nil
}

func claimReferralSpecificEvents(cfg config.Config) error {
	evType := cfg.EventTypes().Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evType == nil {
		return nil
	}

	db := cfg.DB().Clone()
	events, err := pg.NewEvents(db).FilterByType(evtypes.TypeReferralSpecific).FilterByStatus(data.EventFulfilled).Select()
	if err != nil {
		return fmt.Errorf("failed to select passport scan events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	eventsMap := make(map[string][]data.Event, len(events))
	nullifiers := make([]string, len(events))
	for i, event := range events {
		nullifiers[i] = event.Nullifier
		eventsMap[event.Nullifier] = append(eventsMap[event.Nullifier], event)
	}

	balances, err := pg.NewBalances(db).FilterByNullifier(nullifiers...).FilterDisabled().Select()
	if err != nil {
		return fmt.Errorf("failed to select balances for claim passport scan event: %w", err)
	}

	if len(balances) == 0 {
		return nil
	}

	countryCodesMap := make(map[string]int64, len(balances))
	for _, balance := range balances {
		// normally should never happen
		if balance.Country == nil {
			return fmt.Errorf("balance have fulfilled passport scan event, but have no country")
		}
		if _, ok := countryCodesMap[*balance.Country]; !ok {
			countryCodesMap[*balance.Country] = 0
		}
	}

	countryCodesSlice := make([]string, 0, len(countryCodesMap))
	for k := range countryCodesMap {
		countryCodesSlice = append(countryCodesSlice, k)
	}

	countries, err := pg.NewCountries(db).FilterByCodes(countryCodesSlice...).Select()
	if err != nil {
		return fmt.Errorf("failed to select countries for claim passport scan events: %w", err)
	}

	for _, country := range countries {
		if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
			delete(countryCodesMap, country.Code)
			continue
		}
		countryCodesMap[country.Code] = country.ReserveLimit - country.Reserved
	}

	for _, balance := range balances {
		// country should exists because of previous validation
		if _, ok := countryCodesMap[*balance.Country]; !ok {
			continue
		}

		countryCodesMap[*balance.Country] -= evType.Reward
		if countryCodesMap[*balance.Country] <= 0 {
			delete(countryCodesMap, *balance.Country)
		}

		for _, event := range eventsMap[balance.Nullifier] {
			_, err = claimEventWithPoints(cfg, event, evType.Reward, &balance)
			if err != nil {
				return fmt.Errorf("failed to claim passport scan event: %w", err)
			}

		}
	}

	return nil
}

// claimEventWithPoints requires event to exist
func claimEventWithPoints(cfg config.Config, event data.Event, reward int64, balance *data.Balance) (claimed *data.Event, err error) {
	err = pg.NewEvents(cfg.DB().Clone()).Transaction(func() error {
		db := cfg.DB().Clone()
		// Upgrade level logic when threshold is reached
		refsCount, level := cfg.Levels().LvlUp(balance.Level, reward+balance.Amount)
		if level != balance.Level {
			count, err := pg.NewReferrals(db).FilterByNullifier(balance.Nullifier).Count()
			if err != nil {
				return fmt.Errorf("failed to get referral count: %w", err)
			}

			refToAdd := prepareReferralsToAdd(balance.Nullifier, uint64(refsCount), count)
			if err = pg.NewReferrals(db).Insert(refToAdd...); err != nil {
				return fmt.Errorf("failed to insert referrals: %w", err)
			}

			err = pg.NewBalances(db).FilterByNullifier(balance.Nullifier).Update(map[string]any{
				data.ColLevel: level,
			})
			if err != nil {
				return fmt.Errorf("failed to update level: %w", err)
			}
		}

		claimed, err = pg.NewEvents(db).FilterByID(event.ID).Update(data.EventClaimed, nil, &reward)
		if err != nil {
			return fmt.Errorf("update event status: %w", err)
		}

		err = pg.NewBalances(db).FilterByNullifier(balance.Nullifier).Update(map[string]any{
			data.ColAmount: pg.AddToValue(data.ColAmount, reward),
		})
		if err != nil {
			return fmt.Errorf("update balance amount: %w", err)
		}

		err = pg.NewCountries(db).FilterByCodes(*balance.Country).Update(map[string]any{
			data.ColReserved: pg.AddToValue(data.ColReserved, reward),
		})
		if err != nil {
			return fmt.Errorf("increase country reserve: %w", err)
		}

		return nil
	})

	return claimed, nil
}

func prepareReferralsToAdd(nullifier string, count, index uint64) []data.Referral {
	refCodes := referralid.NewMany(nullifier, count, index)
	refs := make([]data.Referral, len(refCodes))

	for i, code := range refCodes {
		refs[i] = data.Referral{
			ID:        code,
			Nullifier: nullifier,
			UsageLeft: 1,
		}
	}

	return refs
}
