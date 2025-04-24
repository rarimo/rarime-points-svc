/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type RootInclusionEventState struct {
	Key
	Attributes RootInclusionEventStateAttributes `json:"attributes"`
}
type RootInclusionEventStateResponse struct {
	Data     RootInclusionEventState `json:"data"`
	Included Included                `json:"included"`
}

type RootInclusionEventStateListResponse struct {
	Data     []RootInclusionEventState `json:"data"`
	Included Included                  `json:"included"`
	Links    *Links                    `json:"links"`
	Meta     json.RawMessage           `json:"meta,omitempty"`
}

func (r *RootInclusionEventStateListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *RootInclusionEventStateListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustRootInclusionEventState - returns RootInclusionEventState from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustRootInclusionEventState(key Key) *RootInclusionEventState {
	var rootInclusionEventState RootInclusionEventState
	if c.tryFindEntry(key, &rootInclusionEventState) {
		return &rootInclusionEventState
	}
	return nil
}
