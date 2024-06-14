/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type BalanceAttributes struct {
	// Amount of points
	Amount int64 `json:"amount"`
	// Unix timestamp of balance creation
	CreatedAt int32 `json:"created_at"`
	// Whether the user was not referred by anybody, but the balance with some events was reserved. It happens when the user fulfills some event before the balance creation.
	IsDisabled *bool `json:"is_disabled,omitempty"`
	// Whether the user has scanned passport. Returned only for the single user.
	IsVerified *bool `json:"is_verified,omitempty"`
	// The level indicates user permissions and features
	Level int `json:"level"`
	// Rank of the user in the full leaderboard. Returned only for the single user.
	Rank *int `json:"rank,omitempty"`
	// Referral codes. Returned only for the single user.
	ReferralCodes *[]ReferralCode `json:"referral_codes,omitempty"`
	// Unix timestamp of the last points accruing
	UpdatedAt int32 `json:"updated_at"`
}
