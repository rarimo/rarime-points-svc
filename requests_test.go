package main_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

const requestTimeout = time.Second // use bigger on debug with breakpoints to prevent fails

const (
	ukrCode = "5589842"
	usaCode = "5591873"
	gbrCode = "4670034"
	deuCode = "4474197"
	canCode = "4407630"
	fraCode = "4608577"
	indCode = "4804164"
	mcoCode = "5063503"

	genesisBalance = "0x0000000000000000000000000000000000000000000000000000000000000000"

	balancesEndpoint = "public/balances"
	eventsEndpoint   = "public/events"
)

var (
	apiURL      = ""
	genesisCode = "kPRQYQUcWzW"
)

var (
	balancesPath = func() string {
		return "public/balances"
	}
	balancesSpecificPath = func(nullifier string) string {
		return "public/balances/" + nullifier
	}
	verifyPassportPath = func(nullifier string) string {
		return balancesSpecificPath(nullifier) + "/verifypassport"
	}
	witdhrawalsPath = func(nullifier string) string {
		return balancesSpecificPath(nullifier) + "/withdrawals"
	}
	countriesConfigPath = func() string {
		return "public/countries_config"
	}
	eventTypesPath = func() string {
		return "public/event_types"
	}
	eventsPath = func() string {
		return "public/events"
	}
	eventsSpecificPath = func(id string) string {
		return "public/events/" + id
	}
	fulfillEventPath = func() string {
		return "private/events"
	}
	editReferralsPath = func() string {
		return "private/referrals"
	}
)

var baseProof = zkptypes.ZKProof{
	Proof: &zkptypes.ProofData{
		A:        []string{"0", "0", "0"},
		B:        []([]string){{"0", "0"}, {"0", "0"}, {"0", "0"}},
		C:        []string{"0", "0", "0"},
		Protocol: "groth16",
	},
	PubSignals: []string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"},
}

func TestMain(m *testing.M) {
	if err := setUp(); err != nil {
		panic(fmt.Errorf("failed to setup: %w", err))
	}

	exitVal := m.Run()

	if err := tearDown(); err != nil {
		panic(fmt.Errorf("failed to teardown: %w", err))
	}
	os.Exit(exitVal)
}

func setUp() error {
	if err := setApiURL(); err != nil {
		return fmt.Errorf("failed to set api URL: %w", err)
	}

	return nil
}

func tearDown() error {
	return nil
}

func setApiURL() error {
	var cfg struct {
		Addr string `fig:"addr,required"`
	}

	err := figure.Out(&cfg).From(kv.MustGetStringMap(kv.MustFromEnv(), "listener")).Please()
	if err != nil {
		return fmt.Errorf("failed to figure out listener from service config: %w", err)
	}

	apiURL = fmt.Sprintf("http://%s/integrations/rarime-points-svc/v1/", cfg.Addr)
	return nil
}

func TestCreateBalance(t *testing.T) {
	t.Run("SimpleBalance", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000001"
		body := createBalanceBody(nullifier, genesisCode)
		_, respCode := postPatchRequest(t, balancesPath(), body, nullifier, false)
		if respCode != http.StatusOK {
			t.Errorf("failed to create simple balance: want %d got %d", http.StatusOK, respCode)
		}
	})

	t.Run("SameBalance", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000001"
		body := createBalanceBody(nullifier, genesisCode)
		_, respCode := postPatchRequest(t, balancesPath(), body, nullifier, false)
		if respCode != http.StatusConflict {
			t.Errorf("want %d got %d", http.StatusConflict, respCode)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000002"
		body := createBalanceBody(nullifier, genesisCode)
		_, respCode := postPatchRequest(t, balancesPath(), body, "0x1"+nullifier[3:], false)
		if respCode != http.StatusUnauthorized {
			t.Errorf("want %d got %d", http.StatusUnauthorized, respCode)
		}
	})

	t.Run("IncorrectCode", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000002"
		body := createBalanceBody(nullifier, "someAntoherCode")
		_, respCode := postPatchRequest(t, balancesPath(), body, nullifier, false)
		if respCode != http.StatusNotFound {
			t.Errorf("want %d got %d", http.StatusNotFound, respCode)
		}
	})
}

