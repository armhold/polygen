package polygen

import (
	"image"
	"log"
	"sort"
	"math/rand"
	"time"
	"sync"
	"io/ioutil"
	"fmt"
	"encoding/gob"
	"bytes"
)

type Evolver struct {
	refImgRGBA *image.RGBA
	dstImgFile string
	previews   []*SafeImage
	checkpoint string
	candidates []*Candidate
}


func NewEvolver(refImg image.Image, dstImageFile string, checkpoint string) *Evolver {
	result := &Evolver{
		dstImgFile: dstImageFile,
		checkpoint: checkpoint,
	}

	result.refImgRGBA = ConvertToRGBA(refImg)

	return result
}


func (e *Evolver) RestoreSavedCandidates(checkpoint string) error {
	b, err := ioutil.ReadFile(checkpoint)
	if err != nil {
		return fmt.Errorf("error reading candidates file: %s", err)
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(b))

	var candidates []*Candidate
	err = decoder.Decode(&candidates)
	if err != nil {
		return fmt.Errorf("error decoding candidates: %s", err)
	}

	e.candidates = candidates

	// TODO: could parallelize this, but probably not worth the code
	for _, c := range e.candidates {
		c.RenderImage()
		e.evaluateCandidate(c)
	}

	return nil
}

func (e *Evolver) saveCheckpoint() error {
	log.Printf("checkpointing candidates to %s...", e.checkpoint)

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(e.candidates)
	if err != nil {
		return fmt.Errorf("error encoding candidates: %s", err)
	}

	err = ioutil.WriteFile(e.checkpoint, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing candidates to file: %s", err)
	}

	return nil
}

func (e *Evolver) Run(maxGen int, previews []*SafeImage) {
	w := e.refImgRGBA.Bounds().Dx()
	h := e.refImgRGBA.Bounds().Dy()

	var mostFit *Candidate
	if len(e.candidates) > 0 {
		mostFit = e.candidates[0]
	} else {
		e.candidates = make([]*Candidate, PopulationCount)
		mostFit = RandomCandidate(w, h)
		e.candidates[0] = mostFit
	}

	e.evaluateCandidate(mostFit)

	startTime := time.Now()
	generationsSinceChange := 0

	for gen := 0; gen < maxGen; gen++ {
		c := make(chan *Candidate, PopulationCount)
		var wg sync.WaitGroup

		processCandidate := func(cand *Candidate) {
			cand.MutateInPlace()
			e.evaluateCandidate(cand)
			c <- cand
			wg.Done()
		}

		wg.Add(PopulationCount - 1)

		for i := 1; i < PopulationCount; i++ {
			copy := mostFit.CopyOf()
			go processCandidate(copy)
		}

		wg.Wait()

		for i := 1; i < PopulationCount; i++ {
			e.candidates[i] = <- c
		}

		// after sort, the best will be at [0], worst will be at [len() - 1]
		sort.Sort(ByFitness(e.candidates))

		if gen % 10 == 0 {
			printStats(e.candidates, gen, generationsSinceChange, startTime)
		}

		for i := 0; i < len(previews); i++ {
			previews[i].Update(e.candidates[i].img)
		}

		currBest := e.candidates[0]

		if currBest.Fitness < mostFit.Fitness {
			generationsSinceChange = 0
			mostFit = currBest
		} else {
			generationsSinceChange++
		}

		if gen % 500 == 0 {
			err := e.saveCheckpoint()
			if err != nil {
				log.Fatalf("error saving checkpoint file: %s", err)
			}
		}
	}

	mostFit.DrawAndSave(e.dstImgFile)
	log.Printf("after %d generations, fitness is: %d, saved to %s", maxGen, mostFit.Fitness, e.dstImgFile)
}


func (e *Evolver) evaluateCandidate(c *Candidate) {
//	diff, err := Compare(e.refImgRGBA, c.img)
	diff, err := FastCompare(e.refImgRGBA, c.img)

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
