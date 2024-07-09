package main_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/jsonapi"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/rarime-points-svc/internal/config"
	"github.com/rarimo/rarime-points-svc/internal/data"
	"github.com/rarimo/rarime-points-svc/internal/data/evtypes"
	"github.com/rarimo/rarime-points-svc/internal/data/pg"
	"github.com/rarimo/rarime-points-svc/internal/service/requests"
	"github.com/rarimo/rarime-points-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/kv"
)

const (
	requestTimeout        = time.Second           // use bigger on debug with breakpoints to prevent fails
	defaultConfigFile     = "config-testing.yaml" // run service with this config for consistency with tests
	defaultReferralsCount = 5

	ukrCode = "UKR"
	usaCode = "USA"
	gbrCode = "GBR"
	deuCode = "DEU"
	canCode = "CAN"
	fraCode = "FRA"
	indCode = "IND"
	mcoCode = "MCO"
	belCode = "BEL"
	mngCode = "MNG"

	genesisBalance = "0x0000000000000000000000000000000000000000000000000000000000000000"
	rarimoAddress  = "rarimo1h2077nfkksek386y8ks5m2wgd60wl3035n8gv0"

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
}

func setUp() {
	if os.Getenv(kv.EnvViperConfigFile) == "" {
		err := os.Setenv(kv.EnvViperConfigFile, defaultConfigFile)
		if err != nil {
			panic(fmt.Errorf("failed to set env: %w", err))
		}
	}

	var err error
	apiURL, err = getApiURL()
	if err != nil {
		panic(fmt.Errorf("failed to get Api URL: %w", err))
	}

	globalCfg = config.New(kv.MustFromEnv())
	initGenesisRef()
	// let's not introduce counting function just for test
	balances, err := pg.NewBalances(globalCfg.DB()).Select()
	if err != nil {
		panic(fmt.Errorf("failed to select balances: %w", err))
	}

	// to prevent repeating cleanups, more balances are created
	currentNullifierIndex = len(balances)
	nullifiers = make([]string, 100)
	for i := range nullifiers {
		hash := sha256.Sum256([]byte{byte(i + len(balances))})
		nullifiers[i] = hexutil.Encode(hash[:])
	}
}

func getApiURL() (string, error) {
	var cfg struct {
		Addr string `fig:"addr,required"`
	}

	err := figure.Out(&cfg).From(kv.MustGetStringMap(kv.MustFromEnv(), "listener")).Please()
	if err != nil {
		return "", fmt.Errorf("failed to figure out listener from service config: %w", err)
	}

	apiURL := fmt.Sprintf("http://localhost%s/integrations/rarime-points-svc/v1/", cfg.Addr)
	return apiURL, nil
}

func initGenesisRef() {
	gen, err := pg.NewReferrals(globalCfg.DB()).
		FilterConsumed().
		FilterByNullifier(genesisBalance).
		Select()
	if err != nil {
		panic(fmt.Errorf("failed to get genesis balance: %w", err))
	}
	if len(gen) > 1 {
		panic(fmt.Errorf("%d genesis referral codes found", len(gen)))
	}

	if len(gen) == 0 || gen[0].UsageLeft < 20 { // approximate amount to run tests
		refs, err := editReferrals(genesisBalance, 10000)
		if err != nil {
			panic(fmt.Errorf("failed to edit referrals: %w", err))
		}
		genesisCode = refs.Ref
		return
	}

	genesisCode = gen[0].ID
}

