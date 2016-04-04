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
	maxGen int
	polyCount int
	srcImgFile string
	dstImgFile string
	host, port string
	loadFromCheckpoint string
	saveToCheckpoint string
)

func init() {
	flag.IntVar(&maxGen, "max", 100000, "the number of generations")
	flag.IntVar(&polyCount, "poly", 50, "the number of polygons")
	flag.StringVar(&srcImgFile, "source", "", "the source input image file")
	flag.StringVar(&dstImgFile, "dest", "output.png", "the output image file")
	flag.StringVar(&host, "host", "localhost", "which hostname to http listen on")
	flag.StringVar(&port, "port", "8080", "which port to http listen on")
	flag.StringVar(&loadFromCheckpoint, "load", "", "load from checkpoint file")
	flag.StringVar(&saveToCheckpoint, "save", "candidates.tmp", "save to checkpoint file")

	flag.Parse()

	if srcImgFile == "" || dstImgFile == "" {
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

	// a set of thread-safe images that can be updated by the evolver, and displayed via the web
	var previews []*polygen.SafeImage

	totalImages := polygen.PopulationCount
	placeholder := image.Rect(0, 0, 200, 200)
	for i := 0; i < totalImages; i++ {
		img := &polygen.SafeImage{Image: placeholder}
		previews = append(previews, img)
	}

	go polygen.Serve(host + ":" + port, refImg, previews)

	evolver := polygen.NewEvolver(refImg, dstImgFile, saveToCheckpoint)

	if loadFromCheckpoint != "" {
		err := evolver.RestoreSavedCandidates(loadFromCheckpoint)
		if err != nil {
			log.Fatalf("error restoring candidates from checkpoint: %s", err)
		}
	}

	evolver.Run(maxGen, polyCount, previews)
}
