package data

type Referral struct {
	ID         string `db:"id"`
	UserDID    string `db:"user_did"`
	IsConsumed bool   `db:"is_consumed"`
	CreatedAt  int32  `db:"created_at"`
}

type ReferralsQ interface {
	New() ReferralsQ
	Insert(...Referral) error
	Consume(ids ...string) (consumedIDs []string, err error)
	ConsumeFirst(did string, count uint64) error

	Select() ([]Referral, error)
	Get(id string) (*Referral, error)
	Count() (uint64, error)

	FilterByUserDID(string) ReferralsQ
	FilterByIsConsumed(bool) ReferralsQ
}
