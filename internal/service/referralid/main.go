package referralid

import (
	"crypto/sha256"
	"math/big"
)

const bytesCount = 6

// New returns a base62 of the first 6 bytes of SHA-256 of the input
func New(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	first := hash.Sum(nil)[:bytesCount]
	return base62Encode(first)
}

// using alphabet: [0-9a-zA-Z]
func base62Encode(b []byte) string {
	i := new(big.Int)
	i.SetBytes(b)
	return i.Text(62)
}
