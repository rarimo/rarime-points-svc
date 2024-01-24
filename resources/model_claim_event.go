/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type ClaimEvent struct {
	Key
	Attributes ClaimEventAttributes `json:"attributes"`
}
type ClaimEventRequest struct {
	Data     ClaimEvent `json:"data"`
	Included Included   `json:"included"`
}

type ClaimEventListRequest struct {
	Data     []ClaimEvent    `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *ClaimEventListRequest) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *ClaimEventListRequest) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustClaimEvent - returns ClaimEvent from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustClaimEvent(key Key) *ClaimEvent {
	var claimEvent ClaimEvent
	if c.tryFindEntry(key, &claimEvent) {
		return &claimEvent
	}
	return nil
}
