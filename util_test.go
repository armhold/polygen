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
