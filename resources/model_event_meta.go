/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type EventMeta struct {
	// Some events require dynamic data, which can be filled into `static` template.
	Dynamic *json.RawMessage `json:"dynamic,omitempty"`
	Static  EventStaticMeta  `json:"static"`
}
