/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import (
	"net/url"
	"time"
)

// Primary event metadata in plain JSON. This is a template to be filled by `dynamic` when it's present.
type EventStaticMeta struct {
	ActionUrl   *url.URL `json:"action_url,omitempty"`
	Description string   `json:"description"`
	// General event expiration date (UTC RFC3339)
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	// Event frequency, which means how often you can fulfill certain task and claim the reward.
	Frequency string   `json:"frequency"`
	Logo      *url.URL `json:"logo,omitempty"`
	// Unique event code name
	Name string `json:"name"`
	// Reward amount in points
	Reward           int64  `json:"reward"`
	ShortDescription string `json:"short_description"`
	// General event starting date (UTC RFC3339)
	StartsAt *time.Time `json:"starts_at,omitempty"`
	Title    string     `json:"title"`
}
