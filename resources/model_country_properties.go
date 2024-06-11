/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CountryProperties struct {
	// ISO 3166-1 alpha-3 country code
	Code string `json:"code"`
	// Whether the users of country are allowed to reserve (claim) tokens
	ReserveAllowed bool `json:"reserve_allowed"`
	// Whether the users of country are allowed to withdraw tokens
	WithdrawalAllowed bool `json:"withdrawal_allowed"`
}
