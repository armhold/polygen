package main

import (
	"flag"
	"github.com/armhold/polygen"
	"math/rand"
	"os"
	"time"
	"image"
	"log"
)

var (
	maxGen               int
	sourceFile, destFile string
	host, port string
	loadFrom, saveTo string
)

func init() {
	flag.IntVar(&maxGen, "maxgen", 10000, "the number of generations")
	flag.StringVar(&sourceFile, "source", "", "the source input image file")
	flag.StringVar(&destFile, "dest", "output.png", "the output image file")
	flag.StringVar(&host, "host", "localhost", "which hostname to http listen on")
	flag.StringVar(&port, "port", "8080", "which port to http listen on")
	flag.StringVar(&loadFrom, "load", "", "load from checkpoint file")
	flag.StringVar(&saveTo, "save", "candidates.tmp", "save to checkpoint file")

	flag.Parse()

	if sourceFile == "" || destFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if port == "" {
		port = "8080"
	}
	if host == "" {
		host = "localhost"
	}

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	refImg := polygen.MustReadImage("images/mona_lisa.jpg")

	var safeImages []*polygen.SafeImage

	totalImages := polygen.PopulationCount
	for i := 0; i < totalImages; i++ {
		img := &polygen.SafeImage{Image: image.Rect(0, 10, 10, 10)}
		safeImages = append(safeImages, img)
	}

	go polygen.Serve(host + ":" + port, refImg, safeImages)

	evolver := polygen.NewEvolver(refImg, destFile, saveTo)

	if loadFrom != "" {
		err := evolver.RestoreSavedCandidates(loadFrom)
		if err != nil {
			log.Fatalf("error restoring candidates from checkpoint: %s", err)
		}
	}

	evolver.Run(maxGen, safeImages)
}
