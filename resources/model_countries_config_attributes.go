/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type CountriesConfigAttributes struct {
	// Country codes where users are eligible to claim
	Allowed []string `json:"allowed"`
	// Country codes where the limit of reservation was reached, making the claim not allowed
	LimitReached []string `json:"limit_reached"`
}
