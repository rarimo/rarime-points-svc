package data

type Referral struct {
	ID        string `db:"id"`
	Nullifier string `db:"nullifier"`
	UsageLeft int32  `db:"usage_left"`
	Status    string `db:"status"`
}

type ReferralsQ interface {
	New() ReferralsQ
	Insert(...Referral) error
	Consume(ids ...string) (consumedIDs []string, err error)
	ConsumeFirst(nullifier string, count uint64) error

	Select() ([]Referral, error)
	Get(id string) (*Referral, error)
	Count() (uint64, error)

	Update(usageLeft int) (*Referral, error)
	DeleteByID(ids ...string) error
	Transaction(f func() error) error

	WithStatus() ReferralsQ
	FilterByNullifier(string) ReferralsQ
	FilterConsumed() ReferralsQ
	FilterByID(nullifier string) ReferralsQ
}
