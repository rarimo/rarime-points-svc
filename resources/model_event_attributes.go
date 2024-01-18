/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type EventAttributes struct {
	// Unix milliseconds timestamp of event creation
	CreatedAt int32 `json:"created_at"`
	// If event has been already claimed
	IsClaimed bool `json:"is_claimed"`
	// Event metadata in JSON format. Configured along with `type`.
	Metadata json.RawMessage `json:"metadata"`
	// Event reward in points. Configured along with `type`.
	Reward int32 `json:"reward"`
	// Event type. Allowed event types are configured on back-end.
	Type string `json:"type"`
}