func TestVerifyPassport(t *testing.T) {
	t.Run("VerifyPassport", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000010"
		createBalance(t, nullifier, genesisCode)

		proof := baseProof
		proof.PubSignals[zk.Citizenship] = ukrCode
		body := verifyPassportBody(nullifier, proof)

		_, respCode := postPatchRequest(t, verifyPassportPath(nullifier), body, nullifier, false)
		if respCode != http.StatusOK {
			t.Errorf("failed to verify passport: want %d got %d", http.StatusOK, respCode)
		}
	})

	t.Run("VerifyPassportSecondTime", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000010"

		proof := baseProof
		proof.PubSignals[zk.Citizenship] = ukrCode
		body := verifyPassportBody(nullifier, proof)

		// depend on previous test, because balance for nullifier must exist
		_, respCode := postPatchRequest(t, verifyPassportPath(nullifier), body, nullifier, false)
		if respCode != http.StatusTooManyRequests {
			t.Errorf("want %d got %d", http.StatusTooManyRequests, respCode)
		}
	})

	t.Run("IncorrectCountryCode", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000020"
		createBalance(t, nullifier, genesisCode)

		proof := baseProof
		proof.PubSignals[zk.Citizenship] = "6974819"
		body := verifyPassportBody(nullifier, proof)

		_, respCode := postPatchRequest(t, verifyPassportPath(nullifier), body, nullifier, false)
		if respCode != http.StatusInternalServerError {
			t.Errorf("want %d got %d", http.StatusInternalServerError, respCode)
		}
	})
}

func TestAutoClaimEvent(t *testing.T) {
	t.Run("SuccessClaimPassportScan", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000100"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, usaCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
		if eventID == "" {
			t.Log("passport scan event absent")
			return
		}
		if eventStatus != string(data.EventClaimed) {
			t.Fatalf("want passport scan event status %s got %s", data.EventClaimed, eventStatus)
		}
	})

	t.Run("ReservedDisallowed", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000200"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, gbrCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
		if eventID == "" {
			t.Log("passport scan event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want passport scan event status %s got %s", data.EventFulfilled, eventStatus)
		}
	})

	// this test depend on `SuccessClaimPassportScan`
	t.Run("ReserveLimitReached", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000300"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, usaCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
		if eventID == "" {
			t.Log("passport scan event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want passport scan event status %s got %s", data.EventFulfilled, eventStatus)
		}
	})

	t.Run("CountryBanned", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000400"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, canCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
		if eventID == "" {
			t.Log("passport scan event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want passport scan event status %s got %s", data.EventFulfilled, eventStatus)
		}
	})

	t.Run("ReferralSpecific", func(t *testing.T) {
		referrer := "0x0000000000000000000000000000000000000000000000000000000000000500"
		referrerBalance := createBalance(t, referrer, genesisCode)
		verifyPassport(t, referrer, ukrCode)
		if referrerBalance.Data.Attributes.ReferralCodes == nil || len(*referrerBalance.Data.Attributes.ReferralCodes) == 0 {
			t.Fatal("referrer's referral codes must exists")
		}

		if (*referrerBalance.Data.Attributes.ReferralCodes)[0].Status != data.StatusActive {
			t.Fatal("first referrer's referral code inactive")
		}

		referred := "0x0000000000000000000000000000000000000000000000000000000000000600"
		createBalance(t, referred, (*referrerBalance.Data.Attributes.ReferralCodes)[0].Id)
		verifyPassport(t, referred, ukrCode)

		eventID, eventStatus := getEventFromList(getEvents(t, referrer), evtypes.TypeReferralSpecific)
		if eventID == "" {
			t.Log("referral specific event absent")
			return
		}
		if eventStatus != string(data.EventClaimed) {
			t.Fatalf("want referral specific event status %s got %s", data.EventClaimed, eventStatus)
		}
	})

	// User can have a lot unclaimed fulfilled referral specific events if user not scan passport
	t.Run("ReferralSpecifics", func(t *testing.T) {
		referrer := "0x0000000000000000000000000000000000000000000000000000000000000700"
		referrerBalance := createBalance(t, referrer, genesisCode)
		if referrerBalance.Data.Attributes.ReferralCodes == nil || len(*referrerBalance.Data.Attributes.ReferralCodes) < 2 {
			t.Fatal("referrer's referral codes must exists")
		}

		if (*referrerBalance.Data.Attributes.ReferralCodes)[0].Status != data.StatusActive {
			t.Fatal("first referrer's referral code inactive")
		}

		referred1 := "0x0000000000000000000000000000000000000000000000000000000000000800"
		createBalance(t, referred1, (*referrerBalance.Data.Attributes.ReferralCodes)[0].Id)
		verifyPassport(t, referred1, ukrCode)

		if (*referrerBalance.Data.Attributes.ReferralCodes)[1].Status != data.StatusActive {
			t.Fatal("second referrer's referral code inactive")
		}

		referred2 := "0x0000000000000000000000000000000000000000000000000000000000000900"
		createBalance(t, referred2, (*referrerBalance.Data.Attributes.ReferralCodes)[1].Id)
		verifyPassport(t, referred2, ukrCode)

		fulfilledEventCount := 0
		referrerEvents := getEvents(t, referrer)
		for _, event := range referrerEvents.Data {
			if event.Attributes.Meta.Static.Name == evtypes.TypeReferralSpecific && event.Attributes.Status == string(data.EventFulfilled) {
				fulfilledEventCount += 1
			}
		}

		if fulfilledEventCount != 2 {
			t.Fatalf("count of fulfilled events for referrer must be 2")
		}

		verifyPassport(t, referrer, ukrCode)

		claimedEventCount := 0
		referrerEvents = getEvents(t, referrer)
		for _, event := range referrerEvents.Data {
			if event.Attributes.Meta.Static.Name == evtypes.TypeReferralSpecific && event.Attributes.Status == string(data.EventClaimed) {
				claimedEventCount += 1
			}
		}

		if claimedEventCount != 2 {
			t.Fatalf("count of claimed events for referrer must be 2")
		}
	})
}

