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
	"os"
)

// Evolver uses a genetic algorithm to evolve a set of polygons to approximate an image.
type Evolver struct {
	refImgRGBA             *image.RGBA
	dstImgFile             string
	checkPointFile         string
	candidates             []*Candidate
	mostFit                *Candidate
	generation             int
	generationsSinceChange int
}

// Checkpoint is used for serializing the current best candidate and corresponding generation
// count to a checkpoint file.
type Checkpoint struct {
	Generation             int
	GenerationsSinceChange int
	MostFit                *Candidate
}

func NewEvolver(refImg image.Image, dstImageFile string, checkPointFile string) (*Evolver, error) {
	result := &Evolver{
		dstImgFile:     dstImageFile,
		checkPointFile: checkPointFile,
		candidates:     make([]*Candidate, PopulationCount),
	}

	result.refImgRGBA = ConvertToRGBA(refImg)

	// if there's an existing checkpoint file, restore from last checkpoint
	if _, err := os.Stat(checkPointFile); !os.IsNotExist(err) {
		err := result.restoreFromCheckpoint()
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Run runs the Evolver until maxGen generations have been evaluated.
// At each generation, the candidate images are rendered & evaluated, and the preview images are
// updated to reflect the current state.
func (e *Evolver) Run(maxGen, polyCount int, previews []*SafeImage) {
	w := e.refImgRGBA.Bounds().Dx()
	h := e.refImgRGBA.Bounds().Dy()

	// no candidate from prev call to RestoreFromCheckpoint()
	if e.mostFit == nil {
		e.mostFit = RandomCandidate(w, h, polyCount)
		e.candidates[0] = e.mostFit
	}

	e.renderAndEvaluate(e.mostFit)

	startTime := time.Now()

	for ; e.generation < maxGen; e.generation++ {
		c := make(chan *Candidate, PopulationCount)
		var wg sync.WaitGroup

		processCandidate := func(cand *Candidate) {
			for i := 0; i < 3; i++ {
				cand.mutateInPlace()
			}

			e.renderAndEvaluate(cand)
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
			e.printStats(startTime)
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
			err := e.mostFit.drawAndSave(e.dstImgFile)
			if err != nil {
				log.Fatalf("error saving output image: %s", err)
			}

			if e.checkPointFile != "" {
				err = e.saveCheckpoint()
				if err != nil {
					log.Fatalf("error saving checkpoint file: %s", err)
				}
			}
			dur := time.Since(cpSave)
			log.Printf("checkpoint took %s", dur)
		}
	}

	e.mostFit.drawAndSave(e.dstImgFile)
	log.Printf("after %d generations, fitness is: %d, saved to %s", maxGen, e.mostFit.Fitness, e.dstImgFile)
}

func (e *Evolver) restoreFromCheckpoint() error {
	b, err := ioutil.ReadFile(e.checkPointFile)
	if err != nil {
		return fmt.Errorf("error reading checkpoint file: %s: %s", e.checkPointFile, err)
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(b))

	var cp Checkpoint
	err = decoder.Decode(&cp)
	if err != nil {
		return fmt.Errorf("error decoding checkpoint file: %s %s", e.checkPointFile, err)
	}

	e.generation = cp.Generation
	e.generationsSinceChange = cp.GenerationsSinceChange
	e.candidates[0] = cp.MostFit
	e.mostFit = cp.MostFit
	e.renderAndEvaluate(e.mostFit)

	return nil
}

func (e *Evolver) saveCheckpoint() error {
	log.Printf("checkpointing to %s...", e.checkPointFile)

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

	err = ioutil.WriteFile(e.checkPointFile, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing checkpoint to file: %s", err)
	}

	return nil
}

func (e *Evolver) renderAndEvaluate(c *Candidate) {
	c.renderImage()

	diff, err := FastCompare(e.refImgRGBA, c.img)

	if err != nil {
		log.Fatalf("error comparing images: %s", err)
	}

	c.Fitness = diff
}

func (e *Evolver) printStats(startTime time.Time) {
	dur := time.Since(startTime)
	best := e.candidates[0].Fitness
	worst := e.candidates[len(e.candidates)-1].Fitness

	log.Printf("dur: %s, gen: %d, since change: %d, best: %d, worst: %d", dur, e.generation, e.generationsSinceChange, best, worst)
}
