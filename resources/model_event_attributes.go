/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type EventAttributes struct {
	// Unix timestamp of event creation
	CreatedAt int32 `json:"created_at"`
	// Whether this event may become expired.
	HasExpiration bool      `json:"has_expiration"`
	Meta          EventMeta `json:"meta"`
	// How many points were accrued. Required only for `claimed` events. This is necessary, as the reward might change over time, while the certain balance should be left intact.
	PointsAmount *int64 `json:"points_amount,omitempty"`
	// See `filter[status]` parameter for explanation
	Status string `json:"status"`
	// Unix timestamp of the event status change
	UpdatedAt int32 `json:"updated_at"`
}