func TestCreateBalance(t *testing.T) {
	var (
		nullifierShared = nextN()
		otRefCode       string
	)

	t.Run("BalanceGenesisCode", func(t *testing.T) {
		resp := createAndValidateBalance(t, nullifierShared, genesisCode)
		otRefCode = (*resp.Data.Attributes.ReferralCodes)[0].Id
	})

	t.Run("BalanceOneTimeCode", func(t *testing.T) {
		createAndValidateBalance(t, nextN(), otRefCode)
	})

	t.Run("SameBalanceConflict", func(t *testing.T) {
		_, err := createBalance(nullifierShared, genesisCode)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "409", apiErr.Status)
	})

	t.Run("NullifierUnauthorized", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		body := createBalanceBody(n1, genesisCode)
		err := requestWithBody(balancesEndpoint, "POST", n2, "", body, nil)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "401", apiErr.Status)
	})

	t.Run("ConsumedCode", func(t *testing.T) {
		_, err := createBalance(nextN(), otRefCode)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "404", apiErr.Status)
	})

	t.Run("IncorrectCode", func(t *testing.T) {
		_, err := createBalance(nextN(), "invalid")
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "404", apiErr.Status)
	})
}

func createAndValidateBalance(t *testing.T, nullifier, code string) resources.BalanceResponse {
	t.Helper()

	resp, err := createBalance(nullifier, code)
	require.NoError(t, err)
	require.Equal(t, nullifier, resp.Data.ID)

	attr := resp.Data.Attributes

	require.NotNil(t, attr.IsDisabled)
	require.NotNil(t, attr.IsVerified)
	require.NotNil(t, attr.ReferralCodes)
	require.NotEmpty(t, *attr.ReferralCodes)

	assert.Equal(t, int64(0), attr.Amount)
	assert.False(t, *attr.IsDisabled)
	assert.False(t, *attr.IsVerified)
	assert.Equal(t, 1, attr.Level)
	assert.NotNil(t, attr.Rank)

	rc := (*attr.ReferralCodes)[0]
	assert.NotEmpty(t, rc.Id)
	assert.Equal(t, data.StatusActive, rc.Status)
	return resp
}

func TestVerifyPassport(t *testing.T) {
	var (
		referrer = nextN()
		referee  = nextN()
		balance1 = createAndValidateBalance(t, referrer, genesisCode)
	)
	createAndValidateBalance(t, referee, (*balance1.Data.Attributes.ReferralCodes)[0].Id)

	var countriesResp resources.CountriesConfigResponse
	err := getRequest("public/countries_config", nil, "", &countriesResp)
	require.NoError(t, err)

	countriesList := countriesResp.Data.Attributes.Countries

	var ukr, usa bool
	// ensure the same behaviour whitelisted and banned countries
	for _, c := range countriesList {
		if c.Code == ukrCode {
			ukr = true
			require.True(t, c.ReserveAllowed)
			require.True(t, c.WithdrawalAllowed)
			continue
		}
		if c.Code == usaCode {
			usa = true
			require.False(t, c.ReserveAllowed)
			require.False(t, c.WithdrawalAllowed)
		}
	}
	require.False(t, !ukr || !usa)

	// passport verification should lead to referral event appearance and claimed passport event
	t.Run("VerifyPassport", func(t *testing.T) {
		resp, err := verifyPassport(referee, ukrCode)
		require.NoError(t, err)
		assert.True(t, resp.Data.Attributes.Claimed)
		getAndValidateBalance(t, referee, true)
	})

	t.Run("VerifyPassportSecondTime", func(t *testing.T) {
		_, err = verifyPassport(referee, ukrCode)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "429", apiErr.Status)
		getAndValidateBalance(t, referee, true)
	})

	t.Run("IncorrectCountryCode", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err = verifyPassport(n, "6974819")
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "400", apiErr.Status)
		getAndValidateBalance(t, n, false)
	})
}

func getAndValidateBalance(t *testing.T, nullifier string, isVerified bool) resources.BalanceResponse {
	resp, err := getBalance(nullifier)
	require.NoError(t, err)

	attr := resp.Data.Attributes
	require.NotNil(t, attr.IsDisabled)
	require.NotNil(t, attr.IsVerified)
	assert.False(t, *attr.IsDisabled)
	assert.Equal(t, isVerified, *attr.IsVerified)

	assert.NotNil(t, attr.Rank)
	assert.NotNil(t, attr.ReferralCodes)
	assert.NotEmpty(t, *attr.ReferralCodes)

	return resp
}

