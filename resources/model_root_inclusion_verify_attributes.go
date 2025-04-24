/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "github.com/iden3/go-rapidsnark/types"

type RootInclusionVerifyAttributes struct {
	// Query ZK root inclusion verification proof. Required for endpoint `/v2/balances/root_verify`.
	Proof types.ZKProof `json:"proof"`
}
