package referralid

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

const bytesCount = 8

// New generates a new deterministic referral ID from nullifier and index.
func New(nullifier string, index uint64) string {
	s := fmt.Sprintf("%s_%d", nullifier, index)
	hash := sha256.New()
	hash.Write([]byte(s))
	first := hash.Sum(nil)[:bytesCount]
	return base62Encode(first)
}

// NewMany generates a bunch of referral IDs for a single nullifier with incrementing
// index. Specify non-zero index argument to start from a specific index, this is
// useful when you have stored referral IDs for this nullifier previously.
func NewMany(nullifier string, count, index uint64) []string {
	ids := make([]string, 0, count)
	for i := index; i < count+index; i++ {
		ids = append(ids, New(nullifier, i))
	}
	return ids
}

// using alphabet: [0-9a-zA-Z]
func base62Encode(b []byte) string {
	i := new(big.Int)
	i.SetBytes(b)
	return i.Text(62)
}
