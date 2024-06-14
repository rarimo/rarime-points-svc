package main_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	zkptypes "github.com/iden3/go-rapidsnark/types"
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

	genesisCode = "kPRQYQUcWzW"

	balancesEndpoint = "public/balances"
	eventsEndpoint   = "public/events"
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

func TestCreateBalance(t *testing.T) {
	t.Run("SimpleBalance", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000001"
		body := createBalanceBody(nullifier, genesisCode)
		_, respCode := postPatchRequest(t, balancesEndpoint, body, nullifier, false)
		if respCode != http.StatusOK {
			t.Errorf("failed to create simple balance: want %d got %d", http.StatusOK, respCode)
		}
	})

	t.Run("SameBalance", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000001"
		body := createBalanceBody(nullifier, genesisCode)
		_, respCode := postPatchRequest(t, balancesEndpoint, body, nullifier, false)
		if respCode != http.StatusConflict {
			t.Errorf("want %d got %d", http.StatusConflict, respCode)
		}
	})

	t.Run("Unauthorized", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000002"
		body := createBalanceBody(nullifier, genesisCode)
		_, respCode := postPatchRequest(t, balancesEndpoint, body, "0x1"+nullifier[3:], false)
		if respCode != http.StatusUnauthorized {
			t.Errorf("want %d got %d", http.StatusUnauthorized, respCode)
		}
	})

	t.Run("IncorrectCode", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000002"
		body := createBalanceBody(nullifier, "someAntoherCode")
		_, respCode := postPatchRequest(t, balancesEndpoint, body, nullifier, false)
		if respCode != http.StatusNotFound {
			t.Errorf("want %d got %d", http.StatusNotFound, respCode)
		}
	})
}

func TestVerifyPassport(t *testing.T) {
	nullifier := "0x0000000000000000000000000000000000000000000000000000000000000002"
	referrer := "0x0000000000000000000000000000000000000000000000000000000000000001"

	balance := getBalance(t, referrer)
	if balance.Data.Attributes.ActiveReferralCodes == nil ||
		len(*balance.Data.Attributes.ActiveReferralCodes) == 0 {
		t.Fatalf("active referral codes for user %s absent", referrer)
	}
	createBalance(t, nullifier, (*balance.Data.Attributes.ActiveReferralCodes)[0])

	proof := baseProof
	proof.PubSignals[zk.Citizenship] = ukrCode
	body := verifyPassportBody(nullifier, proof)

	t.Run("VerifyPassport", func(t *testing.T) {
		_, respCode := postPatchRequest(t, balancesEndpoint+"/"+nullifier+"/verifypassport", body, nullifier, false)
		if respCode != http.StatusNoContent {
			t.Errorf("failed to verify passport: want %d got %d", http.StatusNoContent, respCode)
		}
	})

	t.Run("VerifyOneMore", func(t *testing.T) {
		_, respCode := postPatchRequest(t, balancesEndpoint+"/"+nullifier+"/verifypassport", body, nullifier, false)
		if respCode != http.StatusTooManyRequests {
			t.Errorf("want %d got %d", http.StatusTooManyRequests, respCode)
		}
	})

	t.Run("IncorrectCountryCode", func(t *testing.T) {
		proof.PubSignals[zk.Citizenship] = "6974819"
		body = verifyPassportBody(referrer, proof)
		_, respCode := postPatchRequest(t, balancesEndpoint+"/"+referrer+"/verifypassport", body, referrer, false)
		if respCode != http.StatusInternalServerError {
			t.Errorf("want %d got %d", http.StatusInternalServerError, respCode)
		}
	})
}

func TestClaimEvent(t *testing.T) {
	nullifier1 := "0x0000000000000000000000000000000000000000000000000000000000000010"
	nullifier2 := "0x0000000000000000000000000000000000000000000000000000000000000020"

	balance1 := createBalance(t, nullifier1, genesisCode)
	if balance1.Data.Attributes.ActiveReferralCodes == nil ||
		len(*balance1.Data.Attributes.ActiveReferralCodes) == 0 {
		t.Fatalf("active referral codes for user %s absent", nullifier1)
	}

	passportScanEventID, _ := getEventFromList(getEvents(t, nullifier1), evtypes.TypePassportScan)
	if passportScanEventID == "" {
		t.Fatalf("passport scan event absent for %s", nullifier1)
	}

	t.Run("TryClaimOpenEvent", func(t *testing.T) {
		body := claimEventBody(passportScanEventID)
		_, respCode := postPatchRequest(t, eventsEndpoint+"/"+passportScanEventID, body, nullifier1, true)
		if respCode != http.StatusNotFound {
			t.Errorf("want %d got %d", http.StatusNotFound, respCode)
		}
	})

	createBalance(t, nullifier2, (*balance1.Data.Attributes.ActiveReferralCodes)[0])
	verifyPassport(t, nullifier2, ukrCode)

	refSpecEventID, _ := getEventFromList(getEvents(t, nullifier1), evtypes.TypeReferralSpecific)
	if refSpecEventID == "" {
		t.Fatalf("referral specific event absent for %s", nullifier1)
	}

	t.Run("TryClaimEventWithoutPassport", func(t *testing.T) {
		body := claimEventBody(refSpecEventID)
		_, respCode := postPatchRequest(t, eventsEndpoint+"/"+refSpecEventID, body, nullifier1, true)
		if respCode != http.StatusForbidden {
			t.Errorf("want %d got %d", http.StatusForbidden, respCode)
		}
	})

	passportScanEventID, _ = getEventFromList(getEvents(t, nullifier2), evtypes.TypePassportScan)
	if passportScanEventID == "" {
		t.Fatalf("passport scan event absent for %s", nullifier2)
	}

	t.Run("ClaimEvent", func(t *testing.T) {
		body := claimEventBody(passportScanEventID)
		_, respCode := postPatchRequest(t, eventsEndpoint+"/"+passportScanEventID, body, nullifier2, true)
		if respCode != http.StatusOK {
			t.Errorf("want %d got %d", http.StatusOK, respCode)
		}
	})
}

