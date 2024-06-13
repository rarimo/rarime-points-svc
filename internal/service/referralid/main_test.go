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
			name: "Valid nullifier with index 0",
			did:  "2184ae1f990d26aa5bb84d54dc945ac3cce569cd828269802f0fa5c5c28f30a7",
			want: "6xM70VgX4eh",
		},
		{
			name:  "Valid nullifier with index 1",
			did:   "2184ae1f990d26aa5bb84d54dc945ac3cce569cd828269802f0fa5c5c28f30a7",
			index: 1,
			want:  "eLHv3hj5txB",
		},
		{
			name:  "Valid nullifier with index 258",
			did:   "2184ae1f990d26aa5bb84d54dc945ac3cce569cd828269802f0fa5c5c28f30a7",
			index: 258,
			want:  "1hhJaHQB13G",
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
		name      string
		nullifier string
		index     uint64
		count     uint64
		want      []string
	}{
		{
			name:      "Valid nullifier for basic balance creation",
			nullifier: "2184ae1f990d26aa5bb84d54dc945ac3cce569cd828269802f0fa5c5c28f30a7",
			count:     5,
			want:      []string{"6xM70VgX4eh", "eLHv3hj5txB", "8Mu12YhyDVQ", "4l3LwW9p77V", "bLnCgkUOPWT"},
		},
		{
			name:      "Valid nullifier for start from non-zero index",
			nullifier: "2184ae1f990d26aa5bb84d54dc945ac3cce569cd828269802f0fa5c5c28f30a7",
			index:     2,
			count:     3,
			want:      []string{"8Mu12YhyDVQ", "4l3LwW9p77V", "bLnCgkUOPWT"},
		},
		{
			name:      "Valid nullifier, no count",
			nullifier: "2184ae1f990d26aa5bb84d54dc945ac3cce569cd828269802f0fa5c5c28f30a7",
			index:     8,
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
			got := NewMany(tt.nullifier, tt.count, tt.index)
			if !equal(got, tt.want) {
				t.Errorf("NewMany() = %s, want %s", got, tt.want)
			}
		})
	}
}
