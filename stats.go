package polygen

import (
	"log"
	"time"
)

// Stats is used by the Evolver for printing runtime statistics.
type Stats struct {
	startTime           time.Time
	prevTime            time.Time
	candidatesEvaluated int
}

func NewStats() *Stats {
	timeNow := time.Now()

	return &Stats{startTime: timeNow, prevTime: timeNow}
}

// Increments the number of candidates that have been evaluated since last call to Print().
func (s *Stats) Increment(count int) {
	s.candidatesEvaluated += count
}

func (s *Stats) Print(best, worst *Candidate, generation, generationsSinceChange int) {
	timeNow := time.Now()
	durOverall := timeNow.Sub(s.startTime)

	durSinceLastStats := time.Since(s.prevTime)
	cps := float64(s.candidatesEvaluated) / durSinceLastStats.Seconds()

	s.prevTime = timeNow
	s.candidatesEvaluated = 0

	log.Printf("dur: %s, gen: %d, since change: %d, candidates/sec: %.2f, best: %d, worst: %d", durOverall, generation, generationsSinceChange, cps, best.Fitness, worst.Fitness)
}
