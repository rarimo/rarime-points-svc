package requests

import (
	"encoding/json"
	"net/http"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
)

func NewLikenessRegistryVerifyRequest(r *http.Request) (req resources.LikenessRegistryRequest, err error) {
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = newDecodeError("body", err)
		return
	}

	var (
		proof = req.Data.Attributes.Proof
	)

	return req, val.Errors{
		"data/attributes/proof/proof":       val.Validate(proof.Proof, val.Required),
		"data/attributes/proof/pub_signals": val.Validate(proof.PubSignals, val.Required, val.Length(3, 3)),
	}.Filter()
}
