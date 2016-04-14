package polygen

import (
	"testing"
	"reflect"
)

func BenchmarkMutateInPlace(b *testing.B) {
	c := randomCandidate(200, 200, 50)

	for i := 0; i < b.N; i++ {
		c.mutateInPlace()
	}
}

func BenchmarkRenderImage(b *testing.B) {
	c := randomCandidate(200, 200, 50)

	for i := 0; i < b.N; i++ {
		c.renderImage()
	}
}

func TestCopy(t *testing.T) {
	c1 := randomCandidate(100, 100, 10)

	// don't care about these two fields
	c1.img = nil
	c1.Fitness = 0

	c2 := c1.copyOf()

	if ! reflect.DeepEqual(c1, c2) {
		t.Fatalf("c1 != c2: %+v, %+v", c1, c2)
	}
}
