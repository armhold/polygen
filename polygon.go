package main

import (
    "image/color"
	"time"
	"math/rand"
	"log"
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

var (
	Mutations = []int{MutationColor, MutationPoint, MutationAddPoint, MutationDeletePoint}
)

type Point struct {
	X, Y int
}

type Polygon struct {
	Points []*Point
	color.Color
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

func main() {
	p := &Polygon{}
	p.AddPoint(&Point{10, 20})
	p.AddPoint(&Point{30, 100})

	p.Mutate()
}
