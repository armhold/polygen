package polygen

import (
	"image"
	"log"
	"sort"
	"math/rand"
	"time"
)

func Evolve(maxGen int, referenceImg image.Image, destFile string, safeImage *SafeImage) {
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

	for i := 0; i < maxGen; i++ {
		//log.Printf("generation %d", i)

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
		//for _, candidate := range population {
		//	log.Print(candidate)
		//}

		if i % 10 == 0 {
			printStats(population, i, startTime)
		}

		//bestChild := population[parentCount]


		// evict the least-fit
		population = population[:PopulationCount]

		mostFit := population[0]
		//mostFit.DrawAndSave(destFile)
		//safeImage.Update(mostFit.img)
		safeImage.Update(mostFit.img)

	}

	mostFit := population[0]
	mostFit.DrawAndSave(destFile)
	log.Printf("after %d generations, fitness is: %d, saved to %s", maxGen, mostFit.Fitness, destFile)
}

func evaluateCandidate(c *Candidate, referenceImg *image.RGBA) {
	//diff, err := Compare(referenceImg, c.img)
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

func printStats(sortedPop []*Candidate, generations int, startTime time.Time) {
	dur := time.Since(startTime)
	best := sortedPop[0].Fitness
	worst := sortedPop[len(sortedPop)-1].Fitness

	log.Printf("dur: %s, generations: %d, best: %d, worst: %d", dur, generations, best, worst)
}