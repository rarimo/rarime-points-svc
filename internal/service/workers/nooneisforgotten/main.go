package nooneisforgotten

import (
	"errors"
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
		panic(fmt.Errorf("failed to claim referral specific events: %w", err))
	}

	sig <- struct{}{}
}

// updatePassportScanEvents is needed so that if the passport
// scan events were not fulfilled or claimed because the event was disabled,
// expired or no autoclaimed, fulfill and, if possible, claim them.
// First, there is an attempt to claim as many events as
// possible and to fulfill the rest of the events.
//
// Event will not be claimed if AutoClaim is disabled.
func updatePassportScanEvents(db *pgdb.DB, types evtypes.Types, levels config.Levels) error {
	evType := types.Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType == nil {
		return nil
	}

	if evtypes.FilterInactive(*evType) {
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

	// if autoclaim disabled, then event definitely active - fulfill passport scan events
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

	if len(countriesList) == 0 {
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

		// Not all events can be claimed, because limit can be reached in half path
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
				evType.Reward,
				evType.IgnoreCountryLimit)
			if err != nil {
				return fmt.Errorf("failed to do claim event updates for passport scan: %w", err)
			}

			// we mark claimed events to fulfill event which can't be claimed because of country limitations
			countriesBalancesMap[country.Code][i].EventStatus = data.EventClaimed
		}
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

// updateReferralUserEvents is used to add events for referrers
// for friends who have scanned the passport, if they have not been added.
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

// claimReferralSpecificEvents claim fulfilled events for invited
// friends which have passport scanned, if it possible
func claimReferralSpecificEvents(db *pgdb.DB, types evtypes.Types, levels config.Levels) error {
	evType := types.Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evType == nil || !evType.AutoClaim {
		return nil
	}

	events, err := pg.NewEvents(db).
		FilterByType(evtypes.TypeReferralSpecific).
		FilterByStatus(data.EventFulfilled).
		Select()
	if err != nil {
		return fmt.Errorf("failed to select passport scan events: %w", err)
	}

	// we need to have maps which link nullifiers to events slice and countries to balances slice
	nullifiersEventsMap := make(map[string][]data.Event, len(events))
	nullifiers := make([]string, 0, len(events))
	for _, event := range events {
		if _, ok := nullifiersEventsMap[event.Nullifier]; !ok {
			nullifiersEventsMap[event.Nullifier] = make([]data.Event, 0, len(events))
			nullifiers = append(nullifiers, event.Nullifier)
		}
		nullifiersEventsMap[event.Nullifier] = append(nullifiersEventsMap[event.Nullifier], event)
	}

	if len(nullifiers) == 0 {
		return nil
	}

	balances, err := pg.NewBalances(db).FilterByNullifier(nullifiers...).Select()
	if err != nil {
		return fmt.Errorf("failed to select balances for claim passport scan event: %w", err)
	}
	if len(balances) == 0 {
		return errors.New("critical: events present, but no balances with nullifier")
	}

	countriesBalancesMap := make(map[string][]data.Balance, len(balances))
	for _, balance := range balances {
		if !balance.ReferredBy.Valid || balance.Country == nil {
			continue
		}
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

	// toClaim - event ids which must be claimed
	toClaim := make([]string, 0, len(events))
	for _, country := range countries {
		// if country have limitations - skip this
		if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
			continue
		}

		limit := country.ReserveLimit - country.Reserved
		for _, balance := range countriesBalancesMap[country.Code] {
			// if limit reached we need stop
			if limit <= 0 {
				break
			}

			toAccrue := int64(0)
			for _, event := range nullifiersEventsMap[balance.Nullifier] {
				limit -= evType.Reward
				toClaim = append(toClaim, event.ID)
				toAccrue += evType.Reward
				if limit <= 0 {
					break
				}
			}

			err = handlers.DoClaimEventUpdates(
				levels,
				pg.NewReferrals(db),
				pg.NewBalances(db),
				pg.NewCountries(db),
				balance,
				evType.Reward,
				evType.IgnoreCountryLimit)
			if err != nil {
				return fmt.Errorf("failed to do claim event updates for referral specific event: %w", err)
			}
		}
	}

	if len(toClaim) == 0 {
		return nil
	}

	_, err = pg.NewEvents(db).FilterByID(toClaim...).Update(data.EventClaimed, nil, &evType.Reward)
	if err != nil {
		return fmt.Errorf("update event status: %w", err)
	}

	return nil
}
