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
	Deactivate(id string) error

	Select() ([]Referral, error)
	Get(id string) (*Referral, error)
	Count() (uint, error)

	FilterByUserDID(string) ReferralsQ
	FilterByIsConsumed(bool) ReferralsQ
}
