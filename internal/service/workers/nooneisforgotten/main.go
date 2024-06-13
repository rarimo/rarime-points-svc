package nooneisforgotten

import (
	"fmt"
	"math"
	"sort"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/handlers"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func Run(cfg config.Config, sig chan struct{}) {
	db := cfg.DB().Clone()
	if err := pg.NewEvents(db).Transaction(func() error {
		return updatePassportScanEvents(db, cfg.EventTypes(), cfg.Levels())
	}); err != nil {
		panic(fmt.Errorf("failed to update passport scan events: %w", err))
	}

	if err := pg.NewEvents(db).Transaction(func() error {
		return updateReferralUserEvents(db, cfg.EventTypes())
	}); err != nil {
		panic(fmt.Errorf("failed to update referral user events"))
	}

	if err := pg.NewEvents(db).Transaction(func() error {
		return claimReferralSpecificEvents(db, cfg.EventTypes(), cfg.Levels())
	}); err != nil {
		panic(fmt.Errorf("failed to claim referral specific events"))
	}

	sig <- struct{}{}
}

func updatePassportScanEvents(db *pgdb.DB, types evtypes.Types, levels config.Levels) error {
	evType := types.Get(evtypes.TypePassportScan)
	if evType == nil {
		return nil
	}

	if !evType.AutoClaim && evtypes.FilterInactive(*evType) {
		return nil
	}

	balances, err := pg.NewBalances(db).WithoutPassportEvent()
	if err != nil {
		return fmt.Errorf("failed to select balances without points for passport scan: %w", err)
	}

	toFulfill := make([]string, 0, len(balances))
	countriesBalancesMap := make(map[string][]data.WithoutPassportEventBalance, len(balances))
	countriesList := make([]string, 0, len(balances))
	for _, balance := range balances {
		if balance.EventStatus == data.EventOpen {
			toFulfill = append(toFulfill, balance.EventID)
		}

		// country must exist because of db query logic
		if _, ok := countriesBalancesMap[*balance.Country]; !ok {
			countriesList = append(countriesList, *balance.Country)
			countriesBalancesMap[*balance.Country] = make([]data.WithoutPassportEventBalance, 0, len(balances))
		}
		countriesBalancesMap[*balance.Country] = append(countriesBalancesMap[*balance.Country], balance)
	}

	if !evType.AutoClaim {
		if len(toFulfill) != 0 {
			_, err = pg.NewEvents(db).
				FilterByID(toFulfill...).
				Update(data.EventFulfilled, nil, nil)
			if err != nil {
				return fmt.Errorf("failed to update passport scan events: %w", err)
			}
		}

		return nil
	}

	countries, err := pg.NewCountries(db).FilterByCodes(countriesList...).Select()
	if err != nil {
		return fmt.Errorf("failed to select countries: %w", err)
	}

	// we need sort, because firstly claim already fulfilled event
	// and then open events
	for _, country := range countries {
		if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
			continue
		}

		sort.SliceStable(countriesBalancesMap[country.Code], func(i, j int) bool {
			if countriesBalancesMap[country.Code][i].EventStatus == countriesBalancesMap[country.Code][j].EventStatus {
				return false
			}
			if countriesBalancesMap[country.Code][i].EventStatus == data.EventOpen {
				return true
			}
			return false
		})

		countToClaim := int(math.Min(
			float64(len(countriesBalancesMap[country.Code])),
			math.Ceil(float64(country.ReserveLimit-country.Reserved)/float64(evType.Reward))))

		for i := 0; i < countToClaim; i++ {
			// if event is inactive we claim only fulfilled events
			if countriesBalancesMap[country.Code][i].EventStatus == data.EventOpen && evtypes.FilterInactive(*evType) {
				break
			}

			eventID := countriesBalancesMap[country.Code][i].EventID
			_, err = pg.NewEvents(db).FilterByID(eventID).Update(data.EventClaimed, nil, &evType.Reward)
			if err != nil {
				return fmt.Errorf("update event status: %w", err)
			}

			err = handlers.DoClaimEventUpdates(
				levels,
				pg.NewReferrals(db),
				pg.NewBalances(db),
				pg.NewCountries(db),
				countriesBalancesMap[country.Code][i].Balance,
				evType.Reward)
			if err != nil {
				return fmt.Errorf("failed to do claim event updates for passport scan: %w", err)
			}

			countriesBalancesMap[country.Code][i].EventStatus = data.EventClaimed
		}
	}

	if evtypes.FilterInactive(*evType) {
		return nil
	}

	toFulfill = make([]string, 0, len(balances))
	for _, balances := range countriesBalancesMap {
		for _, balance := range balances {
			if balance.EventStatus == data.EventOpen {
				toFulfill = append(toFulfill, balance.EventID)
			}
		}
	}

	if len(toFulfill) == 0 {
		return nil
	}

	_, err = pg.NewEvents(db).
		FilterByID(toFulfill...).
		Update(data.EventFulfilled, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to update passport scan events: %w", err)
	}

	return nil
}

