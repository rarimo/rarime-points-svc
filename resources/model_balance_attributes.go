/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "time"

type BalanceAttributes struct {
	// Amount of points
	Amount int32 `json:"amount"`
	// UTC time (RFC3339) of the last points accruing
	UpdatedAt time.Time `json:"updated_at"`
	// DID of the points owner
	UserDid string `json:"user_did"`
}
