package main

import (
	"image"
	"image/color"
	"log"

	"flag"
	"github.com/armhold/polygen"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image/draw"
	"math/rand"
	"os"
	"time"
)

var (
	maxGen     int
	polyCount  int
	srcImgFile string
	dstImgFile string
	cpArg string
)

func init() {
	flag.IntVar(&maxGen, "max", 100000, "the number of generations")
	flag.IntVar(&polyCount, "poly", 50, "the number of polygons")
	flag.StringVar(&srcImgFile, "source", "images/mona_lisa.jpg", "the source input image file")
	flag.StringVar(&dstImgFile, "dest", "output.png", "the output image file")
	flag.StringVar(&cpArg, "cp", "", "checkpoint file")

	flag.Parse()

	if srcImgFile == "" || dstImgFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	refImg := polygen.MustReadImage(srcImgFile)

	cp := polygen.DeriveCheckpointFile(srcImgFile, cpArg, polyCount)

	evolver, err := polygen.NewEvolver(refImg, dstImgFile, cp)
	if err != nil {
		log.Fatal(err)
	}

	driver.Main(func(s screen.Screen) {
		// double-width: we will display the refImg and the best candidate side-by-side
		winSize := image.Point{refImg.Bounds().Dx() * 2, refImg.Bounds().Dy()}
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: winSize.X, Height: winSize.Y})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		b, err := s.NewBuffer(winSize)
		if err != nil {
			log.Fatal(err)
		}
		defer b.Release()

		t, err := s.NewTexture(winSize)
		if err != nil {
			log.Fatal(err)
		}
		defer t.Release()
		t.Upload(image.Point{}, b, b.Bounds())

		// draw the refImg
		draw.Draw(b.RGBA(), b.Bounds(), refImg, refImg.Bounds().Min, draw.Src)

		// start the evolver; it will send paint events to the Window when best candidate has changed
		go evolver.Run(maxGen, polyCount, w)

		log.Printf("refImage bounds: %v\n", refImg.Bounds())
		log.Printf("winSize        : %v\n", winSize)

		var sz size.Event
		for {
			e := w.NextEvent()

			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case polygen.CandidateChangedEvent:
				ccEvent := polygen.CandidateChangedEvent(e)
				img := ccEvent.Image

				// offset next to the refImg, then paint the candidate
				bounds := b.Bounds()
				r := image.Rect(bounds.Min.X+refImg.Bounds().Dx(), bounds.Min.Y, bounds.Max.X, bounds.Max.Y)
				draw.Draw(b.RGBA(), r, img, img.Bounds().Min, draw.Src)
				doPaint(w, b, sz)

			case paint.Event:
				doPaint(w, b, sz)

			case size.Event:
				sz = e
				log.Printf("size: %v\n", sz)

			case error:
				log.Print(e)
			}

		}
	})
}

func doPaint(w screen.Window, b screen.Buffer, sz size.Event) {
	w.Fill(sz.Bounds(), color.Black, screen.Src)
	w.Upload(image.Point{}, b, b.Bounds())
	w.Publish()
}
