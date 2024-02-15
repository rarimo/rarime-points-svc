/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

// Primary event metadata in plain JSON. This is a template to be filled by `dynamic` when it's present.
type EventStaticMeta struct {
	Description string `json:"description"`
	// General event expiration date (UTC RFC3339)
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Event frequency, which means how often you can fulfill certain task and claim the reward.
	Frequency string `json:"frequency"`
	// Unique event code name
	Name string `json:"name"`
	// If true, the event will not be created with `open` status automatically when user creates the balance.
	NoAutoOpen bool `json:"no_auto_open"`
	// Reward amount in points
	Reward uint64 `json:"reward"`
	Title  string `json:"title"`
}