func updateReferralUserEvents(db *pgdb.DB, types evtypes.Types) error {
	evTypeRef := types.Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evTypeRef == nil {
		return nil
	}

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

	if len(toInsert) == 0 {
		return nil
	}

	if err = pg.NewEvents(db).Insert(toInsert...); err != nil {
		return fmt.Errorf("failed to insert referred user events: %w", err)
	}

	return nil
}

func claimReferralSpecificEvents(db *pgdb.DB, types evtypes.Types, levels config.Levels) error {
	evType := types.Get(evtypes.TypeReferralSpecific)
	if evType == nil {
		return nil
	}
	if !evType.AutoClaim {
		return nil
	}

	events, err := pg.NewEvents(db).FilterByType(evtypes.TypeReferralSpecific).FilterByStatus(data.EventFulfilled).Select()
	if err != nil {
		return fmt.Errorf("failed to select passport scan events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	nullifiersEventsMap := make(map[string][]data.Event, len(events))
	nullifiers := make([]string, 0, len(events))
	for _, event := range events {
		if _, ok := nullifiersEventsMap[event.Nullifier]; !ok {
			nullifiersEventsMap[event.Nullifier] = make([]data.Event, 0, len(events))
			nullifiers = append(nullifiers, event.Nullifier)
		}
		nullifiersEventsMap[event.Nullifier] = append(nullifiersEventsMap[event.Nullifier], event)
	}

	balances, err := pg.NewBalances(db).FilterByNullifier(nullifiers...).FilterDisabled().Select()
	if err != nil {
		return fmt.Errorf("failed to select balances for claim passport scan event: %w", err)
	}
	if len(balances) == 0 {
		return fmt.Errorf("critical: events present, but no balances with nullifier")
	}

	countriesBalancesMap := make(map[string][]data.Balance, len(balances))
	for _, balance := range balances {
		// country can't be nil because of db query logic
		if _, ok := countriesBalancesMap[*balance.Country]; !ok {
			countriesBalancesMap[*balance.Country] = make([]data.Balance, 0, len(balances))
		}

		countriesBalancesMap[*balance.Country] = append(countriesBalancesMap[*balance.Country], balance)
	}

	countryCodes := make([]string, 0, len(countriesBalancesMap))
	for k := range countriesBalancesMap {
		countryCodes = append(countryCodes, k)
	}

	countries, err := pg.NewCountries(db).FilterByCodes(countryCodes...).Select()
	if err != nil {
		return fmt.Errorf("failed to select countries for claim passport scan events: %w", err)
	}

	toClaim := make([]string, 0, len(events))
	for _, country := range countries {
		if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
			continue
		}

		limit := country.ReserveLimit - country.Reserved
		for _, balance := range countriesBalancesMap[country.Code] {
			if limit <= 0 {
				break
			}

			toAccrue := int64(0)
			for _, event := range nullifiersEventsMap[balance.Nullifier] {
				limit -= evType.Reward
				toClaim = append(toClaim, event.ID)
				toAccrue += evType.Reward
			}

			err = handlers.DoClaimEventUpdates(
				levels,
				pg.NewReferrals(db),
				pg.NewBalances(db),
				pg.NewCountries(db),
				balance,
				evType.Reward)
			if err != nil {
				return fmt.Errorf("failed to do claim event updates for referral specific event: %w", err)
			}
		}
	}

	_, err = pg.NewEvents(db).FilterByID(toClaim...).Update(data.EventClaimed, nil, &evType.Reward)
	if err != nil {
		return fmt.Errorf("update event status: %w", err)
	}

	return nil
}
