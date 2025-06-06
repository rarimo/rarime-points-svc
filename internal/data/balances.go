package data

import (
	"database/sql"

	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	ColAmount      = "amount"
	ColLevel       = "level"
	ColCountry     = "country"
	ColIsPassport  = "is_passport_proven"
	ColAnonymousID = "anonymous_id"
)

type Balance struct {
	Nullifier        string         `db:"nullifier"`
	Amount           int64          `db:"amount"`
	CreatedAt        int32          `db:"created_at"`
	UpdatedAt        int32          `db:"updated_at"`
	ReferredBy       sql.NullString `db:"referred_by"`
	Rank             *int           `db:"rank"`
	Level            int            `db:"level"`
	Country          *string        `db:"country"`
	IsPassportProven bool           `db:"is_passport_proven"`
	AnonymousID      *string        `db:"anonymous_id"`
}

type BalancesQ interface {
	New() BalancesQ
	Insert(Balance) error
	Update(map[string]any) error
	Transaction(f func() error) error

	Page(*pgdb.OffsetPageParams) BalancesQ
	Select() ([]Balance, error)
	Get() (*Balance, error)
	// GetWithRank returns balance with rank, filtered by nullifier. No other filters can be applied.
	GetWithRank(nullifier string) (*Balance, error)
	SelectWithRank() ([]Balance, error)

	Count() (int64, error)

	// WithoutPassportEvent returns balances which already
	// have scanned passport, but there no claimed events
	// for this. Filters are not applied.
	WithoutPassportEvent() ([]WithoutPassportEventBalance, error)
	WithoutReferralEvent() ([]ReferredReferrer, error)

	FilterByCreatedAtBefore(timestamp int) BalancesQ
	FilterByIsPassportProven(isProven bool) BalancesQ
	FilterByNullifier(...string) BalancesQ
	FilterDisabled() BalancesQ
	FilterByAnonymousID(id string) BalancesQ
	FilterUnverified() BalancesQ
}

type WithoutPassportEventBalance struct {
	Balance
	EventID     string      `db:"event_id"`
	EventStatus EventStatus `db:"event_status"`
}

type ReferredReferrer struct {
	Referred string `db:"referred"`
	Referrer string `db:"referrer"`
}
