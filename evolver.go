package polygen

import (
	"image"
	"log"
	"sort"
)

type Individual interface {
	Fitness() int
	BreedWith(Individual) Individual
}

func Evolve(maxGen int, sourceFile, destFile string) {
	referenceImg := ConvertToRGBA(MustReadImage(sourceFile))

	w := referenceImg.Bounds().Dx()
	h := referenceImg.Bounds().Dy()

	var population []*Candidate

	for i := 0; i < PopulationCount; i++ {
		population = append(population, RandomCandidate(w, h))
	}

	for i := 0; i < maxGen; i++ {
		log.Printf("generation %d", i)

		evaluatePopulation(population, referenceImg)

		// after sort, the 2 best populations will be at [0] and [1], worst will be at [len() - 1]
		sort.Sort(ByFitness(population))
		for _, candidate := range population {
			log.Print(candidate)
		}

		offspring := population[0].Mate(population[1])

		// evict the least fit individual
		population[len(population)-1] = offspring

		population[0].DrawAndSave(destFile)
	}
	//log.Printf("population: %+v", population)
}

func evaluatePopulation(population []*Candidate, referenceImg *image.RGBA) {
	for _, candidate := range population {
		diff, err := Compare(referenceImg, candidate.img)

		if err != nil {
			log.Fatalf("error comparing images: %s", err)
		}

		candidate.Fitness = diff
	}
}
