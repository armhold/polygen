package polygen

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"math/rand"
	"log"
	"encoding/gob"
)

const (
	MutationAlpha            = iota
	MutationColor            = iota
	MutationPoint            = iota
	MutationZOrder           = iota
	MutationAddOrDeletePoint = iota
)

const (
	MutationChance           = 0.25
	PopulationCount          = 10
	MaxPolygonPoints         = 6
	MinPolygonPoints         = 3
	PointMutationMaxDistance = 100
)

var (
	Mutations = []int{MutationColor, MutationPoint, MutationAlpha, MutationZOrder, MutationAddOrDeletePoint}
)


func init() {
	// need to give an example of a concrete type for the color.Color interface
	gob.Register(RandomColor())
}


type Candidate struct {
	W, H     int
	Polygons []*Polygon
	Fitness  int64
	img      *image.RGBA
}

type Point struct {
	X, Y int
}

type Polygon struct {
	Points []Point
	color.Color
}

func (p* Polygon) Copy() *Polygon {
	result := &Polygon{Color: p.Color}
	for i := 0; i < len(p.Points); i++ {
		result.Points = append(result.Points, p.Points[i])
	}

	return result
}

func RandomCandidate(w, h, polyCount int) *Candidate {
	result := &Candidate{W: w, H: h}
	for i := 0; i < polyCount; i++ {
		result.Polygons = append(result.Polygons, RandomPolygon(w, h))
	}

	result.RenderImage()

	return result
}

func RandomPolygon(maxW, maxH int) *Polygon {
	result := &Polygon{}
	result.Color = color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256))}

	numPoints := RandomInt(MinPolygonPoints, MaxPolygonPoints+1)

	for i := 0; i < numPoints; i++ {
		result.AddPoint(RandomPoint(maxW, maxH))
	}

	return result
}

func RandomPoint(maxW, maxH int) Point {
	return Point{rand.Intn(maxW), rand.Intn(maxH)}
}

// does not copy image- we assume the copy will be mutated after
func (c *Candidate) CopyOf() *Candidate {
	result := &Candidate{W: c.W, H: c.H}
	for i := 0; i < len(c.Polygons); i++ {
		result.Polygons = append(result.Polygons, c.Polygons[i].Copy())
  	}

	return result
}

func (c *Candidate) MutateInPlace() {
	shouldShufflePolygons := false

	// make 3 mutations
	for i := 0; i < 3 ; i++ {
		locus := rand.Intn(len(c.Polygons))
		pgon := c.Polygons[locus]
		switch randomMutation() {
		case MutationColor:
			pgon.Color = MutateColor(pgon.Color)

		case MutationAlpha:
			pgon.Color = MutateAlpha(pgon.Color)

		case MutationPoint:
			pi := rand.Intn(len(pgon.Points))
			pgon.Points[pi].MutateNearby(c.W, c.H)

		case MutationZOrder:
			shouldShufflePolygons = true

		case MutationAddOrDeletePoint:
			if len(pgon.Points) == MinPolygonPoints {
				// can't delete
				pgon.AddPoint(RandomPoint(c.W, c.H))
			} else if len(pgon.Points) == MaxPolygonPoints {
				// can't add
				pgon.DeleteRandomPoint()
			} else {
				// we can do either add or delete
				if NextBool() {
					pgon.AddPoint(RandomPoint(c.W, c.H))
				} else {
					pgon.DeleteRandomPoint()
				}
			}


		default:
			log.Fatal("fell through")
		}
	}

	if shouldShufflePolygons {
		shufflePolygonZOrder(c.Polygons)
	}

	c.RenderImage()
}

func (p *Polygon) AddPoint(point Point) {
	p.Points = append(p.Points, point)
}


func (p *Polygon) DeleteRandomPoint() {
	i := rand.Intn(len(p.Points))
	p.Points = append(p.Points[:i], p.Points[i+1:]...)
}

func (p *Point) MutateNearby(maxW, maxH int) {
	xDelta := rand.Intn(PointMutationMaxDistance + 1)
	if NextBool() {
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
	if NextBool() {
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

func RandomColor() color.Color {
	c := color.NRGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: uint8(rand.Intn(256)) }
	return color.RGBAModel.Convert(c)
}

func MutateColor(c color.Color) color.Color {
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

func MutateAlpha(c color.Color) color.Color {
	// get the non-premultiplied rgba values
	nrgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	nrgba.A = uint8(rand.Intn(256))

	return color.RGBAModel.Convert(nrgba)
}

func randomMutation() int {
	return Mutations[rand.Intn(len(Mutations))]
}

func (cd *Candidate) RenderImage() {
	cd.img = image.NewRGBA(image.Rect(0, 0, cd.W, cd.H))
	gc := draw2dimg.NewGraphicContext(cd.img)

	// paint the whole thing black to start
	gc.SetFillColor(color.Black)
	gc.MoveTo(0, 0)
	gc.LineTo(float64(cd.W), 0)
	gc.LineTo(float64(cd.W), float64(cd.H))
	gc.LineTo(0, float64(cd.H))
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

func (cd *Candidate) DrawAndSave(destFile string) {
	draw2dimg.SaveToPngFile(destFile, cd.img)
}

func (cd *Candidate) String() string {
	return fmt.Sprintf("fitness: %d", cd.Fitness)
}

func shouldMutate() bool {
	return rand.Float32() < MutationChance
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
