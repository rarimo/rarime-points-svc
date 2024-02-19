package referralid

import "testing"

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "Valid DID",
			s:    "did:iden3:readonly:tUDjWxnVJNi7t3FudukqrUcNwF5KVGoWgim5pp2jV",
			want: "zgsScguZ",
		},
		{
			name: "Arbitrary string",
			s:    "any string $@",
			want: "Etv79RQ0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.s); got != tt.want {
				t.Errorf("New() = %s, want %s", got, tt.want)
			}
		})
	}
}
