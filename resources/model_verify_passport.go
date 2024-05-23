/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type VerifyPassport struct {
	Key
	Attributes VerifyPassportAttributes `json:"attributes"`
}
type VerifyPassportRequest struct {
	Data     VerifyPassport `json:"data"`
	Included Included       `json:"included"`
}

type VerifyPassportListRequest struct {
	Data     []VerifyPassport `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
	Meta     json.RawMessage  `json:"meta,omitempty"`
}

func (r *VerifyPassportListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *VerifyPassportListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustVerifyPassport - returns VerifyPassport from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustVerifyPassport(key Key) *VerifyPassport {
	var verifyPassport VerifyPassport
	if c.tryFindEntry(key, &verifyPassport) {
		return &verifyPassport
	}
	return nil
}
