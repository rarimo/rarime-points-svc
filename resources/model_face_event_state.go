/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type FaceEventState struct {
	Key
	Attributes FaceEventStateAttributes `json:"attributes"`
}
type FaceEventStateResponse struct {
	Data     FaceEventState `json:"data"`
	Included Included       `json:"included"`
}

type FaceEventStateListResponse struct {
	Data     []FaceEventState `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
	Meta     json.RawMessage  `json:"meta,omitempty"`
}

func (r *FaceEventStateListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *FaceEventStateListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustFaceEventState - returns FaceEventState from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustFaceEventState(key Key) *FaceEventState {
	var faceEventState FaceEventState
	if c.tryFindEntry(key, &faceEventState) {
		return &faceEventState
	}
	return nil
}