func TestEventsAutoClaim(t *testing.T) {
	t.Run("PassportScanAutoclaim", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respVerifyStatus, err := verifyPassport(n, canCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respBalance, err := getBalance(n)
		require.NoError(t, err)
		require.Equal(t, 2, respBalance.Data.Attributes.Level)
		require.Equal(t, int64(5), respBalance.Data.Attributes.Amount)
		require.NotNil(t, respBalance.Data.Attributes.ReferralCodes)
		require.Equal(t, 10, len(*respBalance.Data.Attributes.ReferralCodes))
	})

	// this test depend on previous test
	t.Run("PassportScanLimitReached", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		respVerifyStatus, err := verifyPassport(n, canCode)
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

		respVerifyStatus, err := verifyPassport(n, usaCode)
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

		respEvents, err := getEvents(n1, evtypes.TypeReferralSpecific)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventClaimed), respEvents.Data[0].Attributes.Status)
	})

	// User can have a lot unclaimed fulfilled referral specific events if user not scan passport
	t.Run("ReferralSpecificAutoclaimMany", func(t *testing.T) {
		n1, n2, n3 := nextN(), nextN(), nextN()
		respBalance, err := createBalance(n1, genesisCode)
		require.NoError(t, err)
		require.NotNil(t, respBalance.Data.Attributes.ReferralCodes)
		require.GreaterOrEqual(t, len(*respBalance.Data.Attributes.ReferralCodes), 2)

		_, err = createBalance(n2, (*respBalance.Data.Attributes.ReferralCodes)[0].Id)
		require.NoError(t, err)
		_, err = verifyPassport(n2, ukrCode)
		require.NoError(t, err)

		_, err = createBalance(n3, (*respBalance.Data.Attributes.ReferralCodes)[1].Id)
		require.NoError(t, err)
		_, err = verifyPassport(n3, ukrCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n1, evtypes.TypeReferralSpecific)
		require.NoError(t, err)
		require.Equal(t, 2, len(respEvents.Data))
		fulfilledEventCount := 0
		for _, event := range respEvents.Data {
			if event.Attributes.Status == string(data.EventFulfilled) {
				fulfilledEventCount++
			}
		}
		require.Equal(t, 2, fulfilledEventCount)

		respVerifyStatus, err := verifyPassport(n1, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respEvents, err = getEvents(n1, evtypes.TypeReferralSpecific)
		require.NoError(t, err)
		require.Equal(t, 2, len(respEvents.Data))
		claimedEventCount := 0
		for _, event := range respEvents.Data {
			if event.Attributes.Status == string(data.EventClaimed) {
				claimedEventCount++
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

		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventFulfilled), respEvents.Data[0].Attributes.Status)

		_, err = claimEvent(respEvents.Data[0].ID, n)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})

	t.Run("IncorrectEventID", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, ukrCode)
		require.NoError(t, err)

		_, err = claimEvent("e174d6e2-0c81-4771-99a1-8447532143b8", n)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "404", apiErr.Status)
	})

	t.Run("EventClaim", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, fraCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventFulfilled), respEvents.Data[0].Attributes.Status)

		respEvent, err := claimEvent(respEvents.Data[0].ID, n)
		require.NoError(t, err)
		require.Equal(t, string(data.EventClaimed), respEvent.Data.Attributes.Status)
	})

	// this test depend on previous test
	t.Run("ReserveLimitReached", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, fraCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventFulfilled), respEvents.Data[0].Attributes.Status)

		_, err = claimEvent(respEvents.Data[0].ID, n)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})

	t.Run("ReserveDisallowed", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, indCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventFulfilled), respEvents.Data[0].Attributes.Status)

		_, err = claimEvent(respEvents.Data[0].ID, n)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})

	t.Run("CountryBanned", func(t *testing.T) {
		n := nextN()
		_, err := createBalance(n, genesisCode)
		require.NoError(t, err)

		_, err = verifyPassport(n, mcoCode)
		require.NoError(t, err)

		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventFulfilled), respEvents.Data[0].Attributes.Status)

		_, err = claimEvent(respEvents.Data[0].ID, n)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})
}

