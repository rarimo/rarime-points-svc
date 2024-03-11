package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func VerifyPassport(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVerifyPassport(r)
	if err != nil {
		Log(r).WithError(err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	log := Log(r).WithFields(map[string]any{
		"user_did": req.UserDID,
		"hash":     req.Hash,
		"expiry":   req.Expiry.String(),
	})

	balance, err := BalancesQ(r).FilterByPassportHash(req.Hash).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance != nil {
		if balance.DID != req.UserDID {
			log.Error("passport_hash already in use")
			ape.RenderErr(w, problems.Conflict())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	balance, err = BalancesQ(r).FilterByDID(req.UserDID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by DID")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var reward int64
	logMsgScan := "PassportScan event type is disabled or expired, not accruing points"
	evType := EventTypes(r).Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType != nil {
		var success bool
		reward, success = EventTypes(r).CalculatePassportScanReward(req.SharedData...)
		if !success {
			log.WithError(err).Error("Failed to calculate PassportScanReward, incorrect fields")
			ape.RenderErr(w, problems.NotFound())
			return
		}
	}

	if balance == nil {
		log.Debug("Balance not found, creating new one")
		events := EventTypes(r).PrepareEvents(req.UserDID, evtypes.FilterNotOpenable)

		for i := 0; i < len(events); i++ {
			if events[i].Type == evtypes.TypePassportScan {
				events[i].PointsAmount = &reward
				events[i].Status = data.EventFulfilled
			}
		}

		err = EventsQ(r).Transaction(func() error {
			balance = &data.Balance{
				DID:             req.UserDID,
				PassportHash:    sql.NullString{String: req.Hash, Valid: true},
				PassportExpires: sql.NullTime{Time: req.Expiry, Valid: true},
			}

			if err = BalancesQ(r).Insert(*balance); err != nil {
				return fmt.Errorf("add balance: %w", err)
			}

			log.Debugf("%d events will be added for user_did=%s", len(events), req.UserDID)
			if err = EventsQ(r).Insert(events...); err != nil {
				return fmt.Errorf("add open events: %w", err)
			}
			return nil
		})

		if err != nil {
			log.WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = EventsQ(r).Transaction(func() error {
		// If you make this endpoint public, you should check the passport hash for
		// uniqueness and provide a better validation. Think about other changes too.
		err = BalancesQ(r).FilterByDID(req.UserDID).SetPassport(req.Hash, req.Expiry)
		if err != nil {
			return fmt.Errorf("set passport for balance by DID: %w", err)
		}

		if evType != nil {
			passportScanEvent, err := EventsQ(r).
				FilterByUserDID(req.UserDID).
				FilterByType(evtypes.TypePassportScan).
				FilterByStatus(data.EventOpen).
				Get()
			if err != nil {
				return fmt.Errorf("get passport_scan event by DID: %w", err)
			}
			logMsgOpenE := "PassportScan event not open"
			if passportScanEvent != nil {
				_, err = EventsQ(r).
					FilterByUserDID(req.UserDID).
					FilterByType(evtypes.TypePassportScan).
					Update(data.EventFulfilled, json.RawMessage(passportScanEvent.Meta), &reward)
				if err != nil {
					return fmt.Errorf("update reward for passport_scan event by DID: %w", err)
				}
				logMsgOpenE = "PassportScan event open"
				logMsgScan = "PassportScan event reward update successful"
			}
			log.Debug(logMsgOpenE)
		}
		log.Debug(logMsgScan)

		evType = EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
		if evType == nil {
			log.Debug("Referral event type is disabled or expired, not accruing points to referrer")
			return nil
		}

		refDID := balance.ReferredBy.String
		if refDID == "" {
			return nil
		}

		err = EventsQ(r).Insert(data.Event{
			UserDID: refDID,
			Type:    evType.Name,
			Status:  data.EventFulfilled,
			Meta:    data.Jsonb(fmt.Sprintf(`{"did": "%s"}`, req.UserDID)),
		})
		if err != nil {
			return fmt.Errorf("add event for referrer: %w", err)
		}

		return nil
	})

	if err != nil {
		log.WithError(err).Error("Failed to set passport and add event for referrer")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
