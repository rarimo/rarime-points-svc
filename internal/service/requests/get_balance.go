package requests

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval/v4"
)

type GetBalance struct {
	Nullifier     string
	Rank          bool `url:"rank"`
	ReferralCodes bool `url:"referral_codes"`
}

func NewGetBalance(r *http.Request) (getBalance GetBalance, err error) {
	getBalance.Nullifier = strings.ToLower(chi.URLParam(r, "nullifier"))

	if err = urlval.Decode(r.URL.Query(), &getBalance); err != nil {
		err = newDecodeError("query", err)
		return
	}

	err = validation.Errors{"nullifier": validation.Validate(getBalance.Nullifier, validation.Required, validation.Match(nullifierRegexp))}.
		Filter()
	return
}