func TestLevels(t *testing.T) {
	var (
		nullifier = nextN()

		evTypePassport = globalCfg.EventTypes().Get(evtypes.TypePassportScan)
		evTypeWeekly   = globalCfg.EventTypes().Get(evtypes.TypeFreeWeekly)

		lvl1Cfg = globalCfg.Levels()[1]
		lvl2Cfg = globalCfg.Levels()[2]
		lvl3Cfg = globalCfg.Levels()[3]

		amountClaim1  = evTypePassport.Reward
		amountClaim2  = evTypePassport.Reward + evTypeWeekly.Reward
		lvl2Referrals = lvl1Cfg.Referrals + lvl2Cfg.Referrals
		lvl3Referrals = lvl2Referrals + lvl3Cfg.Referrals
	)
	require.NotNil(t, evTypePassport)
	require.NotNil(t, evTypeWeekly)
	// ensure that levels are set
	require.Equal(t, 1, lvl1Cfg.Level)
	require.Equal(t, 2, lvl2Cfg.Level)
	require.Equal(t, 3, lvl3Cfg.Level)
	// rewards must be equal to level threshold in order to upgrade level for each of 2 claimed events
	require.Equal(t, amountClaim1, int64(lvl2Cfg.Threshold))
	require.Equal(t, amountClaim2, int64(lvl3Cfg.Threshold))
	require.False(t, evTypeWeekly.AutoClaim)

	createAndValidateBalance(t, nullifier, genesisCode)
	passportResp, err := verifyPassport(nullifier, ukrCode)
	require.NoError(t, err)
	assert.True(t, passportResp.Data.Attributes.Claimed)

	status := data.EventClaimed
	if !evTypePassport.AutoClaim {
		status = data.EventFulfilled
	}
	eventID := getAndValidateSingleEvent(t, nullifier, evtypes.TypePassportScan, status)

	if !evTypePassport.AutoClaim {
		claimEventAndValidate(t, eventID, nullifier, amountClaim1)
	}

	balance := getAndValidateBalance(t, nullifier, true)
	balanceAttr := balance.Data.Attributes
	assert.Equal(t, 2, balanceAttr.Level)
	assert.Equal(t, amountClaim1, balanceAttr.Amount)

	refCodes := balanceAttr.ReferralCodes
	require.NotNil(t, refCodes)
	assert.Equal(t, lvl2Referrals, len(*refCodes))

	eventID = getAndValidateSingleEvent(t, nullifier, evtypes.TypeFreeWeekly, data.EventFulfilled)
	claimEventAndValidate(t, eventID, nullifier, 1)

	balance = getAndValidateBalance(t, nullifier, true)
	balanceAttr = balance.Data.Attributes
	assert.Equal(t, 3, balanceAttr.Level)
	assert.Equal(t, amountClaim2, balanceAttr.Amount)

	refCodes = balanceAttr.ReferralCodes
	require.NotNil(t, refCodes)
	assert.Equal(t, lvl3Referrals, len(*refCodes))
}

func getAndValidateSingleEvent(t *testing.T, nullifier, evType string, status data.EventStatus) string {
	resp, err := getEvents(nullifier, evType)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)

	event := resp.Data[0]
	attr := event.Attributes

	require.NotEmpty(t, event.ID)
	assert.Equal(t, evType, attr.Meta.Static.Name)
	assert.Equal(t, string(status), attr.Status)
	return event.ID
}

