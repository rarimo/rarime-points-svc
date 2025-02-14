package handlers

import (
	"net/http"

	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
)

func MaintenanceHandler(w http.ResponseWriter, r *http.Request) {
	ape.Render(w, resources.MaintenanceResponse{
		Data: resources.Maintenance{
			Key: resources.Key{
				Type: resources.MAINTENANCE,
			},
			Attributes: resources.MaintenanceAttributes{
				Maintenance: MaintenanceConfig(r).IsMaintenance,
			},
		},
	})
}
