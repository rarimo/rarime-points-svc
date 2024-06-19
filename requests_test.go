package main_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/jsonapi"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	requestTimeout    = time.Second // use bigger on debug with breakpoints to prevent fails
	defaultConfigFile = "config.local.yaml"

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
	globalCfg             config.Config
	apiURL                string
	genesisCode           string
	nullifiers            []string
	currentNullifierIndex int
)

var baseProof = zkptypes.ZKProof{
	Proof: &zkptypes.ProofData{
		A:        []string{"0", "0", "0"},
		B:        [][]string{{"0", "0"}, {"0", "0"}, {"0", "0"}},
		C:        []string{"0", "0", "0"},
		Protocol: "groth16",
	},
	PubSignals: []string{"0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"},
}

func TestMain(m *testing.M) {
	var exitVal int
	defer func() {
		if r := recover(); r != nil {
			log.Printf("tests panicked: %v\n%s", r, debug.Stack())
			exitVal = 1
		}
		os.Exit(exitVal)
	}()

	setUp()
	exitVal = m.Run()
	tearDown()
}

func tearDown() {
}

func setUp() {
	if os.Getenv(kv.EnvViperConfigFile) == "" {
		err := os.Setenv(kv.EnvViperConfigFile, defaultConfigFile)
		if err != nil {
			panic(fmt.Errorf("failed to set env: %w", err))
		}
	}

	globalCfg = config.New(kv.MustFromEnv())
	apiURL = fmt.Sprintf("http://%s/integrations/rarime-points-svc/v1", globalCfg.Listener().Addr().String())

	refs, err := editReferrals(genesisBalance, 20)
	if err != nil {
		panic(fmt.Errorf("failed to edit referrals: %w", err))
	}
	genesisCode = refs.Ref

	nullifiers = make([]string, 20)
	for i := range nullifiers {
		hash := sha256.Sum256([]byte{byte(i)})
		nullifiers[i] = hexutil.Encode(hash[:])
	}
}

func TestCreateBalance(t *testing.T) {
	var (
		nullifierShared = nextN()
		otRefCode       string
	)

	validBalanceChecks := func(t *testing.T, nullifier, code string) {
		resp, err := createBalance(nullifier, code)
		require.NoError(t, err)
		require.Equal(t, nullifier, resp.Data.ID)

		attr := resp.Data.Attributes

		require.NotNil(t, attr.IsDisabled)
		require.NotNil(t, attr.IsVerified)
		require.NotNil(t, attr.ReferralCodes)
		require.NotEmpty(t, *attr.ReferralCodes)

		assert.Equal(t, 0, attr.Amount)
		assert.False(t, *attr.IsDisabled)
		assert.False(t, *attr.IsVerified)
		assert.Equal(t, 1, attr.Level)
		assert.NotNil(t, attr.Rank)

		otRefCode = (*attr.ReferralCodes)[0].Id
		require.NotEmpty(t, otRefCode)
	}

	// fixme @violog: looks like fail on assert/require won't stop outer tests, must check before proceeding
	t.Run("BalanceGenesisCode", func(t *testing.T) {
		validBalanceChecks(t, nullifierShared, genesisCode)
	})

	t.Run("BalanceOneTimeCode", func(t *testing.T) {
		validBalanceChecks(t, nextN(), otRefCode)
	})

	t.Run("SameBalanceConflict", func(t *testing.T) {
		_, err := createBalance(nullifierShared, genesisCode)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "409", apiErr.Status)
	})

	t.Run("NullifierUnauthorized", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		body := createBalanceBody(n1, genesisCode)
		err := requestWithBody(balancesEndpoint, "POST", n2, body, nil)

		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "401", apiErr.Status)
	})

	t.Run("ConsumedCode", func(t *testing.T) {
		_, err := createBalance(nextN(), otRefCode)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "404", apiErr.Status)
	})

	t.Run("IncorrectCode", func(t *testing.T) {
		_, err := createBalance(nextN(), "invalid")
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "404", apiErr.Status)
	})
}

