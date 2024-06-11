package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetCountriesConfig(w http.ResponseWriter, r *http.Request) {
	countries, err := CountriesQ(r).Select()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get enabled countries")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	properties := make([]resources.CountryProperties, 0, len(countries))
	for _, c := range countries {
		prop := resources.CountryProperties{
			Code:              c.Code,
			ReserveAllowed:    c.ReserveAllowed,
			WithdrawalAllowed: c.WithdrawalAllowed,
		}
		// when the limit is reached, reserve is not allowed despite the config
		if c.Reserved >= c.ReserveLimit {
			prop.ReserveAllowed = false
		}
		properties = append(properties, prop)
	}

	ape.Render(w, resources.CountriesConfigResponse{
		Data: resources.CountriesConfig{
			Key: resources.Key{
				Type: resources.COUNTRIES_CONFIG,
			},
			Attributes: resources.CountriesConfigAttributes{
				Countries: properties,
			},
		},
	})
}