func claimEventAndValidate(t *testing.T, id, nullifier string, reward int64) {
	resp, err := claimEvent(id, nullifier)
	require.NoError(t, err)
	attr := resp.Data.Attributes
	assert.Equal(t, string(data.EventClaimed), attr.Status)
	require.NotNil(t, attr.PointsAmount)
	assert.Equal(t, reward, *attr.PointsAmount)
}

// test only default config because main logic already tested in another tests (autoclaim, claim, verifypassport)
func TestCountryPoolsDefault(t *testing.T) {
	n := nextN()
	createAndValidateBalance(t, n, genesisCode)

	t.Run("DefaultUnderLimit", func(t *testing.T) {
		resp, err := verifyPassport(n, deuCode)
		require.NoError(t, err)
		assert.True(t, resp.Data.Attributes.Claimed)
		getAndValidateBalance(t, n, true)

	})

	t.Run("DefaultOverLimit", func(t *testing.T) {
		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		require.Equal(t, string(data.EventFulfilled), respEvents.Data[0].Attributes.Status)

		_, err = claimEvent(respEvents.Data[0].ID, n)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "403", apiErr.Status)
	})
}

func TestReferralCodeStatuses(t *testing.T) {
	t.Run("ActiveCode", func(t *testing.T) {
		n := nextN()
		respBalance := createAndValidateBalance(t, n, genesisCode)
		require.Equal(t, 5, len(*respBalance.Data.Attributes.ReferralCodes))
		for _, v := range *respBalance.Data.Attributes.ReferralCodes {
			require.Equal(t, data.StatusActive, v.Status)
		}
	})

	t.Run("BannedCode", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		respBalance := createAndValidateBalance(t, n1, genesisCode)
		respVerifyStatus, err := verifyPassport(n1, usaCode)
		require.NoError(t, err)
		require.False(t, respVerifyStatus.Data.Attributes.Claimed)

		refCode := (*respBalance.Data.Attributes.ReferralCodes)[0].Id
		createAndValidateBalance(t, n2, refCode)

		respBalance = getAndValidateBalance(t, n1, true)
		for _, v := range *respBalance.Data.Attributes.ReferralCodes {
			if v.Id == refCode && v.Status == data.StatusBanned {
				return
			}
		}
		t.Fatal("Banned referral code absent")
	})

	t.Run("LimitedCode", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		respBalance := createAndValidateBalance(t, n1, genesisCode)
		respVerifyStatus, err := verifyPassport(n1, gbrCode)
		require.NoError(t, err)
		require.False(t, respVerifyStatus.Data.Attributes.Claimed)

		refCode := (*respBalance.Data.Attributes.ReferralCodes)[0].Id
		createAndValidateBalance(t, n2, refCode)

		respBalance = getAndValidateBalance(t, n1, true)
		for _, v := range *respBalance.Data.Attributes.ReferralCodes {
			if v.Id == refCode && v.Status == data.StatusLimited {
				return
			}
		}
		t.Fatal("Limited referral code absent")
	})

	t.Run("AwaitingCode", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		respBalance := createAndValidateBalance(t, n1, genesisCode)

		refCode := (*respBalance.Data.Attributes.ReferralCodes)[0].Id
		createAndValidateBalance(t, n2, refCode)
		respVerifyStatus, err := verifyPassport(n2, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respBalance = getAndValidateBalance(t, n1, false)
		for _, v := range *respBalance.Data.Attributes.ReferralCodes {
			if v.Id == refCode && v.Status == data.StatusAwaiting {
				return
			}
		}
		t.Fatal("Awaiting referral code absent")
	})

	t.Run("RewardedCode", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		respBalance := createAndValidateBalance(t, n1, genesisCode)
		respVerifyStatus, err := verifyPassport(n1, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		refCode := (*respBalance.Data.Attributes.ReferralCodes)[0].Id
		createAndValidateBalance(t, n2, refCode)
		respVerifyStatus, err = verifyPassport(n2, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		respBalance = getAndValidateBalance(t, n1, true)
		for _, v := range *respBalance.Data.Attributes.ReferralCodes {
			if v.Id == refCode && v.Status == data.StatusRewarded {
				return
			}
		}
		t.Fatal("Rewarded referral code absent")
	})

	t.Run("ConsumedCode", func(t *testing.T) {
		n1, n2 := nextN(), nextN()
		respBalance := createAndValidateBalance(t, n1, genesisCode)
		respVerifyStatus, err := verifyPassport(n1, ukrCode)
		require.NoError(t, err)
		require.True(t, respVerifyStatus.Data.Attributes.Claimed)

		refCode := (*respBalance.Data.Attributes.ReferralCodes)[0].Id
		createAndValidateBalance(t, n2, refCode)

		respBalance = getAndValidateBalance(t, n1, true)
		for _, v := range *respBalance.Data.Attributes.ReferralCodes {
			if v.Id == refCode && v.Status == data.StatusConsumed {
				return
			}
		}
		t.Fatal("Consumed referral code absent")
	})
}

