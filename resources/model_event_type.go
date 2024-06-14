/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type EventType struct {
	Key
	Attributes EventStaticMeta `json:"attributes"`
}
type EventTypeResponse struct {
	Data     EventType `json:"data"`
	Included Included  `json:"included"`
}

type EventTypeListResponse struct {
	Data     []EventType     `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *EventTypeListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *EventTypeListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustEventType - returns EventType from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustEventType(key Key) *EventType {
	var eventType EventType
	if c.tryFindEntry(key, &eventType) {
		return &eventType
	}
	return nil
}
