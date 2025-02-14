/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type Maintenance struct {
	Key
	Attributes MaintenanceAttributes `json:"attributes"`
}
type MaintenanceResponse struct {
	Data     Maintenance `json:"data"`
	Included Included    `json:"included"`
}

type MaintenanceListResponse struct {
	Data     []Maintenance   `json:"data"`
	Included Included        `json:"included"`
	Links    *Links          `json:"links"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

func (r *MaintenanceListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *MaintenanceListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustMaintenance - returns Maintenance from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustMaintenance(key Key) *Maintenance {
	var maintenance Maintenance
	if c.tryFindEntry(key, &maintenance) {
		return &maintenance
	}
	return nil
}
