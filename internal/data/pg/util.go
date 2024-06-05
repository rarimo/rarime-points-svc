package pg

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

// AddToValue is a shorthand for expression to add/subtract to the current column value
func AddToValue(col string, add int64) squirrel.Sqlizer {
	return squirrel.Expr(fmt.Sprintf("%s + ?", col), add)
}
