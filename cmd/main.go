package main

import (
	"flag"
	"github.com/armhold/polygen"
	"math/rand"
	"os"
	"time"
)

var (
	maxGen               int
	sourceFile, destFile string
)

func init() {
	flag.IntVar(&maxGen, "maxgen", 10000, "the number of generations")
	flag.StringVar(&sourceFile, "source", "", "the source input image file")
	flag.StringVar(&destFile, "dest", "output.png", "the output image file")
	flag.Parse()

	if sourceFile == "" || destFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	polygen.Evolve(maxGen, sourceFile, destFile)
}
