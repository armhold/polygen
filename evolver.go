package polygen

import (
	"log"
	"image"
	"sort"
)

type Individual interface {
	Fitness() int
	BreedWith(Individual) Individual
}

func Evolve(maxGen int, sourceFile, destFile string) {
	referenceImg := MustReadImage(sourceFile)
	w := referenceImg.Bounds().Dx()
	h := referenceImg.Bounds().Dy()

	var population []*PolygonSet

	for i := 0; i < PopulationCount; i++ {
		individual := &PolygonSet{}
		for j := 0; j < PolygonsPerIndividual; j++ {
			individual.Polygons = append(individual.Polygons, RandomPolygon(w, h))
		}

		population = append(population, individual)
	}

	population[0].DrawAndSave(w, h, destFile)


	for i := 0; i < maxGen; i++ {
		log.Printf("generation %d", i)

		evaluatePopulation(population, referenceImg)
		sort.Sort(ByFitness(population))

		for _, polygonSet := range population {
			log.Print(polygonSet)
		}
	}
	//log.Printf("population: %+v", population)
}


func evaluatePopulation(population []*PolygonSet, referenceImg image.Image) {
	w := referenceImg.Bounds().Dx()
	h := referenceImg.Bounds().Dy()

	for _, polygonSet := range population {
		img := polygonSet.RenderImage(w, h)
		diff, err := Compare(referenceImg, img)

		if err != nil {
			log.Fatalf("error comparing images: %s", err)
		}

		polygonSet.Fitness = diff
	}
}
