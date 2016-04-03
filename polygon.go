package polygen

import (
	"fmt"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"math/rand"
	"log"
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
	PolygonsPerIndividual    = 50
	MaxPolygonPoints         = 6
	MinPolygonPoints         = 3
	PointMutationMaxDistance = 5
)

var (
	Mutations = []int{MutationColor, MutationPoint, MutationZOrder, MutationAddOrDeletePoint, MutationAlpha}
)


func init() {
	log.Printf("MutationAlpha: %d", MutationAlpha)
	log.Printf("MutationColor: %d", MutationColor)
	log.Printf("MutationPoint: %d", MutationPoint)
	log.Printf("MutationZOrder: %d", MutationZOrder)
	log.Printf("MutationAddOrDeletePoint: %d", MutationAddOrDeletePoint)
}


type Candidate struct {
	w, h     int
	Polygons []*Polygon
	Fitness  int64
	img      *image.RGBA
}

type Point struct {
	X, Y int
}

type Polygon struct {
	Points []*Point
	color.Color
}

func RandomCandidate(w, h int) *Candidate {
	result := &Candidate{w: w, h: h, Polygons: make([]*Polygon, PolygonsPerIndividual)}
	for i := 0; i < len(result.Polygons); i++ {
		result.Polygons[i] = RandomPolygon(w, h)
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

func RandomPoint(maxW, maxH int) *Point {
	return &Point{rand.Intn(maxW), rand.Intn(maxH)}
}

func (m1 *Candidate) Mate(m2 *Candidate) *Candidate {
	w, h := m1.w, m1.h
	crossover := rand.Intn(len(m1.Polygons))
	polygons := make([]*Polygon, len(m1.Polygons))

	shouldShufflePolygons := false
	for i := 0; i < len(polygons); i++ {
		var p Polygon

		if i < crossover {
			p = *m1.Polygons[i] // NB copy the polygon, not the pointer
		} else {
			p = *m2.Polygons[i]
		}

		if shouldMutate() {
			switch randomMutation() {
			case MutationColor:
				p.Color = MutateColor(p.Color)

			case MutationAlpha:
				p.Color = MutateAlpha(p.Color)

			case MutationPoint:
				i := rand.Intn(len(p.Points))
				p.Points[i].MutateNearby(w, h)

			case MutationZOrder:
				shouldShufflePolygons = true

			case MutationAddOrDeletePoint:
				if len(p.Points) == MinPolygonPoints {
					// can't delete
					p.AddPoint(RandomPoint(w, h))
				} else if len(p.Points) == MaxPolygonPoints {
					// can't add
					p.DeleteRandomPoint()
				} else {
					// we can do either add or delete
					if NextBool() {
						p.AddPoint(RandomPoint(w, h))
					} else {
						p.DeleteRandomPoint()
					}
				}
			}
		}

		polygons[i] = &p
	}

	if shouldShufflePolygons {
		shufflePolygonZOrder(polygons)
	}

	result := &Candidate{w: w, h: h, Polygons: polygons}
	result.RenderImage()
	return result
}

func (c *Candidate) MutatedCopy() *Candidate {
	result := &Candidate{w: c.w, h: c.h, Polygons: make([]*Polygon, PolygonsPerIndividual)}
	polygons := result.Polygons

	for i := 0; i < len(result.Polygons); i++ {
		p := *c.Polygons[i]  // copy
		polygons[i] = &p
	}

	shouldShufflePolygons := false

	// make 3 mutations
	for i := 0; i < 3 ; i++ {
		locus := rand.Intn(len(polygons))
		p := polygons[locus]
		switch randomMutation() {
		case MutationColor:
			p.Color = MutateColor(p.Color)

		case MutationAlpha:
			p.Color = MutateAlpha(p.Color)

		case MutationPoint:
			i := rand.Intn(len(p.Points))
			p.Points[i].MutateNearby(c.w, c.h)

		case MutationZOrder:
			shouldShufflePolygons = true

		case MutationAddOrDeletePoint:
			if len(p.Points) == MinPolygonPoints {
				// can't delete
				p.AddPoint(RandomPoint(c.w, c.h))
			} else if len(p.Points) == MaxPolygonPoints {
				// can't add
				p.DeleteRandomPoint()
			} else {
				// we can do either add or delete
				if NextBool() {
					p.AddPoint(RandomPoint(c.w, c.h))
				} else {
					p.DeleteRandomPoint()
				}
			}
		}
	}

	if shouldShufflePolygons {
		shufflePolygonZOrder(polygons)
	}

	result.RenderImage()

	return result
}

func (p *Polygon) AddPoint(point *Point) {
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
	cd.img = image.NewRGBA(image.Rect(0, 0, cd.w, cd.h))
	gc := draw2dimg.NewGraphicContext(cd.img)

	// paint the whole thing black to start
	gc.SetFillColor(color.Black)
	gc.MoveTo(0, 0)
	gc.LineTo(float64(cd.w), 0)
	gc.LineTo(float64(cd.w), float64(cd.h))
	gc.LineTo(0, float64(cd.h))
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
