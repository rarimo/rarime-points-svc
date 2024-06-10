/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "github.com/iden3/go-rapidsnark/types"

type WithdrawAttributes struct {
	// Rarimo address to withdraw to. Can be any valid address.
	Address string `json:"address"`
	// Amount of points to withdraw
	Amount int64 `json:"amount"`
	// Iden3 ZK passport verification proof.
	Proof types.ZKProof `json:"proof"`
}