func TestVerifyPassport(t *testing.T) {
	t.Run("VerifyPassport", func(t *testing.T) {
		nullifier := "0x0000000000000000000000000000000000000000000000000000000000000010"
		createBalance(t, nullifier, genesisCode)

		proof := baseProof
		proof.PubSignals[zk.Citizenship] = ukrCode
		body := verifyPassportBody(nullifier, proof)
	})
	t.Run("VerifyPassport", func(t *testing.T) {
		_, respCode := requestWithBody(t, balancesEndpoint+"/"+nullifier+"/verifypassport", body, nullifier, false)
		if respCode != http.StatusNoContent {
			t.Errorf("failed to verify passport: want %d got %d", http.StatusNoContent, respCode)
		}
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

func TestEventsAutoClaim(t *testing.T) {
	t.Run("PassportScanAutoclaim", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respVerifyStatus, err := verifyPassport(n, usaCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respBalance, err := getBalance(n)
		require.NoError(t, err)
		require.Equal(t, 2, respBalance.Data.Attributes.Level)
		require.Equal(t, 5, respBalance.Data.Attributes.Amount)
		require.NotNil(t, respBalance.Data.Attributes.ReferralCodes)
		require.Equal(t, 10, len(*respBalance.Data.Attributes.ReferralCodes))
	})

	// this test depend on previous test
	t.Run("PassportScanLimitReached", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respVerifyStatus, err := verifyPassport(n, usaCode)
		require.NoError(t, err)
		require.False(t, respVerifyStatus.Data.Attributes.Claimed)
	})

	t.Run("PassportScanReserveDisallowed", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respVerifyStatus, err := verifyPassport(n, gbrCode)
		require.NoError(t, err)
		require.False(t, respVerifyStatus.Data.Attributes.Claimed)
	})

	t.Run("PassportScanCountryBanned", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respVerifyStatus, err := verifyPassport(n, canCode)
		require.NoError(t, err)
		require.False(t, respVerifyStatus.Data.Attributes.Claimed)
	})

	t.Run("ReferralSpecificAutoclaim", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		respBalance, err := createBalance(n1, genesisCode)
		require.NoError(t, err)
		require.NotNil(t, respBalance.Data.Attributes.ReferralCodes)
		require.NotEmpty(t, (*respBalance.Data.Attributes.ReferralCodes))

		respVerifyStatus, err := verifyPassport(n1, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respBalance, err = createBalance(n2, (*respBalance.Data.Attributes.ReferralCodes)[0].Id)
		require.NoError(t, err)

		_, err = verifyPassport(n2, ukrCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n1)
		require.NoError(t, err)
		_, status := getEventFromList(respEvents, evtypes.TypeReferralSpecific)
		require.Equal(t, data.EventClaimed, status)
	})

	// User can have a lot unclaimed fulfilled referral specific events if user not scan passport
	t.Run("ReferralSpecificsAutoclaim", func(t *testing.T) {
		n1, n2, n3 := nextN(), nextN(), nextN()
		respBalance, err := createBalance(n1, genesisCode)
		require.NoError(t, err)
		require.NotNil(t, respBalance.Data.Attributes.ReferralCodes)
		require.GreaterOrEqual(t, (*respBalance.Data.Attributes.ReferralCodes), 2)

		respBalance, err = createBalance(n2, (*respBalance.Data.Attributes.ReferralCodes)[0].Id)
		require.NoError(t, err)
		_, err = verifyPassport(n2, ukrCode)
		require.NoError(t, err)

		respBalance, err = createBalance(n3, (*respBalance.Data.Attributes.ReferralCodes)[1].Id)
		require.NoError(t, err)
		_, err = verifyPassport(n3, ukrCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n1)
		require.NoError(t, err)
		fulfilledEventCount := 0
		for _, event := range respEvents.Data {
			if event.Attributes.Meta.Static.Name == evtypes.TypeReferralSpecific && event.Attributes.Status == string(data.EventFulfilled) {
				fulfilledEventCount += 1
			}
		}
		require.Equal(t, 2, fulfilledEventCount)

		respVerifyStatus, err := verifyPassport(n1, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respEvents, err = getEvents(n1)
		require.NoError(t, err)
		claimedEventCount := 0
		for _, event := range respEvents.Data {
			if event.Attributes.Meta.Static.Name == evtypes.TypeReferralSpecific && event.Attributes.Status == string(data.EventFulfilled) {
				claimedEventCount += 1
			}
		}
		require.Equal(t, 2, claimedEventCount)
	})
}

func TestClaimEvent(t *testing.T) {
	t.Run("WithoutPassport", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n)
		require.NoError(t, err)
		eventID, status := getEventFromList(respEvents, evtypes.TypeFreeWeekly)
		require.Equal(t, data.EventFulfilled, status)

		_, err = claimEvent(eventID, n)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})

	t.Run("IncorrectEventID", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, ukrCode)
		require.NoError(t, err)

		_, err = claimEvent("event", n)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "404", apiErr.Status)
	})

	t.Run("EventClaim", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, fraCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n)
		require.NoError(t, err)
		eventID, status := getEventFromList(respEvents, evtypes.TypeFreeWeekly)
		require.Equal(t, data.EventFulfilled, status)

		respEvent, err := claimEvent(eventID, n)
		require.NoError(t, err)
		require.Equal(t, data.EventClaimed, respEvent.Data.Attributes.Status)
	})

	// this test depend on previous test
	t.Run("ReserveLimitReached", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, fraCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n)
		require.NoError(t, err)
		eventID, status := getEventFromList(respEvents, evtypes.TypeFreeWeekly)
		require.Equal(t, data.EventFulfilled, status)

		_, err = claimEvent(eventID, n)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})

	t.Run("ReserveDisallowed", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, indCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n)
		require.NoError(t, err)
		eventID, status := getEventFromList(respEvents, evtypes.TypeFreeWeekly)
		require.Equal(t, data.EventFulfilled, status)

		_, err = claimEvent(eventID, n)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})

	t.Run("CountryBanned", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, mcoCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n)
		require.NoError(t, err)
		eventID, status := getEventFromList(respEvents, evtypes.TypeFreeWeekly)
		require.Equal(t, data.EventFulfilled, status)

		_, err = claimEvent(eventID, n)
		var apiErr jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
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
		_, respCode := requestWithBody(t, eventsEndpoint+"/"+freeWeeklyEventID, body, nullifier, true)
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

