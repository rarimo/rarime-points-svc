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
	IsLeaderboard bool // deny ascending order for leaderboard
}

func (p *CursorParams) Validate() error {
	var orderRule val.Rule = val.In(pgdb.OrderTypeAsc, pgdb.OrderTypeDesc)
	if p.IsLeaderboard {
		orderRule = val.Empty
	}

	return val.Errors{
		pageParamCursor: val.Validate(p.Cursor, val.Max(uint64(math.MaxInt32))),
		pageParamLimit:  val.Validate(p.Limit, val.Max(maxLimit)),
		pageParamOrder:  val.Validate(p.Order, orderRule),
	}.Filter()
}
