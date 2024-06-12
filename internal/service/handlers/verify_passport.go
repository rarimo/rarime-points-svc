package handlers

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/jsonapi"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/decentralized-auth-svc/pkg/auth"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func VerifyPassport(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVerifyPassport(r)
	if err != nil {
		Log(r).WithError(err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	balance, errs := getAndVerifyBalanceEligibility(r, req.Data.ID, &req.Data.Attributes.Proof)
	if len(errs) > 0 {
		ape.RenderErr(w, errs...)
		return
	}

	if balance.Country != nil {
		Log(r).Debugf("Balance %s already verified", balance.Nullifier)
		ape.RenderErr(w, problems.TooManyRequests())
		return
	}

	err = EventsQ(r).Transaction(func() error {
		return doPassportScanUpdates(r, *balance, req.Data.Attributes.Proof)
	})
	if err != nil {
		Log(r).WithError(err).Error("Failed to execute transaction")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	event, err := EventsQ(r).FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypePassportScan).
		FilterByStatus(data.EventClaimed).Get()
	if err != nil {
		Log(r).WithError(err).Error("Failed to get claimed event")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	var res resources.PassportEventStateResponse
	res.Data.ID = req.Data.ID
	res.Data.Type = resources.PASSPORT_EVENT_STATE
	res.Data.Attributes.Claimed = (event != nil)

	ape.Render(w, res)
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
	// for withdrawal
	if proof == nil {
		return balance, nil
	}

	// never panics because of request validation
	proof.PubSignals[zk.Nullifier] = mustHexToInt(nullifier)
	err = Verifier(r).VerifyProof(*proof)
	if err != nil {
		return nil, problems.BadRequest(err)
	}

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
func doPassportScanUpdates(r *http.Request, balance data.Balance, proof zkptypes.ZKProof) error {
	country, err := updateBalanceCountry(r, balance, proof)
	if err != nil {
		return fmt.Errorf("update balance country: %w", err)
	}
	if !country.ReserveAllowed || !country.WithdrawalAllowed || country.Reserved >= country.ReserveLimit {
		Log(r).Infof("User %s scanned passport which country has restrictions: %+v", balance.Nullifier, country)
	}

	// because for claim event must be country code
	balance.Country = &country.Code

	if err = fulfillPassportScanEvent(r, balance, *country); err != nil {
		return fmt.Errorf("fulfill passport scan event: %w", err)
	}

	if err = addEventForReferrer(r, balance); err != nil {
		return fmt.Errorf("add event for referrer: %w", err)
	}

	if err = claimReferralSpecificEvents(r, balance, true); err != nil {
		return fmt.Errorf("failed to claim referral specific events: %w", err)
	}

	return nil
}

func updateBalanceCountry(r *http.Request, balance data.Balance, proof zkptypes.ZKProof) (*data.Country, error) {
	country, err := getOrCreateCountry(CountriesQ(r), proof)
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

	err = BalancesQ(r).FilterByNullifier(balance.Nullifier).Update(map[string]any{
		data.ColCountry: country.Code,
	})
	if err != nil {
		return nil, fmt.Errorf("update balance country: %w", err)
	}

	return country, nil
}

func fulfillPassportScanEvent(r *http.Request, balance data.Balance, country data.Country) error {
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

	event, err = EventsQ(r).
		FilterByID(event.ID).
		Update(data.EventFulfilled, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	if !evTypePassport.AutoClaim || !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
		return nil
	}

	_, err = claimEventWithPoints(r, *event, evTypePassport.Reward, &balance)
	return err
}

// if autoClaim true - events will be claimed if autoclaim enabled for event
// if false - claim events without check if autoclaim enabled
func claimReferralSpecificEvents(r *http.Request, balance data.Balance, autoClaim bool) error {
	evTypeRef := EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evTypeRef == nil {
		Log(r).Debug("Referral specific event type is inactive")
		return nil
	}

	if autoClaim && !evTypeRef.AutoClaim {
		return nil
	}

	country, err := CountriesQ(r).FilterByCodes(*balance.Country).Get()
	if err != nil {
		return fmt.Errorf("failed to get referrer country: %w", err)
	}

	if country == nil {
		return fmt.Errorf("failed to get referrer country: must be present in database")
	}

	events, err := EventsQ(r).FilterByNullifier(balance.Nullifier).
		FilterByType(evtypes.TypeReferralSpecific).
		FilterByStatus(data.EventFulfilled).Select()
	if err != nil {
		return fmt.Errorf("get fulfilled referral specific events: %w", err)
	}

	if len(events) == 0 {
		return nil
	}

	if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
		return nil
	}

	for _, event := range events {
		if _, err = claimEventWithPoints(r, event, evTypeRef.Reward, &balance); err != nil {
			return fmt.Errorf("failed to claim referral specific event: %w", err)
		}
	}
	return nil

}

func addEventForReferrer(r *http.Request, balance data.Balance) error {
	evTypeRef := EventTypes(r).Get(evtypes.TypeReferralSpecific, evtypes.FilterInactive)
	if evTypeRef == nil {
		Log(r).Debug("Referral event type is inactive, not fulfilling event for referrer")
		return nil
	}

	// ReferredBy always valid because of the previous logic
	referral, err := ReferralsQ(r).Get(balance.ReferredBy.String)
	if err != nil {
		return fmt.Errorf("get referral by ID: %w", err)
	}

	event, err := EventsQ(r).InsertOne(data.Event{
		Nullifier: referral.Nullifier,
		Type:      evTypeRef.Name,
		Status:    data.EventFulfilled,
		Meta:      data.Jsonb(fmt.Sprintf(`{"nullifier": "%s"}`, balance.Nullifier)),
	})

	if err != nil {
		return fmt.Errorf("failed to insert event for referrer: %w", err)
	}

	if !evTypeRef.AutoClaim {
		Log(r).Debugf("auto claim for referral specific disabled")
		return nil
	}

	referrerBalance, err := BalancesQ(r).FilterByNullifier(referral.Nullifier).Get()
	if err != nil {
		return fmt.Errorf("failed to get referrer balance: %w", err)
	}

	if referrerBalance == nil {
		return fmt.Errorf("referrer balance not exists: %s", referral.Nullifier)
	}

	// genesis address have codes, but haven't referred_by
	if !referrerBalance.ReferredBy.Valid || referrerBalance.Country == nil {
		Log(r).Debug("Referrer have invalid referred_by or not scan passport")
		return nil
	}

	country, err := CountriesQ(r).FilterByCodes(*referrerBalance.Country).Get()
	if err != nil {
		return fmt.Errorf("failed to get referrer country: %w", err)
	}

	if country == nil {
		return fmt.Errorf("failed to get referrer country: must be present in database")
	}

	if !country.ReserveAllowed || country.Reserved >= country.ReserveLimit {
		Log(r).Debug("Referrer country have ReserveAllowed false or limit reached")
		return nil
	}

	_, err = claimEventWithPoints(r, *event, evTypeRef.Reward, referrerBalance)
	if err != nil {
		return fmt.Errorf("failed to claim referral_specific event for referrer: %w", err)
	}

	return nil
}

func getOrCreateCountry(q data.CountriesQ, proof zkptypes.ZKProof) (*data.Country, error) {
	code, err := extractCountry(proof)
	if err != nil {
		return nil, fmt.Errorf("extract country: %w", err)
	}

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

// extractCountry extracts 3-letter country code from the proof.
func extractCountry(proof zkptypes.ZKProof) (string, error) {
	b, ok := new(big.Int).SetString(proof.PubSignals[zk.Citizenship], 10)
	if !ok {
		b = new(big.Int)
	}

	code := string(b.Bytes())

	return code, validation.Errors{
		"code": validation.Validate(
			code,
			validation.Required,
			validation.When(code != data.DefaultCountryCode, is.CountryCode3),
		)}.Filter()
}

func mustHexToInt(s string) string {
	return new(big.Int).SetBytes(hexutil.MustDecode(s)).String()
}
