/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Withdrawal struct {
	Key
	Attributes WithdrawalAttributes `json:"attributes"`
}
type WithdrawalResponse struct {
	Data     Withdrawal `json:"data"`
	Included Included   `json:"included"`
}

type WithdrawalListResponse struct {
	Data     []Withdrawal    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *WithdrawalListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *WithdrawalListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustWithdrawal - returns Withdrawal from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustWithdrawal(key Key) *Withdrawal {
	var withdrawal Withdrawal
	if c.tryFindEntry(key, &withdrawal) {
		return &withdrawal
	}
	return nil
}
