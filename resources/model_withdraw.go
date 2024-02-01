/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Withdraw struct {
	Key
	Attributes WithdrawAttributes `json:"attributes"`
}
type WithdrawRequest struct {
	Data     Withdraw `json:"data"`
	Included Included `json:"included"`
}

type WithdrawListRequest struct {
	Data     []Withdraw      `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *WithdrawListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *WithdrawListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustWithdraw - returns Withdraw from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustWithdraw(key Key) *Withdraw {
	var withdraw Withdraw
	if c.tryFindEntry(key, &withdraw) {
		return &withdraw
	}
	return nil
}
