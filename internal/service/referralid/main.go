package referralid

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

const bytesCount = 8

// New generates a new deterministic referral ID from DID and index.
func New(did string, index uint64) string {
	s := fmt.Sprintf("%s_%d", did, index)
	hash := sha256.New()
	hash.Write([]byte(s))
	first := hash.Sum(nil)[:bytesCount]
	return base62Encode(first)
}

// NewMany generates a bunch of referral IDs for a single DID with incrementing
// index. Specify non-zero index argument to start from a specific index, this is
// useful when you have stored referral IDs for this DID previously.
func NewMany(did string, count, index uint64) []string {
	ids := make([]string, count)
	for i := index; index < count; index++ {
		ids[i] = New(did, i)
	}
	return ids
}

// using alphabet: [0-9a-zA-Z]
func base62Encode(b []byte) string {
	i := new(big.Int)
	i.SetBytes(b)
	return i.Text(62)
}
