/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type WithdrawAttributes struct {
	// Rarimo address to withdraw to. Can be any valid address.
	Address string `json:"address"`
	// Amount of points to withdraw
	Amount int64 `json:"amount"`
	// JSON encoded ZK passport verification proof.
	Proof json.RawMessage `json:"proof"`
}
