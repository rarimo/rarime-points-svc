/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type PassportEventState struct {
	Key
	Attributes PassportEventStateAttributes `json:"attributes"`
}
type PassportEventStateResponse struct {
	Data     PassportEventState `json:"data"`
	Included Included           `json:"included"`
}

type PassportEventStateListResponse struct {
	Data     []PassportEventState `json:"data"`
	Included Included             `json:"included"`
	Links    *Links               `json:"links"`
	Meta     json.RawMessage      `json:"meta,omitempty"`
}

func (r *PassportEventStateListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *PassportEventStateListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustPassportEventState - returns PassportEventState from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustPassportEventState(key Key) *PassportEventState {
	var passportEventState PassportEventState
	if c.tryFindEntry(key, &passportEventState) {
		return &passportEventState
	}
	return nil
}
