package polygen

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"sort"
	"sync"
	"time"
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
		candidates: make([]*Candidate, PopulationCount),
	}

	result.refImgRGBA = ConvertToRGBA(refImg)

	return result
}

func (e *Evolver) RestoreSavedCandidate(checkpoint string) error {
	b, err := ioutil.ReadFile(checkpoint)
	if err != nil {
		return fmt.Errorf("error reading candidates file: %s", err)
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(b))

	var candidate Candidate
	err = decoder.Decode(&candidate)
	if err != nil {
		return fmt.Errorf("error decoding candidate: %s", err)
	}

	candidate.RenderImage()
	e.evaluateCandidate(&candidate)
	e.candidates[0] = &candidate

	return nil
}

func (e *Evolver) saveCheckpoint() error {
	log.Printf("checkpointing candidate to %s...", e.checkpoint)

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(e.candidates[0])
	if err != nil {
		return fmt.Errorf("error encoding candidate: %s", err)
	}

	err = ioutil.WriteFile(e.checkpoint, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing candidate to file: %s", err)
	}

	return nil
}


// Run creates a
func (e *Evolver) Run(maxGen, polyCount int, previews []*SafeImage) {
	w := e.refImgRGBA.Bounds().Dx()
	h := e.refImgRGBA.Bounds().Dy()

	// may already have a candidate from prev call to RestoreSavedCandidate()
	mostFit := e.candidates[0]

	if mostFit == nil {
		mostFit = RandomCandidate(w, h, polyCount)
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
			e.candidates[i] = <-c
		}

		// after sort, the best will be at [0], worst will be at [len() - 1]
		sort.Sort(ByFitness(e.candidates))

		if gen%10 == 0 {
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

		if gen%500 == 0 {
			err := mostFit.DrawAndSave(e.dstImgFile)
			if err != nil {
				log.Fatalf("error saving output image: %s", err)
			}

			if e.checkpoint != "" {
				err = e.saveCheckpoint()
				if err != nil {
					log.Fatalf("error saving checkpoint file: %s", err)
				}
			}
		}
	}

	mostFit.DrawAndSave(e.dstImgFile)
	log.Printf("after %d generations, fitness is: %d, saved to %s", maxGen, mostFit.Fitness, e.dstImgFile)
}

func (e *Evolver) evaluateCandidate(c *Candidate) {
	diff, err := FastCompare(e.refImgRGBA, c.img)

	if err != nil {
		log.Fatalf("error comparing images: %s", err)
	}

	c.Fitness = diff
}

func printStats(sortedPop []*Candidate, generations, generationsSinceChange int, startTime time.Time) {
	dur := time.Since(startTime)
	best := sortedPop[0].Fitness
	worst := sortedPop[len(sortedPop)-1].Fitness

	log.Printf("dur: %s, gen: %d, since change: %d, best: %d, worst: %d", dur, generations, generationsSinceChange, best, worst)
}
