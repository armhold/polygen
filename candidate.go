package polygen

import (
	"encoding/gob"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	MutationAlpha            = iota
	MutationColor            = iota
	MutationPoint            = iota
	MutationZOrder           = iota
	MutationAddOrDeletePoint = iota
)

const (
	PopulationCount          = 10
	MaxPolygonPoints         = 6
	MinPolygonPoints         = 3
	PointMutationMaxDistance = 5
	MutationsPerIteration    = 1 // originally had 3, but 1 seems to work best here
)

var (
	Mutations = []int{MutationColor, MutationPoint, MutationAlpha, MutationZOrder, MutationAddOrDeletePoint}
)

func init() {
	// need to give an example of a concrete type for the color.Color interface
	gob.Register(randomColor())
}

// Candidate is a potential solution (set of polygons) to the problem of how to best represent the reference image.
type Candidate struct {
	W, H     int
	Polygons []*Polygon
	img      *image.RGBA
	Fitness  uint64
}

// Polygon is a set of points with a given fill color.
type Polygon struct {
	Points []Point
	color.Color
}

// Point defines a vertex in a Polygon.
type Point struct {
	X, Y int
}

func (p *Polygon) copyOf() *Polygon {
	result := &Polygon{Color: p.Color}
	for i := 0; i < len(p.Points); i++ {
		result.Points = append(result.Points, p.Points[i])
	}

	return result
}

func randomCandidate(w, h, polyCount int) *Candidate {
	result := &Candidate{W: w, H: h}
	for i := 0; i < polyCount; i++ {
		result.Polygons = append(result.Polygons, randomPolygon(w, h))
	}

	result.renderImage()

	return result
}

func randomPolygon(maxW, maxH int) *Polygon {
	result := &Polygon{}
	result.Color = randomColor()

	numPoints := RandomInt(MinPolygonPoints, MaxPolygonPoints+1)

	for i := 0; i < numPoints; i++ {
		result.addPoint(randomPoint(maxW, maxH))
	}

	return result
}

func randomPoint(maxW, maxH int) Point {
	return Point{rand.Intn(maxW), rand.Intn(maxH)}
}

// Copies the Candidate, minus the img (we assume the copy will be mutated/rendered after).
func (c *Candidate) copyOf() *Candidate {
	result := &Candidate{W: c.W, H: c.H}
	for i := 0; i < len(c.Polygons); i++ {
		result.Polygons = append(result.Polygons, c.Polygons[i].copyOf())
	}

	return result
}

// mutateInPlace chooses a random polygon from the candidate and makes a random mutation to it.
func (c *Candidate) mutateInPlace() {
	locus := rand.Intn(len(c.Polygons))
	poly := c.Polygons[locus]
	switch randomMutation() {
	case MutationColor:
		poly.Color = mutateColor(poly.Color)

	case MutationAlpha:
		poly.Color = mutateAlpha(poly.Color)

	case MutationPoint:
		pointIndex := rand.Intn(len(poly.Points))
		poly.Points[pointIndex].mutateNearby(c.W, c.H)

	case MutationZOrder:
		shufflePolygonZOrder(c.Polygons)

	case MutationAddOrDeletePoint:
		if len(poly.Points) == MinPolygonPoints {
			// can't delete
			poly.addPoint(randomPoint(c.W, c.H))
		} else if len(poly.Points) == MaxPolygonPoints {
			// can't add
			poly.deleteRandomPoint()
		} else {
			// we can do either add or delete
			if RandomBool() {
				poly.addPoint(randomPoint(c.W, c.H))
			} else {
				poly.deleteRandomPoint()
			}
		}

	default:
		log.Fatal("fell through")
	}
}

func (p *Polygon) addPoint(point Point) {
	p.Points = append(p.Points, point)
}

func (p *Polygon) deleteRandomPoint() {
	i := rand.Intn(len(p.Points))
	p.Points = append(p.Points[:i], p.Points[i+1:]...)
}

// mutateNearby alters the point by moving it a few pixels.
func (p *Point) mutateNearby(maxW, maxH int) {
	xDelta := rand.Intn(PointMutationMaxDistance + 1)
	if RandomBool() {
		xDelta = -xDelta
	}

	x := p.X + xDelta
	if x < 0 {
		x = 0
	}

	if x >= maxW {
		x = maxW - 1
	}
	p.X = x

	yDelta := rand.Intn(PointMutationMaxDistance + 1)
	if RandomBool() {
		yDelta = -yDelta
	}

	y := p.Y + yDelta
	if y < 0 {
		y = 0
	}

	if y >= maxH {
		y = maxH - 1
	}
	p.Y = y
}

// randomColor returns a color with completely random values for RGBA.
func randomColor() color.Color {
	// start with non-premultiplied RGBA
	c := color.NRGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: uint8(rand.Intn(256))}
	return color.RGBAModel.Convert(c)
}

// mutateColor returns a new color with a single random mutation to one of the RGBA values.
func mutateColor(c color.Color) color.Color {
	// get the non-premultiplied rgba values
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)

	// randomly select one of the r/g/b/a values to mutate
	i := rand.Intn(4)
	val := uint8(rand.Intn(256))

	switch i {
	case 0:
		nrgba.R = val
	case 1:
		nrgba.G = val
	case 2:
		nrgba.B = val
	case 3:
		nrgba.A = val
	}

	return color.RGBAModel.Convert(nrgba)
}

// mutateAlpha a new color whose alpha level has been randomly modified.
func mutateAlpha(c color.Color) color.Color {
	// get the non-premultiplied rgba values
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	nrgba.A = uint8(rand.Intn(256))

	return color.RGBAModel.Convert(nrgba)
}

func randomMutation() int {
	return Mutations[rand.Intn(len(Mutations))]
}

func (cd *Candidate) renderImage() {
	cd.img = image.NewRGBA(image.Rect(0, 0, cd.W, cd.H))
	gc := draw2dimg.NewGraphicContext(cd.img)

	// paint the whole thing black to start
	gc.SetFillColor(color.Black)
	gc.MoveTo(0, 0)
	gc.LineTo(float64(cd.W-1), 0)
	gc.LineTo(float64(cd.W-1), float64(cd.H-1))
	gc.LineTo(0, float64(cd.H-1))
	gc.Close()
	gc.Fill()

	gc.SetLineWidth(1)

	for _, polygon := range cd.Polygons {
		gc.SetStrokeColor(polygon.Color)
		gc.SetFillColor(polygon.Color)

		firstPoint := polygon.Points[0]
		gc.MoveTo(float64(firstPoint.X), float64(firstPoint.Y))

		for _, point := range polygon.Points[1:] {
			gc.LineTo(float64(point.X), float64(point.Y))
		}

		gc.Close()
		//gc.FillStroke()
		gc.Fill()
	}
}

func (cd *Candidate) drawAndSave(destFile string) error {
	log.Printf("saving output image to: %s", destFile)
	return draw2dimg.SaveToPngFile(destFile, cd.img)
}

func shufflePolygonZOrder(polygons []*Polygon) {
	for i := range polygons {
		j := rand.Intn(i + 1)
		polygons[i], polygons[j] = polygons[j], polygons[i]
	}
}

type ByFitness []*Candidate

func (cds ByFitness) Len() int           { return len(cds) }
func (cds ByFitness) Swap(i, j int)      { cds[i], cds[j] = cds[j], cds[i] }
func (cds ByFitness) Less(i, j int) bool { return cds[i].Fitness < cds[j].Fitness }
