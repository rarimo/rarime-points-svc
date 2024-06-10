package nooneisforgotten

import (
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func Run(cfg config.Config, sig chan struct{}) {
	db := cfg.DB().Clone()

	evType := cfg.EventTypes().Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType != nil {
		if err := updatePassportScanEvents(db); err != nil {
			panic(err)
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
