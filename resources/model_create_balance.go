/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type CreateBalance struct {
	Key
	Attributes CreateBalanceAttributes `json:"attributes"`
}
type CreateBalanceRequest struct {
	Data     CreateBalance `json:"data"`
	Included Included      `json:"included"`
}

type CreateBalanceListRequest struct {
	Data     []CreateBalance `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *CreateBalanceListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CreateBalanceListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCreateBalance - returns CreateBalance from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCreateBalance(key Key) *CreateBalance {
	var createBalance CreateBalance
	if c.tryFindEntry(key, &createBalance) {
		return &createBalance
	}
	return nil
}
