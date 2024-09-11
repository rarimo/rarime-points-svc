package cli

import (
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/handlers"
)

func emitEvent(cfg config.Config, timestamp int) {
	log := cfg.Log()
	db := cfg.DB()
	lvls := cfg.Levels()
	evTypes := cfg.EventTypes()

	balancesQ := pg.NewBalances(db)
	eventsQ := pg.NewEvents(db)
	referralsQ := pg.NewReferrals(db)
	countriesQ := pg.NewCountries(db)

	evType := evTypes.Get(evtypes.TypeEarlyTest, evtypes.FilterInactive)

	if evType == nil {
		log.Infof("Event type %s is inactive", evtypes.TypeEarlyTest)
		return
	}

	balances, err := balancesQ.FilterByCreatedAtBefore(timestamp).FilterUnverified().Select()
	if err != nil {
		panic(fmt.Errorf("failed to select balances for early test reward: %w", err))
	}
	if len(balances) == 0 {
		log.Infof("No balances found")
		return
	}

	nullifiers := make([]string, 0, len(balances))
	for _, balance := range balances {
		nullifiers = append(nullifiers, balance.Nullifier)
	}

	emittedEvents, err := eventsQ.New().FilterByType(evtypes.TypeEarlyTest).FilterByNullifier(nullifiers...).Select()
	if err != nil {
		panic(fmt.Errorf("failed to select emitted events: %w", err))
	}

	eventsMap := make(map[string]struct{}, len(emittedEvents))
	for _, event := range emittedEvents {
		eventsMap[event.Nullifier] = struct{}{}
	}

	for _, balance := range balances {
		err = eventsQ.New().Transaction(func() error {
			if _, exists := eventsMap[balance.Nullifier]; exists {
				log.Infof("Event %s is already done for user with nullifier %s ", evtypes.TypeEarlyTest, balance.Nullifier)
				return nil
			}

			err = eventsQ.Insert(data.Event{
				Nullifier: balance.Nullifier,
				Type:      evtypes.TypeEarlyTest,
				Status:    data.EventFulfilled,
			})

			if err != nil {
				return fmt.Errorf("failed to insert %s event: %w", evtypes.TypeEarlyTest, err)
			}

			if !evType.AutoClaim {
				return nil
			}

			_, err = eventsQ.FilterByNullifier(balance.Nullifier).Update(data.EventClaimed, nil, &evType.Reward)
			if err != nil {
				return fmt.Errorf("failed to update %s events for user=%s: %w", evtypes.TypeEarlyTest, balance.Nullifier, err)
			}

			err := handlers.DoClaimEventUpdates(lvls, referralsQ, balancesQ, countriesQ, balance, evType.Reward)
			if err != nil {
				return fmt.Errorf("failed to do lvlup and referrals updates: %w", err)
			}

			return nil
		})
	}
}
