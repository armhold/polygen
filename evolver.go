package polygen

import (
	"image"
	"log"
	"sort"
	"math/rand"
	"time"
	"image/color"
	"image/draw"
)

func Evolve(maxGen int, referenceImg image.Image, destFile string, safeImages []*SafeImage) {
	refImgRGBA := ConvertToRGBA(referenceImg)

	w := refImgRGBA.Bounds().Dx()
	h := refImgRGBA.Bounds().Dy()

	var population []*Candidate

	startTime := time.Now()

	for i := 0; i < PopulationCount; i++ {
		c := RandomCandidate(w, h)
		evaluateCandidate(c, refImgRGBA)
		population = append(population, c)
	}

	generationsSinceChange := 0
	mostFit := population[0]

	for i := 0; i < maxGen; i++ {
		shufflePopulation(population)
		parentCount := len(population)

		for j := 0; j < parentCount; j += 2 {
			m1 := population[j]
			m2 := population[j + 1]

			child := m1.Mate(m2)
			evaluateCandidate(child, refImgRGBA)
			population = append(population, child)
		}

		// after sort, the best will be at [0], worst will be at [len() - 1]
		sort.Sort(ByFitness(population))

		if i % 10 == 0 {
			printStats(population, i, generationsSinceChange, startTime)
		}

//		log.Printf("pop count: %d", len(population))
		for i := 0; i < len(safeImages); i++ {
			safeImages[i].Update(population[i].img)
		}

		// evict the least-fit individuals
		population = population[:PopulationCount]

		prev := mostFit
		mostFit = population[0]

		if mostFit.Fitness < prev.Fitness {
			generationsSinceChange = 0
		} else {
			generationsSinceChange++
		}
	}

	mostFit.DrawAndSave(destFile)
	log.Printf("after %d generations, fitness is: %d, saved to %s", maxGen, mostFit.Fitness, destFile)
}


// for comparison, create a near-perfect copy of the ref image, with only a few pixels changed
func createNearCopy(refImg image.Image) image.Image {
	result := image.NewRGBA(refImg.Bounds())
	b := result.Bounds()

	draw.Draw(result, b, refImg, b.Min, draw.Src)

	for i := 0; i < 5; i++ {
		result.Set(b.Min.X + i, b.Min.Y, color.Black)
	}

	return result
}



func evaluateCandidate(c *Candidate, referenceImg *image.RGBA) {
	// for comparison,
	//almostPerfect := image.NewRGBA(referenceImg.Bounds())
	//draw.Draw(almostPerfect, almostPerfect.Bounds(), referenceImg, almostPerfect.Bounds().Min, draw.Src)
	//almostPerfect.Set(50, 50, color.Black)
	//diff, err := Compare(referenceImg, almostPerfect)

//	diff, err := Compare(referenceImg, c.img)
	diff, err := FastCompare(referenceImg, c.img)

	if err != nil {
		log.Fatalf("error comparing images: %s", err)
	}

	c.Fitness = diff
}

func shufflePopulation(population []*Candidate) {
	for i := range population {
		j := rand.Intn(i + 1)
		population[i], population[j] = population[j], population[i]
	}
}

func printStats(sortedPop []*Candidate, generations, generationsSinceChange int, startTime time.Time) {
	dur := time.Since(startTime)
	best := sortedPop[0].Fitness
	worst := sortedPop[len(sortedPop)-1].Fitness

	log.Printf("dur: %s, gen: %d, since change: %d, best: %d, worst: %d", dur, generations, generationsSinceChange, best, worst)
}