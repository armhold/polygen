package polygen

import (
	"image/color"
	"math/rand"
	"log"
	"image"
	"github.com/llgcode/draw2d/draw2dimg"
	"fmt"
)

const (
	MutationColor = iota
	MutationPoint = iota
	MutationAddPoint = iota
	MutationDeletePoint = iota
)

const (
	MutationChance = 0.05
	PopulationCount = 10
	PolygonsPerIndividual = 100
	MaxPolygonPoints = 6
	MinPolygonPoints = 3
)


var (
	Mutations = []int{MutationColor, MutationPoint, MutationAddPoint, MutationDeletePoint}
)

type PolygonSet struct {
	Polygons []*Polygon
	Fitness int64
}

type Point struct {
	X, Y int
}

type Polygon struct {
	Points []*Point
	color.Color
}

func RandomPolygon(maxW, maxH int) *Polygon {
	result := &Polygon{}
	result.Color = color.RGBA{uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff)), uint8(rand.Intn(0xff))}

	numPoints := RandomInt(MinPolygonPoints, MaxPolygonPoints)

	for i := 0; i < numPoints; i++ {
		result.AddPoint(&Point{rand.Intn(maxW), rand.Intn(maxH)})
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

func (s *PolygonSet) RenderImage(w, h int) image.Image {
	dest := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetLineWidth(1)

	for _, polygon := range s.Polygons {
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

	return dest
}


func (s *PolygonSet) DrawAndSave(w, h int, destFile string) {
	img := s.RenderImage(w, h)
	draw2dimg.SaveToPngFile(destFile, img)
}

func (s *PolygonSet) String() string {
	return fmt.Sprintf("fitness: %d", s.Fitness)
}


type ByFitness []*PolygonSet
func (s ByFitness) Len() int           { return len(s) }
func (s ByFitness) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFitness) Less(i, j int) bool { return s[i].Fitness < s[j].Fitness }
