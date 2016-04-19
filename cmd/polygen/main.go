package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/armhold/polygen"
)

var (
	maxGen     int
	polyCount  int
	srcImgFile string
	dstImgFile string
	cpArg string
	host, port string
)


func init() {
	flag.IntVar(&maxGen, "max", 100000, "the number of generations")
	flag.IntVar(&polyCount, "poly", 50, "the number of polygons")
	flag.StringVar(&srcImgFile, "source", "images/mona_lisa.jpg", "the source input image file")
	flag.StringVar(&dstImgFile, "dest", "output.png", "the output image file")
	flag.StringVar(&cpArg, "cp", "", "checkpoint file")
	flag.StringVar(&host, "host", "localhost", "which hostname to http listen on")
	flag.StringVar(&port, "port", "8080", "which port to http listen on")

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
	refImg := polygen.MustReadImage(srcImgFile)

	// a set of thread-safe images that can be updated by the evolver, and displayed via the web
	var previews []*polygen.SafeImage

	totalImages := polygen.PopulationCount
	placeholder := refImg.Bounds()
	for i := 0; i < totalImages; i++ {
		img := polygen.NewSafeImage(placeholder)
		previews = append(previews, img)
	}

	go polygen.Serve(host+":"+port, refImg, previews)

	cp := polygen.DeriveCheckpointFile(srcImgFile, cpArg, polyCount)

	evolver, err := polygen.NewEvolver(refImg, dstImgFile, cp)
	if err != nil {
		log.Fatal(err)
	}

	evolver.Run(maxGen, polyCount, previews)
}

