package data

// DefaultCountryCode is the special code, where the default settings for
// countries are stored. When a user's country is not found, it must be added to
// DB with its own code and default settings.
const DefaultCountryCode = "default"

const (
	ColReserved  = "reserved"
	ColWithdrawn = "withdrawn"
)

type CountriesQ interface {
	New() CountriesQ
	Insert(countries ...Country) error
	Update(map[string]any) error
	// UpdateMany updates only reserve_limit, reserve_allowed and withdrawal_allowed
	UpdateMany([]Country) error
	Select() ([]Country, error)
	Get() (*Country, error)
	FilterByCodes(codes ...string) CountriesQ
}

type Country struct {
	Code              string `db:"code"`
	ReserveLimit      int64  `db:"reserve_limit"`
	Reserved          int64  `db:"reserved"`
	Withdrawn         int64  `db:"withdrawn"`
	ReserveAllowed    bool   `db:"reserve_allowed"`
	WithdrawalAllowed bool   `db:"withdrawal_allowed"`
}
