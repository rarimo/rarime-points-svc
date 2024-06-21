/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "github.com/iden3/go-rapidsnark/types"

type VerifyPassportAttributes struct {
	// Unique identifier of the passport.
	AnonymousId string `json:"anonymous_id"`
	// ISO 3166-1 alpha-3 country code, must match the one provided in `proof`.
	Country string `json:"country"`
	// Query ZK passport verification proof. Required for endpoint `/v2/balances/{nullifier}/verifypassport`.
	Proof *types.ZKProof `json:"proof,omitempty"`
}
