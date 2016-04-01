package polygen

import (
	"log"
	"net/http"
	"html/template"
	"bytes"
	"strconv"
	"image/png"
	"image"
)

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}


type Page struct {
	 // nothing here yet
}



func rootHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{}
	if err := templates.ExecuteTemplate(w, "index.html", p) ; err != nil {
		log.Println(err)
	}
}

func evolvingImageHandler(safe *SafeImage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		img := safe.Value()
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
		log.Println("unable to encode image.")
	}

	w.Header().Set("Cache-control", "max-age=0, must-revalidate")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image")
	}
}



func Serve(hostPort string, refImg image.Image, evolvingImg *SafeImage) {
	http.HandleFunc("/", rootHandler)
	http.Handle("/image", evolvingImageHandler(evolvingImg))
	http.Handle("/ref", refImageHandler(refImg))

	log.Printf("listening on %s...", hostPort)

	err := http.ListenAndServe(hostPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
