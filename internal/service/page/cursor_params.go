package page

import (
	"math"
	"net/http"
	"strconv"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/rarimo/rarime-points-svc/resources"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	pageParamLimit  = "page[limit]"
	pageParamCursor = "page[cursor]"
	pageParamOrder  = "page[order]"

	maxLimit uint64 = 100
)

// CursorParams is a wrapper around pgdb.CursorPageParams with useful validation and rendering methods
type CursorParams struct {
	pgdb.CursorPageParams
}

func (p *CursorParams) Validate() error {
	return val.Errors{
		pageParamLimit:  val.Validate(p.Limit, val.Max(maxLimit)),
		pageParamOrder:  val.Validate(p.Order, val.In(pgdb.OrderTypeAsc, pgdb.OrderTypeDesc)),
		pageParamCursor: val.Validate(p.Cursor, val.Max(uint64(math.MaxInt32))),
	}.Filter()
}

func (p *CursorParams) GetLinks(r *http.Request, last int32) *resources.Links {
	result := resources.Links{
		Self: p.getLink(r, p.Cursor),
	}
	if last != 0 {
		result.Next = p.getLink(r, uint64(last))
	}
	return &result
}

func (p *CursorParams) getLink(r *http.Request, cursor uint64) string {
	u := r.URL
	query := u.Query()
	query.Set(pageParamCursor, strconv.FormatUint(cursor, 10))
	query.Set(pageParamLimit, strconv.FormatUint(p.Limit, 10))
	query.Set(pageParamOrder, p.Order)
	u.RawQuery = query.Encode()
	return u.String()
}
