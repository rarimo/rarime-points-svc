/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

// Primary event metadata in plain JSON. This is a template to be filled by `dynamic` when it's present.
type EventStaticMeta struct {
	Description string `json:"description"`
	// General event expiration date
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Unique event code name
	Name string `json:"name"`
	// Reward amount in points
	Reward int32  `json:"reward"`
	Title  string `json:"title"`
}
