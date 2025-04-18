package data

type FaceEventBalance struct {
	Nullifier string `db:"nullifier"`
	Amount    int64  `db:"amount"`
	CreatedAt int32  `db:"created_at"`
}

type FaceEventBalanceQ interface {
	New() FaceEventBalanceQ
	Insert(FaceEventBalance) error
	Update(map[string]any) error
	Transaction(f func() error) error

	Get() (*FaceEventBalance, error)

	FilterByNullifier(nullifier ...string) FaceEventBalanceQ
}
