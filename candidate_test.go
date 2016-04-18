package polygen

import (
	"reflect"
	"testing"
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

func TestCandidateCopyOf(t *testing.T) {
	c1 := randomCandidate(100, 100, 10)

	// don't care about these two fields
	c1.img = nil
	c1.Fitness = 0

	c2 := c1.copyOf()

	if !reflect.DeepEqual(c1, c2) {
		t.Fatalf("c1 != c2: %+v, %+v", c1, c2)
	}
}

// check that copyOf() actually copies the polygon's points (vs just copying their pointers). Had a mutability bug here.
func TestPolygonCopyOf(t *testing.T) {
	p1 := randomPolygon(100, 100)
	p2 := p1.copyOf()

	// initially, they should be equal
	if !reflect.DeepEqual(p1, p2) {
		t.Fatalf("p1 != p2: %+v, %+v", p1, p2)
	}

	// but changing a point in p1 should not affect p2
	p1.Points[0].mutateNearby(100, 100)
	if reflect.DeepEqual(p1, p2) {
		t.Fatalf("p1 should have diverged from p2: %+v, %+v", p1, p2)
	}
}
