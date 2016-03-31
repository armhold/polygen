package polygen

type Individual interface {
	Fitness() int
	BreedWith(Individual) Individual
}


func Evolve(sourceFile, destFile string) {
	var population []PolygonSet

	for i := 0; i < PopulationCount; i++ {
		var individual PolygonSet
		for j := 0; j < PolygonsPerIndividual; j++ {
			individual = append(individual, RandomPolygon())
		}

		population = append(population, individual)
	}

	population[0].DrawAndSave(destFile)

	//log.Printf("population: %+v", population)
}

