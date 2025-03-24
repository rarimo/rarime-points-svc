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
	// WithoutExpiredStatus filters out referral codes that have an “expired” status.
	// The status “expired” is assigned to codes that have been used, but the party that used them did not complete the passport scanning procedure by the set time.
	// In this case, usage_left is set to -1, and new codes are generated to replace the expired ones.
	// This allows you to keep a history of all codes used by the user.
	// It can be used only after applying the WithStatus filter, since the statuses are defined in it.
	WithoutExpiredStatus() ReferralsQ

	FilterByNullifier(string) ReferralsQ
	FilterConsumed() ReferralsQ
	FilterByID(id string) ReferralsQ
}
