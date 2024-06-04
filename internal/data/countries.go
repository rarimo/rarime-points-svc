package data

type CountriesQ interface {
	New() CountriesQ
	Insert(countries ...Country) error
	Update(limit, addReserved, addWithdrawn *int64, isDisabled *bool) error
	Select() ([]Country, error)
	Get() (*Country, error)
	FilterByCodes(codes ...string) CountriesQ
	FilterDisabled(bool) CountriesQ
}

type Country struct {
	Code         string `db:"code"`
	ReserveLimit int64  `db:"reserve_limit"`
	Reserved     int64  `db:"reserved"`
	Withdrawn    int64  `db:"withdrawn"`
	IsDisabled   bool   `db:"is_disabled"`
}
