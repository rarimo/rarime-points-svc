/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type EventAttributes struct {
	// UTC time (RFC3339) of event creation
	CreatedAt time.Time `json:"created_at"`
	Meta      EventMeta `json:"meta"`
	// How many points were accrued. Required only for `claimed` events. This is necessary, as the reward might change over time, while the certain balance should be left intact.
	PointsAmount *int32 `json:"points_amount,omitempty"`
	// See `filter[status]` parameter for explanation
	Status string `json:"status"`
}
