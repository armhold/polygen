package polygen

import (
	"image"
	"log"
	"sort"
	"math/rand"
)

type Individual interface {
	Fitness() int
	BreedWith(Individual) Individual
}

func Evolve(maxGen int, referenceImg image.Image, destFile string, safeImage *SafeImage) {
	refImgRGBA := ConvertToRGBA(referenceImg)

	w := refImgRGBA.Bounds().Dx()
	h := refImgRGBA.Bounds().Dy()

	var population []*Candidate

	for i := 0; i < PopulationCount; i++ {
		c := RandomCandidate(w, h)
		evaluateCandidate(c, refImgRGBA)
		population = append(population, c)
	}

	for i := 0; i < maxGen; i++ {
		log.Printf("generation %d", i)

		// after sort, the 2 best populations will be at [0] and [1], worst will be at [len() - 1]
		sort.Sort(ByFitness(population))
		for _, candidate := range population {
			log.Print(candidate)
		}

		offspring := population[0].Mate(population[1])
		evaluateCandidate(offspring, refImgRGBA)

		// evict the least fit individual
		leastFit := population[len(population)-1]
		if leastFit.Fitness > offspring.Fitness {
			population[len(population)-1] = offspring
			log.Printf("evicted, fitness: %d -> %d", leastFit.Fitness, offspring.Fitness)
		} else {
			log.Printf("preserved, fitness: %d vs %d", leastFit.Fitness, offspring.Fitness)
		}
		population[0].DrawAndSave(destFile)
		safeImage.Update(population[len(population) - 1].img)
	}
	//log.Printf("population: %+v", population)
}

func evaluateCandidate(c *Candidate, referenceImg *image.RGBA) {
	diff, err := Compare(referenceImg, c.img)
	//diff, err := FastCompare(referenceImg, c.img)

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
