package referralid

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		did   string
		index uint
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
