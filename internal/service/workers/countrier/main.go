package countrier

import (
	"context"
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type extConfig interface {
	comfig.Logger
	pgdb.Databaser
	evtypes.EventTypeser
	Countrier
}

func Run(_ context.Context, cfg extConfig) {
	log := cfg.Log().WithField("who", "countrier")
	q := pg.NewCountries(cfg.DB().Clone())

	countries, err := q.New().Select() // running only once
	if err != nil {
		panic(fmt.Errorf("failed to select countries: %w", err))
	}

	if len(countries) == 0 {
		log.Info("No countries in database")
	}

	toUpdate, toInsert := compareCountries(cfg.Countries(), countries)

	err = q.New().UpdateMany(toUpdate)
	if err != nil {
		panic(fmt.Errorf("failed to update countries: %w", err))
	}
	log.Infof("%d countries config was updated", len(toUpdate))

	err = q.New().Insert(toInsert...)
	if err != nil {
		panic(fmt.Errorf("failed to insert countries: %w", err))
	}
	log.Infof("%d countries config was inserted", len(toInsert))
}

func compareCountries(cfgCountries Config, dbCountries []data.Country) (toUpdate []data.Country, toInsert []data.Country) {
	toUpdate = make([]data.Country, 0, len(cfgCountries.m)+len(dbCountries))
	toInsert = make([]data.Country, 0, len(cfgCountries.m)+len(dbCountries))
	dbCodes := make(map[string]string, len(dbCountries))
	for _, v := range dbCountries {
		dbCodes[v.Code] = ""
		c := cfgCountries.m[data.DefaultCountryCode]
		if _, ok := cfgCountries.m[v.Code]; ok {
			c = cfgCountries.m[v.Code]
		}

		if v.ReserveLimit != c.ReserveLimit ||
			v.ReserveAllowed != c.ReserveAllowed ||
			v.WithdrawalAllowed != c.WithdrawalAllowed {

			toUpdate = append(toUpdate, data.Country{
				Code:              v.Code,
				ReserveLimit:      c.ReserveLimit,
				ReserveAllowed:    c.ReserveAllowed,
				WithdrawalAllowed: c.WithdrawalAllowed,
			})
		}
	}

	for code, cou := range cfgCountries.m {
		if code == data.DefaultCountryCode {
			continue
		}
		if _, ok := dbCodes[code]; !ok {
			toInsert = append(toInsert, data.Country{
				Code:              code,
				ReserveLimit:      cou.ReserveLimit,
				ReserveAllowed:    cou.ReserveAllowed,
				WithdrawalAllowed: cou.WithdrawalAllowed,
			})
		}
	}

	return toUpdate, toInsert
}