func TestWithdrawals(t *testing.T) {
	t.Run("WithoutPassport", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err := withdraw(n, ukrCode, 10)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "400", apiErr.Status)
	})

	t.Run("BalanceAbsent", func(t *testing.T) {
		n := nextN()
		_, err := withdraw(n, ukrCode, 10)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "404", apiErr.Status)
	})

	t.Run("IncorrectCountryCode", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err := withdraw(n, "6974819", 10)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "500", apiErr.Status)
	})

	t.Run("CountryMismatched", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err := verifyPassport(n, ukrCode)
		require.NoError(t, err)
		_, err = withdraw(n, fraCode, 1)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "400", apiErr.Status)
	})

	t.Run("InsufficientBalance", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err := verifyPassport(n, ukrCode)
		require.NoError(t, err)
		_, err = withdraw(n, ukrCode, 10)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "400", apiErr.Status)
	})

	t.Run("WithdrawNotAllowed", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err := verifyPassport(n, belCode)
		require.NoError(t, err)
		respEvents, err := getEvents(n, evtypes.TypeFreeWeekly)
		require.NoError(t, err)
		require.Equal(t, 1, len(respEvents.Data))
		claimEventAndValidate(t, respEvents.Data[0].ID, n, 1)
		_, err = withdraw(n, belCode, 4)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "400", apiErr.Status)
	})

	t.Run("InsufficientLevelToWithdraw", func(t *testing.T) {
		n := nextN()
		createAndValidateBalance(t, n, genesisCode)
		_, err := verifyPassport(n, ukrCode)
		require.NoError(t, err)
		_, err = withdraw(n, ukrCode, 10)
		var apiErr *jsonapi.ErrorObject
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, "400", apiErr.Status)
	})
}

func withdraw(nullifier, country string, amount int64) (resp resources.WithdrawalResponse, err error) {
	proof := baseProof
	proof.PubSignals[zk.Citizenship] = code2dec(country)
	body := withdrawBody(nullifier, proof, amount)
	err = requestWithBody(balancesEndpoint+"/"+nullifier+"/withdrawals", "POST", nullifier, "", body, &resp)
	return

}

func withdrawBody(nullifier string, proof zkptypes.ZKProof, amount int64) resources.WithdrawRequest {
	return resources.WithdrawRequest{
		Data: resources.Withdraw{
			Key: resources.Key{
				ID:   nullifier,
				Type: resources.WITHDRAW,
			},
			Attributes: resources.WithdrawAttributes{
				Address: rarimoAddress,
				Proof:   proof,
				Amount:  amount,
			},
		},
	}
}

func claimEvent(id, nullifier string) (resp resources.EventResponse, err error) {
	body := claimEventBody(id)
	err = requestWithBody(eventsEndpoint+"/"+id, "PATCH", nullifier, "", body, &resp)
	return
}

