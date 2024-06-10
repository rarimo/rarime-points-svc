/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "encoding/json"

type CountriesConfig struct {
	Key
	Attributes CountriesConfigAttributes `json:"attributes"`
}
type CountriesConfigResponse struct {
	Data     CountriesConfig `json:"data"`
	Included Included        `json:"included"`
}

type CountriesConfigListResponse struct {
	Data     []CountriesConfig `json:"data"`
	Included Included          `json:"included"`
	Links    *Links            `json:"links"`
	Meta     json.RawMessage   `json:"meta,omitempty"`
}

func (r *CountriesConfigListResponse) PutMeta(v interface{}) (err error) {
	r.Meta, err = json.Marshal(v)
	return err
}

func (r *CountriesConfigListResponse) GetMeta(out interface{}) error {
	return json.Unmarshal(r.Meta, out)
}

// MustCountriesConfig - returns CountriesConfig from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustCountriesConfig(key Key) *CountriesConfig {
	var countriesConfig CountriesConfig
	if c.tryFindEntry(key, &countriesConfig) {
		return &countriesConfig
	}
	return nil
}
