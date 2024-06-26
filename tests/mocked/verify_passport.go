package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"net/http"

	"errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/jsonapi"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
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
		"balance.nullifier":    req.Data.ID,
		"balance.anonymous_id": req.Data.Attributes.AnonymousId,
		"balance.country":      req.Data.Attributes.Country,
	})

	var (
		country     = req.Data.Attributes.Country
		anonymousID = req.Data.Attributes.AnonymousId
		proof       = req.Data.Attributes.Proof

		gotSig  = r.Header.Get("Signature")
		wantSig = calculatePassportVerificationSignature(
			CountriesConfig(r).VerificationKey,
			req.Data.ID,
			country,
			anonymousID,
		)
	)

	if gotSig != wantSig {
		log.Warnf("Unauthorized access: HMAC signature mismatch: got %s, want %s", gotSig, wantSig)
		ape.RenderErr(w, problems.Forbidden())
		return
	}
	if proof == nil {
		log.Debug("Proof is not provided: performing logic of joining program instead of full verification")
	}

	balance, errs := getAndVerifyBalanceEligibility(r, req.Data.ID, proof)
	if len(errs) > 0 {
		ape.RenderErr(w, errs...)
		return
	}

	byAnonymousID, err := BalancesQ(r).FilterByAnonymousID(anonymousID).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get balance by anonymous ID")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if byAnonymousID != nil && byAnonymousID.Nullifier != balance.Nullifier {
		log.Warn("Balance with the same anonymous ID already exists")
		ape.RenderErr(w, problems.Conflict())
		return
	}

	if balance.Country != nil {
		if balance.IsPassportProven {
			log.Warnf("Balance %s already verified", balance.Nullifier)
			ape.RenderErr(w, problems.TooManyRequests())
			return
		}
		if proof == nil {
			log.Warnf("Balance %s tried to re-join program", balance.Nullifier)
			ape.RenderErr(w, problems.TooManyRequests())
			return
		}

		var balAID string
		if balance.AnonymousID != nil {
			balAID = *balance.AnonymousID
		}

		err = validation.Errors{
			"data/attributes/country":      validation.Validate(*balance.Country, validation.Required, validation.In(country)),
			"data/attributes/anonymous_id": validation.Validate(anonymousID, validation.Required, validation.In(balAID)),
		}.Filter()
		if err != nil {
			ape.RenderErr(w, problems.BadRequest(err)...)
			return
		}

		err = BalancesQ(r).FilterByNullifier(balance.Nullifier).Update(map[string]any{
			data.ColIsPassport: true,
		})
		if err != nil {
			log.WithError(err).Error("Failed to update balance")
			ape.RenderErr(w, problems.InternalError())
			return
		}

		ape.Render(w, newPassportEventStateResponse(req.Data.ID, nil))
		return
	}

	err = EventsQ(r).Transaction(func() error {
		return doPassportScanUpdates(r, *balance, country, anonymousID, proof != nil)
	})
	if err != nil {
		log.WithError(err).Error("Failed to execute transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	event, err := EventsQ(r).FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypePassportScan).
		FilterByStatus(data.EventClaimed).Get()
	if err != nil {
		log.WithError(err).Error("Failed to get claimed event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, newPassportEventStateResponse(req.Data.ID, event))
}

func newPassportEventStateResponse(id string, event *data.Event) resources.PassportEventStateResponse {
	var res resources.PassportEventStateResponse
	res.Data.ID = id
	res.Data.Type = resources.PASSPORT_EVENT_STATE
	res.Data.Attributes.Claimed = event != nil
	return res
}

func calculatePassportVerificationSignature(key []byte, nullifier, country, anonymousID string) string {
	bNull, err := hex.DecodeString(nullifier[2:])
	if err != nil {
		panic(fmt.Errorf("nullifier was not properly validated as hex: %w", err))
	}
	bAID, err := hex.DecodeString(anonymousID)
	if err != nil {
		panic(fmt.Errorf("anonymousID was not properly validated as hex: %w", err))
	}

	h := hmac.New(sha256.New, key)
	msg := append(bNull, []byte(country)...)
	msg = append(msg, bAID...)
	h.Write(msg)

	return hex.EncodeToString(h.Sum(nil))
}

// getAndVerifyBalanceEligibility provides shared logic to verify that the user
// is eligible to verify passport or withdraw. Some extra checks still exist in
// the flows. You may provide nil proof to handle its verification outside.
func getAndVerifyBalanceEligibility(
	r *http.Request,
	nullifier string,
	proof *zkptypes.ZKProof,
) (balance *data.Balance, errs []*jsonapi.ErrorObject) {

	if !auth.Authenticates(UserClaims(r), auth.UserGrant(nullifier)) {
		return nil, append(errs, problems.Unauthorized())
	}

	balance, err := BalancesQ(r).FilterByNullifier(nullifier).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get balance by nullifier")
		return nil, append(errs, problems.InternalError())
	}

	if errs = checkVerificationEligibility(r, balance); len(errs) > 0 {
		return nil, errs
	}
	// for withdrawal and joining program
	if proof == nil {
		return balance, nil
	}

	// never panics because of request validation
	// proof.PubSignals[zk.Nullifier] = mustHexToInt(nullifier)
	// err = Verifier(r).VerifyProof(*proof)
	// if err != nil {
	//	if errors.Is(err, identity.ErrContractCall) {
	//		Log(r).WithError(err).Error("Failed to verify proof")
	//		return nil, append(errs, problems.InternalError())
	//	}
	//	return nil, problems.BadRequest(err)
	// }

	return balance, nil
}

func checkVerificationEligibility(r *http.Request, balance *data.Balance) (errs []*jsonapi.ErrorObject) {
	switch {
	case balance == nil:
		Log(r).Debug("Balance absent")
		return append(errs, problems.NotFound())
	case !balance.ReferredBy.Valid:
		Log(r).Debug("Balance inactive")
		return append(errs, problems.BadRequest(validation.Errors{
			"referred_by": errors.New("user must be referred to withdraw"),
		})...)
	}

	return nil
}

// doPassportScanUpdates performs all the necessary updates when the passport
// scan proof is provided. This logic is shared between verification and
// withdrawal handlers.
func doPassportScanUpdates(r *http.Request, balance data.Balance, countryCode, anonymousID string, proven bool) error {
	country, err := updateBalanceCountry(r, balance, countryCode, anonymousID, proven)
	if err != nil {
		return fmt.Errorf("update balance country: %w", err)
	}
	if !country.ReserveAllowed || !country.WithdrawalAllowed || country.Reserved >= country.ReserveLimit {
		Log(r).Infof("User %s scanned passport which country has restrictions: %+v", balance.Nullifier, country)
	}

	// because for claim event must be country code
	balance.Country = &country.Code

	// Fulfill passport scan event for user if event active
	// Event can be automaticaly claimed if:
	// 1. Autoclaim enabled for passport scan event
	// 2. Reservation is allowed to the country
	// 3. The country reservation limit has not been reached
	if err = fulfillOrClaimPassportScanEvent(r, balance, *country); err != nil {
		return fmt.Errorf("fulfill passport scan event: %w", err)
	}

	// Type not filtered as inactive because expired events can be claimed
	evTypeRef := EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evTypeRef == nil {
		Log(r).Debug("Referral specific event type is inactive")
		return nil
	}

	// Claim events for invited friends who scanned the passport.
	// This is possible when the user registered in the referral
	// program and invited friends, the friends scanned the passport,
	// but since the user had an unscanned passport, the event
	// could not automatically claimed to him. And now that
	// user has scanned the passport, it is necessary to claim events
	// for user's friends, if possible, that is, the following conditions are met:
	// 1. Autoclaim enabled for passport scan event
	// 2. Reservation is allowed to the country
	// 3. The country reservation limit has not been reached
	if err = claimReferralSpecificEvents(r, evTypeRef, balance.Nullifier); err != nil {
		return fmt.Errorf("failed to claim referral specific events: %w", err)
	}

	// Adds a friend event for the referrer. If the event
	// is inactive, then nothing happens. If active, the
	// fulfilled event is added and, if possible, the event claimed
	if err = addEventForReferrer(r, evTypeRef, balance); err != nil {
		return fmt.Errorf("add event for referrer: %w", err)
	}

	return nil
}

func updateBalanceCountry(r *http.Request, balance data.Balance, code, anonymousID string, proven bool) (*data.Country, error) {
	country, err := getOrCreateCountry(CountriesQ(r), code)
	if err != nil {
		return nil, fmt.Errorf("get or create country: %w", err)
	}
	if balance.Country != nil {
		if *balance.Country == country.Code {
			return country, nil
		}
		// countries mismatch is handled separately in withdrawal flow before calling
		// updateBalanceCountry, so this will never happen
		return nil, errors.New("countries mismatch")
	}

	toUpd := map[string]any{
		data.ColCountry:     country.Code,
		data.ColAnonymousID: anonymousID,
	}
	if proven {
		toUpd[data.ColIsPassport] = true
	}

	err = BalancesQ(r).FilterByNullifier(balance.Nullifier).Update(toUpd)
	if err != nil {
		return nil, fmt.Errorf("update balance country: %w", err)
	}

	return country, nil
}

func fulfillOrClaimPassportScanEvent(r *http.Request, balance data.Balance, country data.Country) error {
	evTypePassport := EventTypes(r).Get(evtypes.TypePassportScan, evtypes.FilterInactive)
	if evTypePassport == nil {
		Log(r).Debug("Passport scan event type is inactive")
		return nil
	}

	event, err := EventsQ(r).FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypePassportScan).
		FilterByStatus(data.EventOpen).Get()
	if err != nil {
		return fmt.Errorf("get open passport scan event: %w", err)
	}

	if event == nil {
		return errors.New("inconsistent state: balance has no country, event type is active, but no open event was found")
	}

	if !evTypePassport.AutoClaim || !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
		_, err = EventsQ(r).
			FilterByID(event.ID).
			Update(data.EventFulfilled, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to update event: %w", err)
		}

		return nil
	}

	_, err = EventsQ(r).FilterByID(event.ID).Update(data.EventClaimed, nil, &evTypePassport.Reward)
	if err != nil {
		return fmt.Errorf("update event status: %w", err)
	}

	err = DoClaimEventUpdates(
		Levels(r),
		ReferralsQ(r),
		BalancesQ(r),
		CountriesQ(r),
		balance,
		evTypePassport.Reward)
	if err != nil {
		return fmt.Errorf("failed to do claim event updates for passport scan: %w", err)
	}

	return nil
}

