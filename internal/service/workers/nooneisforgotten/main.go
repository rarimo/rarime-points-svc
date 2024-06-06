package nooneisforgotten

import (
	"context"
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
)

func Run(_ context.Context, cfg config.Config) {
	evType := cfg.EventTypes().Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType != nil {
		if err := updatePassportScanEvents(cfg); err != nil {
			panic(err)
		}
	}

	evType = cfg.EventTypes().Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evType != nil {
		if err := updateReferrUserEvents(cfg); err != nil {
			panic(err)
		}
	}
}

func updatePassportScanEvents(cfg config.Config) error {
	balances, err := pg.NewBalances(cfg.DB().Clone()).WithoutPassportEvent()
	if err != nil {
		return fmt.Errorf("failed to select balances without points for passport scan: %w", err)
	}

	toInsert := make([]data.Event, 0, len(balances))
	toUpdate := make([]string, 0, len(balances))
	for _, balance := range balances {
		if balance.EventID != nil {
			toUpdate = append(toUpdate, *balance.EventID)
			continue
		}

		toInsert = append(toInsert, data.Event{
			Nullifier: balance.Nullifier,
			Type:      evtypes.TypePassportScan,
			Status:    data.EventFulfilled,
		})
	}
	if len(toInsert) != 0 {
		err = pg.NewEvents(cfg.DB().Clone()).Insert(toInsert...)
		if err != nil {
			return fmt.Errorf("failed to insert passport scan events: %w", err)
		}
	}

	if len(toUpdate) != 0 {
		_, err = pg.NewEvents(cfg.DB().Clone()).FilterByID(toUpdate...).Update(data.EventFulfilled, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to update passport scan events: %w", err)
		}
	}

	return nil
}

func updateReferrUserEvents(cfg config.Config) error {
	referrs, err := pg.NewBalances(cfg.DB().Clone()).WithoutReferralEvent()
	if err != nil {
		return fmt.Errorf("failed to select balances without points for referred users: %w", err)
	}

	toInsert := make([]data.Event, 0, len(referrs))
	for _, referr := range referrs {
		toInsert = append(toInsert, data.Event{
			Nullifier: referr.Referrer,
			Type:      evtypes.TypeReferralSpecific,
			Status:    data.EventFulfilled,
			Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, referr.Referred)),
		})
	}

	if len(toInsert) != 0 {
		err = pg.NewEvents(cfg.DB().Clone()).Insert(toInsert...)
		if err != nil {
			return fmt.Errorf("failed to insert referred user events: %w", err)
		}
	}

	return nil
}
