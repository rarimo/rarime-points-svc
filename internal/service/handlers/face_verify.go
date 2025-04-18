package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

const (
	RootSMT = iota
	FaceChallengedNullifier
)

func FaceVerify(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewFaceScanVerify(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	proof := req.Data.Attributes.Proof
	nullifier := UserClaims(r)[0].Nullifier

	log := Log(r).WithFields(map[string]any{
		"nullifier": nullifier,
		"proof":     proof,
	})

	// if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) ||
	// 	new(big.Int).SetBytes(hexutil.MustDecode(nullifier)).String() != proof.PubSignals[FaceChallengedNullifier] {
	// 	log.Debug("failed to authenticate user")
	// 	ape.RenderErr(w, problems.Unauthorized())
	// 	return
	// }

	balance, err := BalancesQ(r).FilterByNullifier(nullifier).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by nullifier")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if balance == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	evType := EventTypes(r).Get(evtypes.TypeFaceParticipation, evtypes.FilterInactive)
	if evType == nil {
		log.Infof("Event face participation type is inactive")
		ape.RenderErr(w, problems.Forbidden())
		return
	}

	faceUserEvents, err := EventsQ(r).FilterByNullifier(nullifier).FilterByType(evtypes.TypeFaceParticipation).Select()
	if err != nil {
		log.WithError(err).Error("Failed to get user face events")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if len(faceUserEvents) != 0 {
		faceID := struct {
			RootSMT string `json:"root_smt"`
		}{}

		for _, event := range faceUserEvents {
			err := json.Unmarshal(event.Meta, &faceID)
			if err != nil {
				log.WithError(err).Errorf("Failed to parse event meta with eventID=%s", event.ID)
				ape.RenderErr(w, problems.InternalError())
				return
			}

			if faceID.RootSMT == proof.PubSignals[RootSMT] {
				log.Debugf("Face event already fulfilled")
				ape.RenderErr(w, problems.Conflict())
				return
			}
		}
	}

	err = FaceVerifier(r).VerifyProof(proof)
	if err != nil {
		log.WithError(err).Debug("Failed to verify face participation proof")
		if errors.Is(err, config.ErrInvalidRoot) {
			ape.RenderErr(w, problems.BadRequest(validation.Errors{
				"proof": err,
			})...)
			return
		}

		ape.RenderErr(w, problems.InternalError())
		return
	}

	_, err = EventsQ(r).FilterByNullifier(nullifier).FilterByType(evtypes.TypeFaceParticipation).FilterByStatus(data.EventOpen).Update(
		data.EventFulfilled,
		json.RawMessage(fmt.Sprintf(`{"root_smt": "%s"}`, proof.PubSignals[RootSMT])),
		nil,
	)
	if err != nil {
		log.WithError(err).Error("Failed to insert poll event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if !evType.AutoClaim {
		log.Debug("Event fulfilled due to disabled auto-claim")
		ape.Render(w, newEventClaimingStateResponse(balance.Nullifier, false))
		return
	}

	err = EventsQ(r).Transaction(func() error {
		return autoClaimEventsForFaceEventBalance(r, balance)
	})
	if err != nil {
		log.WithError(err).Error("Failed to autoclaim events for user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newEventClaimingStateResponse(balance.Nullifier, true))
}

func createFaceEventBalance(nullifier string, r *http.Request) (*data.FaceEventBalance, error) {

	faceEventBalance := data.FaceEventBalance{
		Nullifier: nullifier,
	}

	err := FaceEventBalancesQ(r).Insert(faceEventBalance)
	if err != nil {
		return nil, fmt.Errorf("add face balance: %w", err)
	}

	return &faceEventBalance, nil
}

func newEventClaimingStateResponse(id string, claimed bool) resources.FaceEventState {
	var res resources.FaceEventState
	res.Key.ID = id
	res.Key.Type = resources.FACE_EVENT_STATE
	res.Attributes.Claimed = claimed
	return res
}

func autoClaimEventsForFaceEventBalance(r *http.Request, balance *data.Balance) error {
	if balance == nil {
		Log(r).Debug("Balance absent. Events not claimed.")
		return nil
	}

	var totalPoints int64
	eventsToClaim, err := EventsQ(r).
		FilterByNullifier(balance.Nullifier).
		FilterByStatus(data.EventFulfilled, data.EventOpen).
		FilterByType(evtypes.TypeFaceParticipation).
		Select()
	if err != nil {
		return fmt.Errorf("failed to select events for user=%s: %w", balance.Nullifier, err)
	}

	eventsMap := map[string][]string{}
	for _, e := range eventsToClaim {
		if _, ok := eventsMap[e.Type]; !ok {
			eventsMap[e.Type] = []string{}
		}
		eventsMap[e.Type] = append(eventsMap[e.Type], e.ID)
	}

	for evName, evIDs := range eventsMap {
		evType := EventTypes(r).Get(evName, evtypes.FilterInactive, evtypes.FilterByAutoClaim(true))
		if evType == nil {
			continue
		}

		_, err = EventsQ(r).FilterByID(evIDs...).Update(data.EventClaimed, nil, &evType.Reward)
		if err != nil {
			return fmt.Errorf("failedt to update %s events for user=%s: %w", evName, balance.Nullifier, err)
		}

		totalPoints += evType.Reward * int64(len(evIDs))
	}

	level, err := doLvlUpAndReferralsUpdate(Levels(r), ReferralsQ(r), *balance, totalPoints)
	if err != nil {
		return fmt.Errorf("failed to do lvlup and referrals updates: %w", err)
	}

	err = BalancesQ(r).FilterByNullifier(balance.Nullifier).Update(map[string]any{
		data.ColLevel: level,
	})
	if err != nil {
		return fmt.Errorf("update level: %w", err)
	}

	err = FaceEventBalancesQ(r).FilterByNullifier(balance.Nullifier).Update(map[string]any{
		data.ColAmount: pg.AddToValue(data.ColAmount, totalPoints),
	})
	if err != nil {
		return fmt.Errorf("update face amount: %w", err)
	}

	return nil
}
