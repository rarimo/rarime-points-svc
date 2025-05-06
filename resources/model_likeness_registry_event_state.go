/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type LikenessRegistryEventState struct {
	Key
	Attributes LikenessRegistryEventStateAttributes `json:"attributes"`
}
type LikenessRegistryEventStateResponse struct {
	Data     LikenessRegistryEventState `json:"data"`
	Included Included                   `json:"included"`
}

type LikenessRegistryEventStateListResponse struct {
	Data     []LikenessRegistryEventState `json:"data"`
	Included Included                     `json:"included"`
	Links    *Links                       `json:"links"`
	Meta     json.RawMessage              `json:"meta,omitempty"`
}

func (r *LikenessRegistryEventStateListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LikenessRegistryEventStateListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLikenessRegistryEventState - returns LikenessRegistryEventState from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLikenessRegistryEventState(key Key) *LikenessRegistryEventState {
	var likenessRegistryEventState LikenessRegistryEventState
	if c.tryFindEntry(key, &likenessRegistryEventState) {
		return &likenessRegistryEventState
	}
	return nil
}
