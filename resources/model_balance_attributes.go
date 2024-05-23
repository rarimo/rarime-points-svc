/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type BalanceAttributes struct {
	// Referral codes which can be used to build a referral link and send it to friends. Returned only for the single user.
	ActiveReferralCodes *[]string `json:"active_referral_codes,omitempty"`
	// Amount of points
	Amount int64 `json:"amount"`
	// Referral codes used by invited users. Returned only for the single user.
	ConsumedReferralCodes *[]string `json:"consumed_referral_codes,omitempty"`
	// Unix timestamp of balance creation
	CreatedAt int32 `json:"created_at"`
	// Whether the user was not referred by anybody, but the balance with some events was reserved. It happens when the user fulfills some event before the balance creation.
	IsDisabled bool `json:"is_disabled"`
	// The level indicates how many possibilities the user has
	Level int `json:"level"`
	// Rank of the user in the full leaderboard. Returned only for the single user.
	Rank *int `json:"rank,omitempty"`
	// Unix timestamp of the last points accruing
	UpdatedAt int32 `json:"updated_at"`
}
