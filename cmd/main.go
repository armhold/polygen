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
	host, port string
)

func init() {
	flag.IntVar(&maxGen, "maxgen", 10000, "the number of generations")
	flag.StringVar(&sourceFile, "source", "", "the source input image file")
	flag.StringVar(&destFile, "dest", "output.png", "the output image file")
	flag.StringVar(&host, "host", "localhost", "which hostname to http listen on")
	flag.StringVar(&port, "port", "8080", "which port to http listen on")

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
	go polygen.Serve(host + ":" + port)
	polygen.Evolve(maxGen, sourceFile, destFile)
}
