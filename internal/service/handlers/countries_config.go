package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetCountriesConfig(w http.ResponseWriter, r *http.Request) {
	countries, err := CountriesQ(r).FilterDisabled(false).Select()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get enabled countries")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	allowed := make([]string, 0, len(countries))
	limitReached := make([]string, 0, len(countries))

	for _, c := range countries {
		if c.Reserved < c.ReserveLimit {
			allowed = append(allowed, c.Code)
			continue
		}
		limitReached = append(limitReached, c.Code)
	}

	ape.Render(w, resources.CountriesConfigResponse{
		Data: resources.CountriesConfig{
			Key: resources.Key{
				Type: resources.COUNTRIES_CONFIG,
			},
			Attributes: resources.CountriesConfigAttributes{
				Allowed:      make([]string, 0, len(countries)),
				LimitReached: make([]string, 0, len(countries)),
			},
		},
	})
}
