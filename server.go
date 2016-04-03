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
)

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}


type Page struct {
	 ImageCount int
}

func rootHandler(safeImages []*SafeImage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &Page{ImageCount: len(safeImages)}
		if err := templates.ExecuteTemplate(w, "index.html", p) ; err != nil {
			log.Println(err)
		}
	}
}

func evolvingImageHandler(safeImages []*SafeImage) http.HandlerFunc {
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

		if imageNum >= len(safeImages) {
			err := "bad image argument"
			log.Print(err)
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		img := safeImages[imageNum].Value()
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



func Serve(hostPort string, refImg image.Image, evolvingImages []*SafeImage) {
	http.Handle("/", rootHandler(evolvingImages))
	http.Handle("/image/", evolvingImageHandler(evolvingImages))
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
