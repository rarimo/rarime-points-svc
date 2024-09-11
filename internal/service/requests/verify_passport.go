package requests

import (
	"encoding/json"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	zkptypes "github.com/iden3/go-rapidsnark/types"
	"github.com/rarimo/rarime-points-svc/resources"
	zk "github.com/rarimo/zkverifier-kit"
)

var (
	nullifierRegexp = regexp.MustCompile("^0x[0-9a-fA-F]{64}$")
	hex32bRegexp    = regexp.MustCompile("^[0-9a-f]{64}$")
	// endpoint is hardcoded to reuse handlers.VerifyPassport
	verifyPassportPathRegexp = regexp.MustCompile("^/integrations/rarime-points-svc/v1/public/balances/0x[0-9a-fA-F]{64}/verifypassport$")
	joinProgramPathRegexp    = regexp.MustCompile("^/integrations/rarime-points-svc/v1/public/balances/0x[0-9a-fA-F]{64}/join_program$")
)

func NewVerifyPassport(r *http.Request) (req resources.VerifyPassportRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, newDecodeError("body", err)
	}

	req.Data.ID = strings.ToLower(req.Data.ID)
	var (
		attr           = req.Data.Attributes
		provingCountry = attr.Country   // validate only when proof is provided
		proof          zkptypes.ZKProof // safe dereference
	)

	if attr.Proof != nil {
		proof = *attr.Proof
		provingCountry, err = ExtractCountry(proof)
		if err != nil {
			return req, err
		}
	}

	return req, val.Errors{
		"data/id": val.Validate(req.Data.ID,
			val.Required,
			val.In(strings.ToLower(chi.URLParam(r, "nullifier"))),
			val.Match(nullifierRegexp)),
		"data/type": val.Validate(req.Data.Type,
			val.Required,
			val.In(resources.VERIFY_PASSPORT)),
		"data/attributes/anonymous_id": val.Validate(attr.AnonymousId, val.Required, val.Match(hex32bRegexp)),
		"data/attributes/country":      val.Validate(attr.Country, val.Required, val.In(provingCountry), is.CountryCode3),
		"data/attributes/proof": val.Validate(attr.Proof,
			val.When(verifyPassportPathRegexp.MatchString(r.URL.Path), val.Required),
			val.When(joinProgramPathRegexp.MatchString(r.URL.Path), val.Nil)),
		"data/attributes/proof/proof":       val.Validate(proof.Proof, val.When(attr.Proof != nil, val.Required)),
		"data/attributes/proof/pub_signals": val.Validate(proof.PubSignals, val.When(attr.Proof != nil, val.Required, val.Length(23, 23))),
	}.Filter()
}

// ExtractCountry extracts country code from the proof, converting decimal UTF-8
// code to ISO 3166-1 alpha-3 code.
func ExtractCountry(proof zkptypes.ZKProof) (string, error) {
	if len(proof.PubSignals) <= zk.Indexes(zk.GlobalPassport)[zk.Citizenship] {
		return "", val.Errors{"country_code": val.ErrLengthTooShort}.Filter()
	}

	getter := zk.PubSignalGetter{Signals: proof.PubSignals, ProofType: zk.GlobalPassport}
	code := getter.Get(zk.Citizenship)

	b, ok := new(big.Int).SetString(code, 10)
	if !ok {
		b = new(big.Int)
	}

	code = string(b.Bytes())

	return code, val.Errors{"country_code": val.Validate(code, val.Required, is.CountryCode3)}.Filter()
}
