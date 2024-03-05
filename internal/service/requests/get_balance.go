package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/urlval/v4"
)

type GetBalanceFilters struct {
	Rank          bool `url:"rank"`
	ReferralCodes bool `url:"referral_codes"`
}

type GetBalance struct {
	DID string
	GetBalanceFilters
}

func NewGetBalance(r *http.Request) (getBalance GetBalance, err error) {
	getBalance.DID = chi.URLParam(r, "did")

	if err = urlval.Decode(r.URL.Query(), &getBalance.GetBalanceFilters); err != nil {
		err = newDecodeError("query", err)
		return
	}

	err = validation.Errors{"did": validation.Validate(getBalance.DID, validation.Required)}.
		Filter()
	return
}
