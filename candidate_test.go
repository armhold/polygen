package polygen

import (
	"testing"
)

func BenchmarkMutateInPlace(b *testing.B) {
	c := RandomCandidate(200, 200, 50)

	for i := 0; i < b.N; i++ {
		c.MutateInPlace()
	}
}

func BenchmarkRenderImage(b *testing.B) {
	c := RandomCandidate(200, 200, 50)

	for i := 0; i < b.N; i++ {
		c.RenderImage()
	}
}
