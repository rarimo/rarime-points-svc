package data

import (
	"database/sql/driver"
	"encoding/json"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Jsonb json.RawMessage

func (j *Jsonb) Value() (driver.Value, error) {
	if j == nil || len(*j) == 0 {
		return nil, nil
	}
	return pgdb.JSONValue(j)
}

// func (j *Jsonb) Scan(src interface{}) error {
// 	return pgdb.JSONScan(src, j)
// }

func (j *Jsonb) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j *Jsonb) Scan(src interface{}) error {
	var data []byte
	switch rawData := src.(type) {
	case []byte:
		data = rawData
	case string:
		data = []byte(rawData)
	case nil:
		data = []byte("null")
	default:
		return errors.New("Unexpected type for jsonb")
	}

	err := json.Unmarshal(data, j)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}

	return nil
}