func TestLevels(t *testing.T) {
	nullifier := "0x0000000000000000000000000000000000000000000000000000000000000100"

	balance := createBalance(t, nullifier, genesisCode)
	if balance.Data.Attributes.Level != 1 {
		t.Fatalf("balance level must be 1, got %d: %s", balance.Data.Attributes.Level, nullifier)
	}

	verifyPassport(t, nullifier, ukrCode)

	passportScanEventID, _ := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
	if passportScanEventID == "" {
		t.Fatalf("passport scan event absent for %s", nullifier)
	}

	claimEvent(t, passportScanEventID, nullifier)

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
	if len(*balance.Data.Attributes.ActiveReferralCodes) != 15 {
		t.Fatalf("balance referral codes must be 15, got %d: %s", len(*balance.Data.Attributes.ActiveReferralCodes), nullifier)
	}
}

func TestCountryPools(t *testing.T) {
	nullifier := "0x0000000000000000000000000000000000000000000000000000000000001000"

	createBalance(t, nullifier, genesisCode)
	verifyPassport(t, nullifier, usaCode)

	t.Run("UnderLimit", func(t *testing.T) {
		passportScanEventID, _ := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
		if passportScanEventID == "" {
			t.Fatalf("passport scan event absent for %s", nullifier)
		}

		claimEvent(t, passportScanEventID, nullifier)
	})

	t.Run("OverLimit", func(t *testing.T) {
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

	nullifier = "0x0000000000000000000000000000000000000000000000000000000000002000"

	createBalance(t, nullifier, genesisCode)
	verifyPassport(t, nullifier, gbrCode)

	t.Run("NotAllowedReserve", func(t *testing.T) {
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

	nullifier = "0x0000000000000000000000000000000000000000000000000000000000003000"

	createBalance(t, nullifier, genesisCode)
	verifyPassport(t, nullifier, deuCode)

	t.Run("DefaultUnderLimit", func(t *testing.T) {
		passportScanEventID, _ := getEventFromList(getEvents(t, nullifier), evtypes.TypePassportScan)
		if passportScanEventID == "" {
			t.Fatalf("passport scan event absent for %s", nullifier)
		}

		claimEvent(t, passportScanEventID, nullifier)
	})

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
		t.Errorf("want %d got %d", http.StatusOK, respCode)
	}

	var event resources.EventResponse
	err := json.Unmarshal(respBody, &event)
	if err != nil {
		t.Fatalf("failed to unmarhal event response: %v", err)
	}

	return event
}

func verifyPassport(t *testing.T, nullifier, country string) {
	proof := baseProof
	proof.PubSignals[zk.Citizenship] = country
	body := verifyPassportBody(nullifier, proof)

	_, respCode := postPatchRequest(t, balancesEndpoint+"/"+nullifier+"/verifypassport", body, nullifier, false)
	if respCode != http.StatusNoContent {
		t.Errorf("failed to verify passport: want %d got %d", http.StatusNoContent, respCode)
	}
}

func getEvents(t *testing.T, nullifier string) resources.EventListResponse {
	respBody, respCode := getRequest(t,
		eventsEndpoint, func() url.Values {
			query := url.Values{}
			query.Add("filter[nullifier]", nullifier)
			return query
		}(), nullifier)
	if respCode != http.StatusOK {
		t.Errorf("failed to get events: want %d got %d", http.StatusOK, respCode)
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
		t.Errorf("failed to get balance: want %d got %d", http.StatusOK, respCode)
	}

	var balance resources.BalanceResponse
	err := json.Unmarshal(respBody, &balance)
	if err != nil {
		t.Fatalf("failed to unmarhal balance response: %v", err)
	}

	return balance
}

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

var apiURL = func() string {
	var cfg struct {
		Addr string `fig:"addr,required"`
	}
	err := figure.Out(&cfg).From(kv.MustGetStringMap(kv.MustFromEnv(), "listener")).Please()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("http://%s/integrations/rarime-points-svc/v1/", cfg.Addr)
}()