// evTypeRef must not be nil
func claimReferralSpecificEvents(r *http.Request, evTypeRef *evtypes.EventConfig, nullifier string) error {
	if !evTypeRef.AutoClaim {
		Log(r).Debugf("auto claim for referral specific disabled")
		return nil
	}

	// balance can't be nil because of previous logic
	balance, err := BalancesQ(r).FilterByNullifier(nullifier).FilterDisabled().Get()
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	// country can't be nill because of previous logic
	country, err := CountriesQ(r).FilterByCodes(*balance.Country).Get()
	if err != nil {
		return fmt.Errorf("failed to get referrer country: %w", err)
	}

	// if user country have restrictions for claim points then not claim events and return
	if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
		Log(r).Debug("Country disallowed for reserve or limit was reached after passport scan")
		return nil
	}

	events, err := EventsQ(r).
		FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypeReferralSpecific).
		FilterByStatus(data.EventFulfilled).
		Select()
	if err != nil {
		return fmt.Errorf("get fulfilled referral specific events: %w", err)
	}

	// Specify how many events can be claimed
	countToClaim := int64(len(events))
	if countToClaim == 0 {
		return nil
	}

	// If, for example, 10 points are awarded for an event,
	// and 2 points remain before reaching the reservation
	// limit, then this event can be claimed. And since there
	// can be many events with invited friends, need to calculate
	// the maximum number of events that can be claimed in order
	// not to exceed the limit.
	if country.Reserved+countToClaim*evTypeRef.Reward >= country.ReserveLimit+evTypeRef.Reward {
		countToClaim = int64(math.Ceil(float64(country.ReserveLimit-country.Reserved) / float64(evTypeRef.Reward)))
	}

	eventsToClaimed := make([]string, countToClaim)
	for i := 0; i < int(countToClaim); i++ {
		eventsToClaimed[i] = events[i].ID
	}

	_, err = EventsQ(r).FilterByID(eventsToClaimed...).Update(data.EventClaimed, nil, &evTypeRef.Reward)
	if err != nil {
		return fmt.Errorf("update event status: %w", err)
	}

	err = DoClaimEventUpdates(
		Levels(r),
		ReferralsQ(r),
		BalancesQ(r),
		CountriesQ(r),
		*balance,
		countToClaim*evTypeRef.Reward)
	if err != nil {
		return fmt.Errorf("failed to do claim event updates for referral specific events: %w", err)
	}

	return nil
}

