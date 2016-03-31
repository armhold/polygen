package polygen

import (
	"image"
	"fmt"
	"math"
)

func Compare(img1, img2 image.Image) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for x := img1.Bounds().Min.X; x < img1.Bounds().Max.X; x++ {
		for y := img1.Bounds().Min.Y; y < img1.Bounds().Max.Y; y++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)

			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()

			// TODO: consider ignoring the Alpha, since the colors are pre-multiplied
			sum := sqDiff(r1, r2) + sqDiff(g1, g2) + sqDiff(b1, b2) + sqDiff(a1, a2)
			accumError += int64(sum)
		}
	}

	return int64(math.Sqrt(float64(accumError))), nil
}


// taken directly from image/color/color.go:
//
// sqDiff returns the squared-difference of x and y, shifted by 2 so that
// adding four of those won't overflow a uint32.
//
// x and y are both assumed to be in the range [0, 0xffff].
func sqDiff(x, y uint32) uint32 {
	var d uint32
	if x > y {
		d = x - y
	} else {
		d = y - x
	}
	return (d * d) >> 2
}
