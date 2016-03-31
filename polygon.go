package polygen

import (
	"image/color"
	"math/rand"
	"log"
	"image"
	"github.com/llgcode/draw2d/draw2dimg"
)

const (
	MutationColor = iota
	MutationPoint = iota
	MutationAddPoint = iota
	MutationDeletePoint = iota
)

const (
	MutationChance = 0.05
	NumGenerations = 10000
	PopulationCount = 10
	PolygonsPerIndividual = 100
	MaxPolygonPoints = 6
	MinPolygonPoints = 3
	ImageWidth = 1000
	ImageHeight = 1000
)


var (
	Mutations = []int{MutationColor, MutationPoint, MutationAddPoint, MutationDeletePoint}
)

type PolygonSet []*Polygon


type Point struct {
	X, Y int
}

type Polygon struct {
	Points []*Point
	color.Color
}

func RandomPolygon() *Polygon {
	result := &Polygon{}
	result.Color = color.RGBA{uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff))}

	numPoints := RandomInt(MinPolygonPoints, MaxPolygonPoints)

	for i := 0; i < numPoints; i++ {
		result.AddPoint(&Point{rand.Intn(ImageWidth), rand.Intn(ImageHeight)})
	}

	return result
}

func (p *Polygon) AddPoint(point *Point) {
	p.Points = append(p.Points, point)
}

func (p *Polygon) Mutate() {
	switch randomMutation() {
	case MutationColor:
		log.Printf("MutationColor")

	case MutationPoint:
		log.Printf("MutationPoint")

	case MutationAddPoint:
		log.Printf("MutationAddPoint")

	case MutationDeletePoint:
		log.Printf("MutationDeletePoint")
	}

}

func randomMutation() int {
	return Mutations[rand.Int() % len(Mutations)]
}

func (s PolygonSet) DrawAndSave(destFile string) {
	dest := image.NewRGBA(image.Rect(0, 0, ImageWidth, ImageHeight))
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetLineWidth(1)

	for _, polygon := range s {
		gc.SetStrokeColor(polygon.Color)
		gc.SetFillColor(polygon.Color)

		firstPoint := polygon.Points[0]
		gc.MoveTo(float64(firstPoint.X), float64(firstPoint.Y))

		for _, point := range polygon.Points[1:] {
			gc.LineTo(float64(point.X), float64(point.Y))
		}

		gc.Close()
		gc.FillStroke()
	}

	draw2dimg.SaveToPngFile(destFile, dest)
}