func getEvents(nullifier string, types ...string) (resp resources.EventListResponse, err error) {
	query := url.Values{}
	query.Add("filter[nullifier]", nullifier)
	query.Add("page[limit]", "100")
	if len(types) > 0 {
		query.Add("filter[meta.static.name]", strings.Join(types, ","))
	}

	err = getRequest(eventsEndpoint, query, nullifier, &resp)
	return
}

func createBalance(nullifier, code string) (resp resources.BalanceResponse, err error) {
	body := createBalanceBody(nullifier, code)
	err = requestWithBody(balancesEndpoint, "POST", nullifier, "", body, &resp)
	return
}

func getBalance(nullifier string) (resp resources.BalanceResponse, err error) {
	query := url.Values{}
	query.Add("referral_codes", "true")
	query.Add("rank", "true")

	err = getRequest(balancesEndpoint+"/"+nullifier, query, nullifier, &resp)
	return
}

func verifyPassport(nullifier, country string) (resp resources.PassportEventStateResponse, err error) {
	proof := baseProof
	proof.PubSignals[zk.Citizenship] = code2dec(country)
	body := verifyPassportBody(nullifier, country, nullifier[2:], &proof)
	err = requestWithBody(balancesEndpoint+"/"+nullifier+"/verifypassport", "POST", nullifier, signature(globalCfg.Countries().VerificationKey, nullifier, country, nullifier[2:]), body, &resp)
	return
}

type editReferralsResponse struct {
	Ref       string `json:"referral"`
	UsageLeft uint64 `json:"usage_left"`
}

func editReferrals(nullifier string, count uint64) (resp editReferralsResponse, err error) {
	req := requests.EditReferralsRequest{Nullifier: nullifier, Count: count}
	err = requestWithBody("private/referrals", "POST", "", "", req, &resp)
	return
}

func verifyPassportBody(nullifier, country, anonymousID string, proof *zkptypes.ZKProof) resources.VerifyPassportRequest {
	return resources.VerifyPassportRequest{
		Data: resources.VerifyPassport{
			Key: resources.Key{
				ID:   nullifier,
				Type: resources.VERIFY_PASSPORT,
			},
			Attributes: resources.VerifyPassportAttributes{
				AnonymousId: anonymousID,
				Country:     country,
				Proof:       proof,
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

func requestWithBody(endpoint, method, user, signature string, body, result any) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal body: %w", err)
	}

	reqBody := bytes.NewReader(bodyJSON)
	req, err := http.NewRequest(method, apiURL+endpoint, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", method, err)
	}

	return doRequest(req, user, signature, result)
}

func getRequest(endpoint string, query url.Values, user string, result any) error {
	req, err := http.NewRequest("GET", apiURL+endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}
	req.URL.RawQuery = query.Encode()

	return doRequest(req, user, "", result)
}

func doRequest(req *http.Request, user, signature string, result any) error {
	reqLog := fmt.Sprintf("%s /%s?%s", req.Method, req.URL.Path, req.URL.Query().Encode())

	if user != "" {
		req.Header.Set("nullifier", user)
	}

	if signature != "" {
		req.Header.Set("Signature", signature)
	}

	resp, err := (&http.Client{Timeout: requestTimeout}).Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request (%s): %w", reqLog, err)
	}
	defer func() { resp.Body.Close() }()

	log.Printf("Req: %s status=%d", reqLog, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read resp body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
	default:
		return &jsonapi.ErrorObject{Status: strconv.Itoa(resp.StatusCode), Title: string(respBody)}
	}

	if result == nil {
		return nil
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

func signature(key []byte, nullifier, country, anonymousID string) string {
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

func code2dec(in string) (out string) {
	return new(big.Int).SetBytes([]byte(in)).String()
}

func dec2code(in string) (out string) {
	b, ok := new(big.Int).SetString(in, 10)
	if !ok {
		b = new(big.Int)
	}

	return string(b.Bytes())
}
