package polygen

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	templates *template.Template
)

func init() {
	// use https://github.com/jteeuwen/go-bindata to re-construct templates/index.html.
	// This is done so that the polygen executable can be installed anywhere, and not have
	// to depend on an external asset like templates/index.html.
	data, err := Asset("templates/index.html")
	if err != nil {
		log.Fatalf("unable to read templates/index.html from bindata: %s", err)
	}

	templates, err = template.New("templates").New("index.html").Parse(string(data))
	if err != nil {
		log.Fatalf("error parsing bindata for templates/index.html: %s", err)
	}
}

type Page struct {
	// TODO
}

func rootHandler(previewCount int) http.HandlerFunc {
	// TODO: use previewCount in template so we don't have to hard-code 10 img tags

	return func(w http.ResponseWriter, r *http.Request) {
		p := &Page{}
		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			log.Println(err)
		}
	}
}

func imageHandler(previews []*SafeImage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := SplitPath(r.URL.Path)

		if len(p) != 2 {
			err := "missing image argument"
			log.Print(err)
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		imageNum, err := strconv.Atoi(p[1])
		if err != nil {
			err := "bad image argument"
			log.Print(err)
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		if imageNum >= len(previews) {
			err := "bad image argument"
			log.Print(err)
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		img := previews[imageNum].Value()
		serveNonCacheableImage(img, w, r)
	}
}

func refImageHandler(referenceImg image.Image) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveNonCacheableImage(referenceImg, w, r)
	}
}

func serveNonCacheableImage(img image.Image, w http.ResponseWriter, r *http.Request) {
	buffer := new(bytes.Buffer)

	// serve images without compression- this results in a 5-10x speedup, and helps get more CPU
	// to the evolver. If you're actually serving images over the network (vs just localhost),
	// you might want to change this to png.DefaultCompression.
	encoder := png.Encoder{CompressionLevel: png.NoCompression}

	if err := encoder.Encode(buffer, img); err != nil {
		s := fmt.Sprintf("unable to encode image: %s", err)
		log.Print(s)
		http.Error(w, s, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-control", "max-age=0, must-revalidate")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image")
	}
}

func Serve(hostPort string, refImg image.Image, previews []*SafeImage) {
	http.Handle("/", rootHandler(len(previews)))
	http.Handle("/image/", imageHandler(previews))
	http.Handle("/ref", refImageHandler(refImg))

	log.Printf("listening on %s...", hostPort)

	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func SplitPath(path string) []string {
	trimmed := strings.TrimFunc(path, func(r rune) bool {
		return r == '/'
	})

	return strings.Split(trimmed, "/")
}
