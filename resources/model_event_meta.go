/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type EventMeta struct {
	// Some events require dynamic data, which can be filled into `static` template.
	Dynamic *json.RawMessage `json:"dynamic,omitempty"`
	// Primary event metadata in plain JSON. This is a template to be filled by `dynamic` when it's present.
	Static json.RawMessage `json:"static"`
}
