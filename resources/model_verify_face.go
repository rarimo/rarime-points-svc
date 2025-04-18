/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type VerifyFace struct {
	Key
	Attributes VerifyFaceAttributes `json:"attributes"`
}
type VerifyFaceRequest struct {
	Data     VerifyFace `json:"data"`
	Included Included   `json:"included"`
}

type VerifyFaceListRequest struct {
	Data     []VerifyFace    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *VerifyFaceListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *VerifyFaceListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustVerifyFace - returns VerifyFace from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustVerifyFace(key Key) *VerifyFace {
	var verifyFace VerifyFace
	if c.tryFindEntry(key, &verifyFace) {
		return &verifyFace
	}
	return nil
}
