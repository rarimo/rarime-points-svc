package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"gitlab.com/distributed_lab/kit/pgdb"
)

type Jsonb json.RawMessage

func (j *Jsonb) Value() (driver.Value, error) {
	if j == nil || len(*j) == 0 {
		return nil, nil
	}
	return pgdb.JSONValue(j)
}

func (j *Jsonb) UnmarshalJSON(data []byte) error {
	if j == nil {
		return fmt.Errorf("json.RawMessage: UnmarshalJSON on nil pointer")
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
		return fmt.Errorf("unexpected type for jsonb: %T", src)
	}

	err := json.Unmarshal(data, j)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}
