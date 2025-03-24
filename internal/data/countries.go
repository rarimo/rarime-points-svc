package data

// DefaultCountryCode is the special code, where the default settings for
// countries are stored. When a user's country is not found, it must be added to
// DB with its own code and default settings.
const DefaultCountryCode = "default"

const (
	ColReserved  = "reserved"
	ColWithdrawn = "withdrawn"
)

const (
	StatusActive   = "active"
	StatusBanned   = "banned"
	StatusLimited  = "limited"
	StatusAwaiting = "awaiting"
	StatusRewarded = "rewarded"
	StatusConsumed = "consumed"
	// The status “expired” is assigned to codes that have been used, but the party that used them did not complete the passport scanning procedure by the set time.
	// In this case, usage_left is set to -1, and new codes are generated to replace the expired ones.
	// This allows you to keep a history of all codes used by the user.
	StatusExpired = "expired"
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
