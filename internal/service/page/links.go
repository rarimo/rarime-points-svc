package page

import (
	"net/http"
	"strconv"

	"github.com/rarimo/points-svc/resources"
)

func (p *CursorParams) GetCursorLinks(r *http.Request, last string) *resources.Links {
	result := resources.Links{
		Self: p.getCursorLink(r, p.Cursor, p.Order),
	}
	if last != "" {
		lastI, _ := strconv.ParseUint(last, 10, 64)
		result.Next = p.getCursorLink(r, lastI, p.Order)
	}
	return &result
}

func (p *CursorParams) getCursorLink(r *http.Request, cursor uint64, order string) string {
	u := r.URL
	query := u.Query()
	query.Set(pageParamCursor, strconv.FormatUint(cursor, 10))
	query.Set(pageParamLimit, strconv.FormatUint(p.Limit, 10))
	query.Set(pageParamOrder, order)
	u.RawQuery = query.Encode()
	return u.String()
}
