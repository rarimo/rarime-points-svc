/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RootInclusionVerify struct {
	Key
	Attributes RootInclusionVerifyAttributes `json:"attributes"`
}
type RootInclusionVerifyRequest struct {
	Data     RootInclusionVerify `json:"data"`
	Included Included            `json:"included"`
}

type RootInclusionVerifyListRequest struct {
	Data     []RootInclusionVerify `json:"data"`
	Included Included              `json:"included"`
	Links    *Links                `json:"links"`
	Meta     json.RawMessage       `json:"meta,omitempty"`
}

func (r *RootInclusionVerifyListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *RootInclusionVerifyListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustRootInclusionVerify - returns RootInclusionVerify from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRootInclusionVerify(key Key) *RootInclusionVerify {
	var rootInclusionVerify RootInclusionVerify
	if c.tryFindEntry(key, &rootInclusionVerify) {
		return &rootInclusionVerify
	}
	return nil
}
