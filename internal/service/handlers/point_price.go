package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
)

func GetPointPrice(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, resources.PointPriceResponse{
		Data: resources.PointPrice{
			Key: resources.Key{
				Type: resources.POINT_PRICE,
			},
			Attributes: resources.PointPriceAttributes{
				Urmo: PointPrice(r),
			},
		},
	})
}
