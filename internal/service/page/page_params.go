package page

import (
	"math"

	val "github.com/go-ozzo/ozzo-validation/v4"
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