func TestClaimEvent(t *testing.T) {
	t.Run("WithoutPassport", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000001000"
		createBalance(t, nullifier, genesisCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
		if eventID == "" {
			t.Log("free weekly event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want free weekly event status %s got %s", data.EventFulfilled, eventStatus)
		}

		body := claimEventBody(eventID)
		_, respCode := postPatchRequest(t, eventsSpecificPath(eventID), body, nullifier, true)
		if respCode != http.StatusForbidden {
			t.Errorf("want %d got %d", http.StatusForbidden, respCode)
		}
	})

	t.Run("SuccessClaim", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000002000"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, fraCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
		if eventID == "" {
			t.Log("free weekly event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want free weekly event status %s got %s", data.EventFulfilled, eventStatus)
		}

		body := claimEventBody(eventID)
		_, respCode := postPatchRequest(t, eventsSpecificPath(eventID), body, nullifier, true)
		if respCode != http.StatusOK {
			t.Errorf("want %d got %d", http.StatusOK, respCode)
		}
	})

	t.Run("ReserveDisallowed", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000003000"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, indCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
		if eventID == "" {
			t.Log("free weekly event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want free weekly event status %s got %s", data.EventFulfilled, eventStatus)
		}

		body := claimEventBody(eventID)
		_, respCode := postPatchRequest(t, eventsSpecificPath(eventID), body, nullifier, true)
		if respCode != http.StatusForbidden {
			t.Errorf("want %d got %d", http.StatusForbidden, respCode)
		}
	})

	// this test depend on `SuccessClaimPassportScan`
	t.Run("ReserveLimitReached", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000004000"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, fraCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
		if eventID == "" {
			t.Log("passport scan event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want free weekly event status %s got %s", data.EventFulfilled, eventStatus)
		}

		body := claimEventBody(eventID)
		_, respCode := postPatchRequest(t, eventsSpecificPath(eventID), body, nullifier, true)
		if respCode != http.StatusForbidden {
			t.Errorf("want %d got %d", http.StatusForbidden, respCode)
		}
	})

	t.Run("CountryBanned", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000005000"
		createBalance(t, nullifier, genesisCode)
		verifyPassport(t, nullifier, mcoCode)

		eventID, eventStatus := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
		if eventID == "" {
			t.Log("passport scan event absent")
			return
		}
		if eventStatus != string(data.EventFulfilled) {
			t.Fatalf("want free weekly event status %s got %s", data.EventFulfilled, eventStatus)
		}

		body := claimEventBody(eventID)
		_, respCode := postPatchRequest(t, eventsSpecificPath(eventID), body, nullifier, true)
		if respCode != http.StatusForbidden {
			t.Errorf("want %d got %d", http.StatusForbidden, respCode)
		}
	})
}

