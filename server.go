package polygen

import (
	"log"
	"net/http"
	"html/template"
	"bytes"
	"strconv"
	"image/png"
	"image"
	"strings"
	"fmt"
)

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

type Page struct {
	// TODO
}

func rootHandler(previewCount int) http.HandlerFunc {
	// TODO: use previewCount in template so we don't have to hard-code 10 img tags

	return func(w http.ResponseWriter, r *http.Request) {
		p := &Page{}
		if err := templates.ExecuteTemplate(w, "index.html", p) ; err != nil {
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
	if err := png.Encode(buffer, img); err != nil {
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
