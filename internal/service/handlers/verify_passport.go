package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/pkg/connector"
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
		"user_did":    req.UserDID,
		"hash":        req.Hash,
		"shared_data": req.SharedData,
	})

	balance, err := BalancesQ(r).FilterByPassportHash(req.Hash).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by Hash")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if balance != nil && balance.DID != req.UserDID {
		log.Error("passport_hash already in use")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	if balance == nil {
		balance, err = BalancesQ(r).FilterByDID(req.UserDID).Get()
		if err != nil {
			log.WithError(err).Error("Failed to get balance by DID")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	var reward *int64
	evType := EventTypes(r).Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evType != nil {
		var success bool
		reward, success = EventTypes(r).CalculatePassportScanReward(req.SharedData...)
		if !success {
			log.Error("Failed to calculate PassportScanReward, incorrect fields")
			ape.RenderErr(w, problems.NotFound())
			return
		}
	}

	if balance == nil {
		log.Debug("Balance not found, creating new one")

		err = createBalanceWithPassportTx(r, req, reward)

		if err != nil {
			log.WithError(err).Error("Failed to create balance with events")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	if !balance.PassportHash.Valid {
		err = setBalancePassportTx(r, req, reward, balance.ReferredBy)

		if err != nil {
			log.WithError(err).Error("Failed to set passport and add event for referrer")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = BalancesQ(r).FilterByDID(req.UserDID).SetPassport(balance.PassportHash.String, time.Now().UTC().AddDate(0, 1, 0))
	if err != nil {
		log.WithError(err).Error("Failed to set expiration date")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createBalanceWithPassportTx(r *http.Request, req connector.VerifyPassportRequest, reward *int64) error {
	events := EventTypes(r).PrepareEvents(req.UserDID, evtypes.FilterNotOpenable)

	log := Log(r).WithFields(map[string]any{
		"user_did":    req.UserDID,
		"hash":        req.Hash,
		"shared_data": req.SharedData,
	})

	for i := 0; i < len(events); i++ {
		if events[i].Type == evtypes.TypePassportScan {
			events[i].PointsAmount = reward
			events[i].Status = data.EventFulfilled
		}
	}

	return EventsQ(r).Transaction(func() (err error) {
		balance := &data.Balance{
			DID:                 req.UserDID,
			PassportHash:        sql.NullString{String: req.Hash, Valid: true},
			PassportExpires:     sql.NullTime{Time: time.Now().UTC().AddDate(0, 1, 0), Valid: true},
			IsWithdrawalAllowed: !req.IsUSA,
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
}

func setBalancePassportTx(r *http.Request, req connector.VerifyPassportRequest, reward *int64, refBy sql.NullString) error {
	log := Log(r).WithFields(map[string]any{
		"user_did":    req.UserDID,
		"hash":        req.Hash,
		"shared_data": req.SharedData,
	})
	return EventsQ(r).Transaction(func() error {
		_, err := BalancesQ(r).FilterByDID(req.UserDID).
			Update(data.Balance{
				ReferredBy:          refBy,
				PassportHash:        sql.NullString{String: req.Hash, Valid: true},
				PassportExpires:     sql.NullTime{Time: time.Now().UTC().AddDate(0, 1, 0), Valid: true},
				IsWithdrawalAllowed: !req.IsUSA})

		if err != nil {
			return fmt.Errorf("set passport for balance by DID: %w", err)
		}

		logMsgScan := "PassportScan event type is disabled or expired, not accruing points"
		if reward != nil {
			logMsgScan = "PassportScan event type available"
			if err = fulFillPassportScanEvent(r, req, reward); err != nil {
				return fmt.Errorf("fulfill passport scan event for user: %w", err)
			}
		}
		log.Debug(logMsgScan)

		if !refBy.Valid {
			log.Debug("User balance incative")
			return nil
		}

		evType := EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
		if evType == nil {
			log.Debug("Referral event type is disabled or expired, not accruing points to referrer")
			return nil
		}

		ref, err := ReferralsQ(r).Get(refBy.String)
		if err != nil {
			return fmt.Errorf("get referral: %w", err)
		}

		// normally should never happen
		if ref == nil {
			return fmt.Errorf("referral code not found")
		}

		err = EventsQ(r).Insert(data.Event{
			UserDID: ref.UserDID,
			Type:    evType.Name,
			Status:  data.EventFulfilled,
			Meta:    data.Jsonb(fmt.Sprintf(`{"did": "%s"}`, req.UserDID)),
		})
		if err != nil {
			return fmt.Errorf("add event for referrer: %w", err)
		}

		return nil
	})
}

func fulFillPassportScanEvent(r *http.Request, req connector.VerifyPassportRequest, reward *int64) error {
	log := Log(r).WithFields(map[string]any{
		"user_did":    req.UserDID,
		"hash":        req.Hash,
		"shared_data": req.SharedData,
	})

	passportScanEvent, err := EventsQ(r).
		FilterByUserDID(req.UserDID).
		FilterByType(evtypes.TypePassportScan).
		FilterByStatus(data.EventOpen).
		Get()

	if err != nil {
		return fmt.Errorf("get passport_scan event by DID: %w", err)
	}

	if passportScanEvent != nil {
		log.Debug("PassportScan event open")

		_, err = EventsQ(r).
			FilterByUserDID(req.UserDID).
			FilterByType(evtypes.TypePassportScan).
			Update(data.EventFulfilled, json.RawMessage(passportScanEvent.Meta), reward)
		if err != nil {
			return fmt.Errorf("update reward for passport_scan event by DID: %w", err)
		}
		log.Debug("PassportScan event reward update successful")
		return nil
	}
	log.Debug("PassportScan event not open")
	return nil
}
