package polygen

import (
	"testing"
	"image"
	"image/color"
	"image/draw"
)


func TestCompareBounds(t *testing.T) {
	img1 := image.Rect(0, 0, 100, 100)
	img2 := image.Rect(0, 0, 100, 100)
	img3 := image.Rect(10, 10, 100, 100)

	_, err := Compare(img1, img2)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}

	_, err = Compare(img1, img3)
	if err == nil {
		t.Fatalf("expected bounds to not be equal")
	}
}

func TestCompareDiff(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)
	img1 := image.NewRGBA(rect)
	img2 := image.NewRGBA(rect)

	blue1 := color.RGBA{0, 0, 255, 255}
	blue2 := color.RGBA{0, 0, 250, 255}

	draw.Draw(img1, img1.Bounds(), &image.Uniform{blue1}, image.ZP, draw.Src)
	draw.Draw(img2, img2.Bounds(), &image.Uniform{blue2}, image.ZP, draw.Src)

	// same img
	diff, _ := Compare(img1, img1)
	expected := int64(0)
	if diff != expected {
		t.Fatalf("expected diff to be %d, got: %d", expected, diff)
	}

	diff, _ = Compare(img1, img2)
	// arbitrary value that we came to by testing
	expected = 64249
	if diff != expected {
		t.Fatalf("expected diff to be %d, got: %d", expected, diff)
	}
}

func BenchmarkCompareDiff(b *testing.B) {
	rect := image.Rect(0, 0, 1000, 1000)
	img1 := image.NewRGBA(rect)
	img2 := image.NewRGBA(rect)

	blue1 := color.RGBA{0, 0, 255, 255}
	blue2 := color.RGBA{0, 0, 250, 255}

	draw.Draw(img1, img1.Bounds(), &image.Uniform{blue1}, image.ZP, draw.Src)
	draw.Draw(img2, img2.Bounds(), &image.Uniform{blue2}, image.ZP, draw.Src)

	for i := 0; i < b.N; i++ {
		Compare(img1, img1)
	}
}
