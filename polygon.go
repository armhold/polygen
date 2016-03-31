package main

import (
	"image/color"
	"time"
	"math/rand"
	"log"
	"image"
	"github.com/llgcode/draw2d/draw2dimg"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

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

type Individual interface {
	Fitness() int
	BreedWith(Individual) Individual
}

type PolygonSet []*Polygon

func Evolve() {
	var population []PolygonSet

	for i := 0; i < PopulationCount; i++ {
		var individual PolygonSet
		for j := 0; j < PolygonsPerIndividual; j++ {
			individual = append(individual, RandomPolygon())
		}

		population = append(population, individual)
	}

	DrawAndSave(population[0])

	//log.Printf("population: %+v", population)
}

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

	numPoints := randomInt(MinPolygonPoints, MaxPolygonPoints)

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

func DrawAndSave(s PolygonSet) {
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

	draw2dimg.SaveToPngFile("output.png", dest)
}

func randomInt(min, max int) int {
	return rand.Intn(max - min) + min
}


func main() {
	Evolve()
}
