package referralid

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		did   string
		index uint64
		want  string
	}{
		{
			name: "Valid DID with index 0",
			did:  "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			want: "bDSCcQB8Hhk",
		},
		{
			name:  "Valid DID with index 1",
			did:   "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			index: 1,
			want:  "9csIL7dW65m",
		},
		{
			name:  "Valid DID with index 258",
			did:   "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			index: 258,
			want:  "73k3bdYaFWM",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.did, tt.index); got != tt.want {
				t.Errorf("New() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestNewMany(t *testing.T) {
	tests := []struct {
		name  string
		did   string
		index uint64
		count uint64
		want  []string
	}{
		{
			name:  "Valid DID for basic balance creation",
			did:   "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			count: 5,
			want:  []string{"bDSCcQB8Hhk", "9csIL7dW65m", "lcRN9LZliVw", "kjmJHA8IdRA", "lcz9ZTJtWgA"},
		},
		{
			name:  "Valid DID for start from non-zero index",
			did:   "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			index: 2,
			count: 3,
			want:  []string{"lcRN9LZliVw", "kjmJHA8IdRA", "lcz9ZTJtWgA"},
		},
		{
			name:  "Valid DID, no count",
			did:   "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			index: 8,
		},
	}

	equal := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMany(tt.did, tt.count, tt.index)
			if !equal(got, tt.want) {
				t.Errorf("NewMany() = %s, want %s", got, tt.want)
			}
		})
	}
}
