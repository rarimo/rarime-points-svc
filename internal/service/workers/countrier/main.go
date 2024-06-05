package countrier

import (
	"context"
	"fmt"

	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type dbaser struct {
	db *pgdb.DB
}

func Run(ctx context.Context, cfg config.Config) {
	log := cfg.Log().WithField("who", "countrier")
	db := dbaser{cfg.DB().Clone()}

	countries, err := db.countriesQ().Select()
	if err != nil {
		panic(fmt.Errorf("failed to select countries: %w", err))
	}

	if len(countries) == 0 {
		log.Info("No countries in database")
	}

	toUpdate, toInsert := compareCountries(cfg.Countries(), countries)

	err = db.countriesQ().UpdateMany(toUpdate)
	if err != nil {
		panic(fmt.Errorf("failed to update countries: %w", err))
	}
	log.Infof("%d countries config was updated", len(toUpdate))

	err = db.countriesQ().Insert(toInsert...)
	if err != nil {
		panic(fmt.Errorf("failed to insert countries: %w", err))
	}
	log.Infof("%d countries config was inserted", len(toInsert))
}

func (db *dbaser) countriesQ() data.CountriesQ {
	return pg.NewCountries(db.db)
}

func compareCountries(cfgCountries config.Countries, dbCountries []data.Country) (toUpdate []data.Country, toInsert []data.Country) {
	toUpdate = make([]data.Country, 0, len(cfgCountries)+len(dbCountries))
	toInsert = make([]data.Country, 0, len(cfgCountries)+len(dbCountries))
	dbCodes := make(map[string]string, len(dbCountries))
	for _, v := range dbCountries {
		dbCodes[v.Code] = ""
		country := cfgCountries[data.DefaultCountryCode]
		if _, ok := cfgCountries[v.Code]; ok {
			country = cfgCountries[v.Code]
		}

		if v.ReserveLimit != country.ReserveLimit ||
			v.ReserveAllowed != country.ReserveAllowed ||
			v.WithdrawalAllowed != country.WithdrawalAllowed {

			toUpdate = append(toUpdate, data.Country{
				Code:              v.Code,
				ReserveLimit:      country.ReserveLimit,
				ReserveAllowed:    country.ReserveAllowed,
				WithdrawalAllowed: country.WithdrawalAllowed,
			})
		}
	}

	for code, country := range cfgCountries {
		if code == data.DefaultCountryCode {
			continue
		}
		if _, ok := dbCodes[code]; !ok {
			toInsert = append(toInsert, data.Country{
				Code:              code,
				ReserveLimit:      country.ReserveLimit,
				ReserveAllowed:    country.ReserveAllowed,
				WithdrawalAllowed: country.WithdrawalAllowed,
			})
		}
	}

	return toUpdate, toInsert
}
