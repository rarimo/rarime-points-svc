/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type EventAttributes struct {
	// UTC time (RFC3339) of event creation
	CreatedAt time.Time `json:"created_at"`
	Meta      EventMeta `json:"meta"`
	// See `filter[status]` parameter for explanation
	Status string `json:"status"`
}