func TestLevels(t *testing.T) {
	nullifier := "0x0000000000000000000000000000000000000000000000000000000000010000"

	balance := createBalance(t, nullifier, genesisCode)
	if balance.Data.Attributes.Level != 1 {
		t.Fatalf("balance level must be 1, got %d: %s", balance.Data.Attributes.Level, nullifier)
	}

	verifyPassport(t, nullifier, ukrCode)

	balance = getBalance(t, nullifier)
	if balance.Data.Attributes.Level != 2 {
		t.Fatalf("balance level must be 2, got %d: %s", balance.Data.Attributes.Level, nullifier)
	}

	freeWeeklyEventID, _ := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
	if freeWeeklyEventID == "" {
		t.Fatalf("free weekly event absent for %s", nullifier)
	}

	claimEvent(t, freeWeeklyEventID, nullifier)

	balance = getBalance(t, nullifier)
	if balance.Data.Attributes.Level != 3 {
		t.Fatalf("balance level must be 3, got %d: %s", balance.Data.Attributes.Level, nullifier)
	}

	// must never panic because of logic getBalance
	if len(*balance.Data.Attributes.ReferralCodes) != 15 {
		t.Fatalf("balance referral codes must be 15, got %d: %s", len(*balance.Data.Attributes.ReferralCodes), nullifier)
	}
}

func TestCountryPoolsDefault(t *testing.T) {
	nullifier := "0x0000000000000000000000000000000000000000000000000000000000100000"

	createBalance(t, nullifier, genesisCode)
	verifyPassport(t, nullifier, deuCode)

	t.Run("DefaultOverLimit", func(t *testing.T) {
		freeWeeklyEventID, _ := getEventFromList(getEvents(t, nullifier), evtypes.TypeFreeWeekly)
		if freeWeeklyEventID == "" {
			t.Fatalf("free weekly event absent for %s", nullifier)
		}

		body := claimEventBody(freeWeeklyEventID)
		_, respCode := postPatchRequest(t, eventsEndpoint+"/"+freeWeeklyEventID, body, nullifier, true)
		if respCode != http.StatusForbidden {
			t.Errorf("want %d got %d", http.StatusForbidden, respCode)
		}
	})
}

func getEventFromList(events resources.EventListResponse, evtype string) (id, status string) {
	for _, event := range events.Data {
		if event.Attributes.Meta.Static.Name == evtype {
			return event.ID, event.Attributes.Status
		}
	}
	return "", ""
}

func claimEvent(t *testing.T, id, nullifier string) resources.EventResponse {
	body := claimEventBody(id)
	respBody, respCode := postPatchRequest(t, eventsEndpoint+"/"+id, body, nullifier, true)
	if respCode != http.StatusOK {
		t.Fatalf("want %d got %d", http.StatusOK, respCode)
	}

	var event resources.EventResponse
	err := json.Unmarshal(respBody, &event)
	if err != nil {
		t.Fatalf("failed to unmarhal event response: %v", err)
	}

	return event
}

func verifyPassport(t *testing.T, nullifier, country string) resources.PassportEventStateResponse {
	proof := baseProof
	proof.PubSignals[zk.Citizenship] = country
	body := verifyPassportBody(nullifier, proof)

	respBody, respCode := postPatchRequest(t, verifyPassportPath(nullifier), body, nullifier, false)
	if respCode != http.StatusOK {
		t.Fatalf("failed to verify passport: want %d got %d", http.StatusOK, respCode)
	}

	var resp resources.PassportEventStateResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal passport event state response: %v", err)
	}

	return resp
}

func getEvents(t *testing.T, nullifier string) resources.EventListResponse {
	respBody, respCode := getRequest(t,
		eventsEndpoint, func() url.Values {
			query := url.Values{}
			query.Add("filter[nullifier]", nullifier)
			return query
		}(), nullifier)
	if respCode != http.StatusOK {
		t.Fatalf("failed to get events: want %d got %d", http.StatusOK, respCode)
	}

	var events resources.EventListResponse
	err := json.Unmarshal(respBody, &events)
	if err != nil {
		t.Fatalf("failed to unmarhal event list response: %v", err)
	}

	return events
}

