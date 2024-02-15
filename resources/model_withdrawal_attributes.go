/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type WithdrawalAttributes struct {
	// Rarimo address which points were withdrawn to. Can be any valid address.
	Address string `json:"address"`
	// Amount of points withdrawn
	Amount int64 `json:"amount"`
	// Unix timestamp of withdrawal creation
	CreatedAt int32 `json:"created_at"`
}