func claimEvent(id, nullifier string) (resp resources.EventResponse, err error) {
	body := claimEventBody(id)
	err = requestWithBody(eventsEndpoint+"/"+id, "PATCH", nullifier, body, &resp)
	return
}

func verifyPassport(nullifier, country string) (resp resources.PassportEventStateResponse, err error) {
	proof := baseProof
	proof.PubSignals[zk.Citizenship] = country
	body := verifyPassportBody(nullifier, proof)

	err = requestWithBody(balancesEndpoint+"/"+nullifier+"/verifypassport", "POST", nullifier, body, &resp)
	return
}

func getEvents(nullifier string) (resp resources.EventListResponse, err error) {
	query := url.Values{}
	query.Add("filter[nullifier]", nullifier)
	query.Add("filter[status]", "claimed,fulfilled,open")

	err = getRequest(eventsEndpoint, query, nullifier, &resp)
	return
}

func createBalance(nullifier, code string) (resp resources.BalanceResponse, err error) {
	body := createBalanceBody(nullifier, code)
	err = requestWithBody(balancesEndpoint, "POST", nullifier, body, &resp)
	return
}

func getBalance(nullifier string) (resp resources.BalanceResponse, err error) {
	query := url.Values{}
	query.Add("referral_codes", "true")
	query.Add("rank", "true")

	err = getRequest(balancesEndpoint+"/"+nullifier, query, nullifier, &resp)
	return
}

type editReferralsResponse struct {
	Ref       string `json:"referral"`
	UsageLeft uint64 `json:"usage_left"`
}

func editReferrals(nullifier string, count uint64) (resp editReferralsResponse, err error) {
	req := requests.EditReferralsRequest{Nullifier: nullifier, Count: count}

	err = requestWithBody(apiURL+"/private/referrals", "POST", "", req, &resp)

	return
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

func requestWithBody(endpoint, method, user string, body, result any) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	reqBody := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest(method, apiURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", method, err)
	}

	return doRequest(req, user, result)
}

func getRequest(endpoint string, query url.Values, user string, result any) error {
	req, err := http.NewRequest("GET", apiURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}
	req.URL.RawQuery = query.Encode()

	return doRequest(req, user, result)
}

func doRequest(req *http.Request, user string, result any) error {
	reqLog := fmt.Sprintf("%s /%s?%s", req.Method, req.URL.Path, req.URL.Query().Encode())

	if user != "" {
		req.Header.Set("nullifier", user)
	}

	resp, err := (&http.Client{Timeout: requestTimeout}).Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request (%s): %w", reqLog, err)
	}
	defer func() { _ = resp.Body.Close() }()

	log.Printf("Req: %s status=%d", reqLog, resp.StatusCode)

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
	default:
		return &jsonapi.ErrorObject{Status: strconv.Itoa(resp.StatusCode)}
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read resp body: %w", err)
	}

	err = json.Unmarshal(respBody, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

func nextN() string {
	currentNullifierIndex++
	return nullifiers[currentNullifierIndex-1]
}
