/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type BalanceAttributes struct {
	// Amount of points
	Amount int `json:"amount"`
	// Rank of the user in the full leaderboard. Returned only for the single user.
	Rank *int `json:"rank,omitempty"`
	// UTC time (RFC3339) of the last points accruing
	UpdatedAt time.Time `json:"updated_at"`
	// DID of the points owner
	UserDid string `json:"user_did"`
}
