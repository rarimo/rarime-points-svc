/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Event struct {
	Key
	Attributes    EventAttributes     `json:"attributes"`
	Relationships *EventRelationships `json:"relationships,omitempty"`
}
type EventResponse struct {
	Data     Event    `json:"data"`
	Included Included `json:"included"`
}

type EventListResponse struct {
	Data     []Event         `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *EventListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *EventListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustEvent - returns Event from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustEvent(key Key) *Event {
	var event Event
	if c.tryFindEntry(key, &event) {
		return &event
	}
	return nil
}
