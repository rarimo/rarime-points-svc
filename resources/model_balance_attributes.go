/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type BalanceAttributes struct {
	// Amount of points
	Amount uint64 `json:"amount"`
	// Unix timestamp of balance creation
	CreatedAt int32 `json:"created_at"`
	// Whether the user has scanned passport
	IsVerified bool `json:"is_verified"`
	// Rank of the user in the full leaderboard. Returned only for the single user.
	Rank *int `json:"rank,omitempty"`
	// Unix timestamp of the last points accruing
	UpdatedAt int32 `json:"updated_at"`
}
