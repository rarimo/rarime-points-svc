/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type ReferralCode struct {
	Key
	Attributes *map[string]interface{} `json:"attributes,omitempty"`
}
type ReferralCodeRequest struct {
	Data     ReferralCode `json:"data"`
	Included Included     `json:"included"`
}

type ReferralCodeListRequest struct {
	Data     []ReferralCode  `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ReferralCodeListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ReferralCodeListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustReferralCode - returns ReferralCode from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustReferralCode(key Key) *ReferralCode {
	var referralCode ReferralCode
	if c.tryFindEntry(key, &referralCode) {
		return &referralCode
	}
	return nil
}
