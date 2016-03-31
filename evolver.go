package polygen

import (
	"log"
	"image"
)

type Individual interface {
	Fitness() int
	BreedWith(Individual) Individual
}


func Evolve(maxGen int, sourceFile, destFile string) {
	referenceImg := MustReadImage(sourceFile)

	var population []PolygonSet

	for i := 0; i < PopulationCount; i++ {
		var individual PolygonSet
		for j := 0; j < PolygonsPerIndividual; j++ {
			individual.Polygons = append(individual.Polygons, RandomPolygon())
		}

		population = append(population, individual)
	}

	population[0].DrawAndSave(destFile)


	for i := 0; i < maxGen; i++ {
		log.Printf("generation %d", i)

		evaluatePopulation(population, referenceImg)

	}
	//log.Printf("population: %+v", population)
}


func evaluatePopulation(population []PolygonSet, referenceImg image.Image) {
	for _, polygonSet := range population {
		img := polygonSet.RenderImage()
		diff, err := Compare(referenceImg, img)

		if err != nil {
			log.Fatalf("error comparing images: %s", err)
		}

		polygonSet.Fitness = diff
	}
}
