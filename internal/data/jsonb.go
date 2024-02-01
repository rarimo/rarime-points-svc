package data

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Jsonb json.RawMessage

func (j *Jsonb) Value() (driver.Value, error) {
	if j == nil || len(*j) == 0 {
		return nil, nil
	}
	return pgdb.JSONValue(j)
}

func (j *Jsonb) Scan(src interface{}) error {
	return pgdb.JSONScan(src, j)
}
