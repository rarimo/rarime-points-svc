/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type JoinProgram struct {
	Key
	Attributes JoinProgramAttributes `json:"attributes"`
}
type JoinProgramRequest struct {
	Data     JoinProgram `json:"data"`
	Included Included    `json:"included"`
}

type JoinProgramListRequest struct {
	Data     []JoinProgram   `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *JoinProgramListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *JoinProgramListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustJoinProgram - returns JoinProgram from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustJoinProgram(key Key) *JoinProgram {
	var joinProgram JoinProgram
	if c.tryFindEntry(key, &joinProgram) {
		return &joinProgram
	}
	return nil
}
