/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "github.com/iden3/go-rapidsnark/types"

type VerifyFaceAttributes struct {
	// Query ZK face verification proof. Required for endpoint `/v2/balances/verifyface`.
	Proof types.ZKProof `json:"proof"`
}