func addEventForReferrer(r *http.Request, evTypeRef *evtypes.EventConfig, balance data.Balance) error {
	if evTypeRef == nil {
		return nil
	}

	// ReferredBy always valid because of the previous logic
	referral, err := ReferralsQ(r).Get(balance.ReferredBy.String)
	if err != nil {
		return fmt.Errorf("get referral by ID: %w", err)
	}
	if referral == nil {
		return fmt.Errorf("critical: referred_by not null, but row in referrals absent")
	}

	referrerBalance, err := BalancesQ(r).FilterByNullifier(referral.Nullifier).Get()
	if err != nil {
		return fmt.Errorf("failed to get referrer balance: %w", err)
	}
	if referrerBalance == nil {
		return fmt.Errorf("critical: referrer balance not exist [%s], while referral code exist", referral.Nullifier)
	}

	if !referrerBalance.ReferredBy.Valid {
		Log(r).Debug("Referrer is genesis balance")
		return nil
	}

	if !evTypeRef.AutoClaim || referrerBalance.Country == nil {
		if referrerBalance.Country == nil {
			Log(r).Debug("Referrer not scan passport yet! Add fulfilled events")
		}
		err = EventsQ(r).Insert(data.Event{
			Nullifier: referral.Nullifier,
			Type:      evTypeRef.Name,
			Status:    data.EventFulfilled,
			Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, balance.Nullifier)),
		})
		if err != nil {
			return fmt.Errorf("failed to insert fulfilled event for referrer: %w", err)
		}

		return nil
	}

	country, err := CountriesQ(r).FilterByCodes(*referrerBalance.Country).Get()
	if err != nil {
		return fmt.Errorf("failed to get referrer country: %w", err)
	}
	if country == nil {
		return fmt.Errorf("critical: country must be present in database")
	}

	if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
		Log(r).Debug("Referrer country have ReserveAllowed false or limit reached")

		err = EventsQ(r).Insert(data.Event{
			Nullifier: referral.Nullifier,
			Type:      evTypeRef.Name,
			Status:    data.EventFulfilled,
			Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, balance.Nullifier)),
		})
		if err != nil {
			return fmt.Errorf("failed to insert fulfilled event for referrer: %w", err)
		}

		return nil
	}

	err = EventsQ(r).Insert(data.Event{
		Nullifier:    referral.Nullifier,
		Type:         evTypeRef.Name,
		Status:       data.EventClaimed,
		PointsAmount: &evTypeRef.Reward,
		Meta:         data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, balance.Nullifier)),
	})
	if err != nil {
		return fmt.Errorf("failed to insert claimed event for referrer: %w", err)
	}

	err = DoClaimEventUpdates(
		Levels(r),
		ReferralsQ(r),
		BalancesQ(r),
		CountriesQ(r),
		*referrerBalance,
		evTypeRef.Reward)
	if err != nil {
		return fmt.Errorf("failed to do claim event updates for referrer referral specific events: %w", err)
	}

	return nil
}

func getOrCreateCountry(q data.CountriesQ, code string) (*data.Country, error) {
	c, err := q.FilterByCodes(code).Get()
	if err != nil {
		return nil, fmt.Errorf("get country by code: %w", err)
	}
	if c != nil {
		return c, nil
	}

	def, err := q.New().FilterByCodes(data.DefaultCountryCode).Get()
	if err != nil {
		return nil, fmt.Errorf("get default country: %w", err)
	}
	if def == nil {
		return nil, errors.New("default country does not exist in DB")
	}

	c = &data.Country{
		Code:              code,
		ReserveLimit:      def.ReserveLimit,
		ReserveAllowed:    def.ReserveAllowed,
		WithdrawalAllowed: def.WithdrawalAllowed,
	}

	if err = q.New().Insert(*c); err != nil {
		return nil, fmt.Errorf("insert country with default values: %w", err)
	}

	return c, nil
}

func mustHexToInt(s string) string {
	return new(big.Int).SetBytes(hexutil.MustDecode(s)).String()
}
