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
	refImgRGBA             *image.RGBA
	dstImgFile             string
	previews               []*SafeImage
	checkpoint             string
	candidates             []*Candidate
	mostFit                *Candidate
	generation             int
	generationsSinceChange int
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

type Checkpoint struct {
	Generation             int
	GenerationsSinceChange int
	MostFit                *Candidate
}

func (e *Evolver) RestoreFromCheckpoint(checkpoint string) error {
	b, err := ioutil.ReadFile(checkpoint)
	if err != nil {
		return fmt.Errorf("error reading checkpoint file: %s", err)
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(b))

	var cp Checkpoint
	err = decoder.Decode(&cp)
	if err != nil {
		return fmt.Errorf("error decoding checkpoint: %s", err)
	}

	e.generation = cp.Generation
	e.generationsSinceChange = cp.GenerationsSinceChange
	e.candidates[0] = cp.MostFit
	e.mostFit = cp.MostFit

	e.mostFit.RenderImage()
	e.evaluateCandidate(e.mostFit)

	return nil
}

func (e *Evolver) saveCheckpoint() error {
	log.Printf("checkpointing to %s...", e.checkpoint)

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	cp := &Checkpoint{
		Generation:             e.generation,
		GenerationsSinceChange: e.generationsSinceChange,
		MostFit:                e.mostFit,
	}

	err := encoder.Encode(cp)
	if err != nil {
		return fmt.Errorf("error encoding checkpoint: %s", err)
	}

	err = ioutil.WriteFile(e.checkpoint, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing checkpoint to file: %s", err)
	}

	return nil
}

// Run creates a
func (e *Evolver) Run(maxGen, polyCount int, previews []*SafeImage) {
	w := e.refImgRGBA.Bounds().Dx()
	h := e.refImgRGBA.Bounds().Dy()

	// no candidate from prev call to RestoreSavedCandidate()
	if e.mostFit == nil {
		e.mostFit = RandomCandidate(w, h, polyCount)
		e.candidates[0] = e.mostFit
	}

	e.evaluateCandidate(e.mostFit)

	startTime := time.Now()

	for ; e.generation < maxGen; e.generation++ {
		c := make(chan *Candidate, PopulationCount)
		var wg sync.WaitGroup

		processCandidate := func(cand *Candidate) {
			cand.MutateInPlace()
			e.evaluateCandidate(cand)
			c <- cand
			wg.Done()
		}

		wg.Add(PopulationCount - 1)

		// mostFit is already in slot 0
		for i := 1; i < PopulationCount; i++ {
			copy := e.mostFit.CopyOf()
			go processCandidate(copy)
		}

		wg.Wait()

		for i := 1; i < PopulationCount; i++ {
			e.candidates[i] = <-c
		}

		// after sort, the best will be at [0], worst will be at [len() - 1]
		sort.Sort(ByFitness(e.candidates))

		if e.generation%10 == 0 {
			printStats(e.candidates, e.generation, e.generationsSinceChange, startTime)
		}

		for i := 0; i < len(previews); i++ {
			previews[i].Update(e.candidates[i].img)
		}

		currBest := e.candidates[0]

		if currBest.Fitness < e.mostFit.Fitness {
			e.generationsSinceChange = 0
			e.mostFit = currBest
		} else {
			e.generationsSinceChange++
		}

		if e.generation%250 == 0 {
			cpSave := time.Now()
			err := e.mostFit.DrawAndSave(e.dstImgFile)
			if err != nil {
				log.Fatalf("error saving output image: %s", err)
			}

			if e.checkpoint != "" {
				err = e.saveCheckpoint()
				if err != nil {
					log.Fatalf("error saving checkpoint file: %s", err)
				}
			}
			dur := time.Since(cpSave)
			log.Printf("checkpoint took %s", dur)
		}
	}

	e.mostFit.DrawAndSave(e.dstImgFile)
	log.Printf("after %d generations, fitness is: %d, saved to %s", maxGen, e.mostFit.Fitness, e.dstImgFile)
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
