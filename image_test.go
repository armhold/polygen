package polygen

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"
)

func TestCompareBounds(t *testing.T) {
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	img2 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	img3 := image.NewRGBA(image.Rect(10, 10, 100, 100))

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

func TestConvert(t *testing.T) {
	rect := image.Rect(0, 0, 100, 100)

	rgbImg := image.NewRGBA(rect)
	result := ConvertToRGBA(rgbImg)

	if reflect.ValueOf(result).Pointer() != reflect.ValueOf(rgbImg).Pointer() {
		t.Fatalf("expected to get the same pointer back for RGBA image")
	}

	cmykImg := image.NewCMYK(rect)
	result = ConvertToRGBA(cmykImg)
	if reflect.ValueOf(result).Pointer() == reflect.ValueOf(cmykImg).Pointer() {
		t.Fatalf("expected to get different pointer back for non-RGBA image")
	}
}
