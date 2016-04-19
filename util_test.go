package polygen

import (
	"testing"
)

func TestRandomInt(t *testing.T) {
	for i := 0; i < 1000; i++ {
		d := RandomInt(1, 4)

		if d < 1 || d > 3 {
			t.Fatalf("expected d to be in [1..3], was: %d", d)
		}
	}
}

func TestDeriveCheckpointFile(t *testing.T) {
	var examples = []struct {
		file  string
		cpArg string
		poly  int
		out   string
	}{
		{"/foo/bar/baz/foo.jpg", "", 50, "foo-50-checkpoint.tmp"},
		{"images/foo.jpg", "", 50, "foo-50-checkpoint.tmp"},
		{"foo.jpg", "", 50, "foo-50-checkpoint.tmp"},
		{"images/foo.jpg", "", 100, "foo-100-checkpoint.tmp"},
		{"images/foo-123.jpg", "", 50, "foo-123-50-checkpoint.tmp"},
		{"foo.jpg", "custom-arg-given.tmp", 50, "custom-arg-given.tmp"},
		{"foo.jpg", "/some/path/custom-arg-given.tmp", 50, "/some/path/custom-arg-given.tmp"},
	}

	for _, tt := range examples {
		expected := tt.out
		actual := DeriveCheckpointFile(tt.file, tt.cpArg, tt.poly)

		if actual != expected {
			t.Errorf("wanted %s, got: %s", expected, actual)
		}
	}
}
