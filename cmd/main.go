package main

import (
	"time"
	"math/rand"
	"github.com/armhold/polygon"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	polygon.Evolve()
}