func createBalance(t *testing.T, nullifier, code string) resources.BalanceResponse {
	body := createBalanceBody(nullifier, code)
	respBody, respCode := postPatchRequest(t, balancesEndpoint, body, nullifier, false)
	if respCode != http.StatusOK {
		t.Fatalf("failed to create simple balance: want %d got %d", http.StatusOK, respCode)
	}

	var balance resources.BalanceResponse
	err := json.Unmarshal(respBody, &balance)
	if err != nil {
		t.Fatalf("failed to unmarhal balance response: %v", err)
	}

	return balance
}

func getBalance(t *testing.T, nullifier string) resources.BalanceResponse {
	respBody, respCode := getRequest(t,
		balancesEndpoint+"/"+nullifier,
		func() url.Values {
			query := url.Values{}
			query.Add("referral_codes", "true")
			query.Add("rank", "true")
			return query
		}(), nullifier)
	if respCode != http.StatusOK {
		t.Fatalf("failed to get balance: want %d got %d", http.StatusOK, respCode)
	}

	var balance resources.BalanceResponse
	err := json.Unmarshal(respBody, &balance)
	if err != nil {
		t.Fatalf("failed to unmarhal balance response: %v", err)
	}

	return balance
}

// func editReferrals(t *testing.T) {

// }

func verifyPassportBody(nullifier string, proof zkptypes.ZKProof) resources.VerifyPassportRequest {
	return resources.VerifyPassportRequest{
		Data: resources.VerifyPassport{
			Key: resources.Key{
				ID:   nullifier,
				Type: resources.VERIFY_PASSPORT,
			},
			Attributes: resources.VerifyPassportAttributes{
				Proof: proof,
			},
		},
	}
}

func createBalanceBody(nullifier, code string) resources.CreateBalanceRequest {
	return resources.CreateBalanceRequest{
		Data: resources.CreateBalance{
			Key: resources.Key{
				ID:   nullifier,
				Type: resources.CREATE_BALANCE,
			},
			Attributes: resources.CreateBalanceAttributes{
				ReferredBy: code,
			},
		},
	}
}

func claimEventBody(id string) resources.Relation {
	return resources.Relation{
		Data: &resources.Key{
			ID:   id,
			Type: resources.CLAIM_EVENT,
		},
	}
}

func postPatchRequest(t *testing.T, endpoint string, body any, user string, patch bool) ([]byte, int) {
	if body == nil {
		t.Fatal("request body not provided")
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal request bode: %v", err)
	}

	log.Printf("  endpoint=/%s  body=%s", endpoint, body)

	reqBody := strings.NewReader(string(bodyJSON))

	reqType := "POST"
	if patch {
		reqType = "PATCH"
	}

	req, err := http.NewRequest(reqType, apiURL+endpoint, reqBody)
	if err != nil {
		t.Fatalf("failed to create post request: %v", err)
	}

	if user != "" {
		req.Header.Set("nullifier", user)
	}

	resp, err := (&http.Client{Timeout: requestTimeout}).Do(req)
	if err != nil {
		t.Fatalf("failed to perform post request: %v", err)
	}
	defer func() {
		resp.Body.Close()
	}()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read resp body: %v", err)
	}

	log.Printf("  endpoint=/%s  body=%s", endpoint, respBody)

	return respBody, resp.StatusCode
}

func getRequest(t *testing.T, endpoint string, query url.Values, user string) ([]byte, int) {
	log.Printf("  endpoint=/%s  query=%+v", endpoint, query)

	req, err := http.NewRequest("GET", apiURL+endpoint, nil)
	if err != nil {
		t.Fatalf("failed to create get request: %v", err)
	}

	req.URL.RawQuery = query.Encode()

	if user != "" {
		req.Header.Set("nullifier", user)
	}

	resp, err := (&http.Client{Timeout: requestTimeout}).Do(req)
	if err != nil {
		t.Fatalf("failed to perform get request: %v", err)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read resp body: %v", err)
	}
	defer func() {
		resp.Body.Close()
	}()

	log.Printf("  endpoint=/%s  body=%s", endpoint, respBody)

	return respBody, resp.StatusCode
}
