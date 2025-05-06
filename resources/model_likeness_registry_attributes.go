/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "github.com/iden3/go-rapidsnark/types"

type LikenessRegistryAttributes struct {
	// Query ZK likeness verification proof. Required for endpoint `/v2/balances/likeness_registry`.
	Proof types.ZKProof `json:"proof"`
}
