/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type LikenessRegistry struct {
	Key
	Attributes LikenessRegistryAttributes `json:"attributes"`
}
type LikenessRegistryRequest struct {
	Data     LikenessRegistry `json:"data"`
	Included Included         `json:"included"`
}

type LikenessRegistryListRequest struct {
	Data     []LikenessRegistry `json:"data"`
	Included Included           `json:"included"`
	Links    *Links             `json:"links"`
	Meta     json.RawMessage    `json:"meta,omitempty"`
}

func (r *LikenessRegistryListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *LikenessRegistryListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustLikenessRegistry - returns LikenessRegistry from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustLikenessRegistry(key Key) *LikenessRegistry {
	var likenessRegistry LikenessRegistry
	if c.tryFindEntry(key, &likenessRegistry) {
		return &likenessRegistry
	}
	return nil
}
