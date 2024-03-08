package handlers

import (
	"net/http"

	"github.com/rarimo/auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetEvent(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewGetEvent(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	event, err := EventsQ(r).FilterByID(req).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get event by ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if event == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(event.UserDID)) {
		ape.RenderErr(w, problems.Unauthorized())
		return
	}

	ape.Render(w, resources.EventResponse{Data: newEventModel(*event, EventTypes(r).Get(event.Type).Resource())})
}
