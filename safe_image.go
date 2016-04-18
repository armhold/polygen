package polygen

import (
	"image"
	"sync"
)

// SafeImage is an image protected by a RWMutex. All access should be via the Update() and Value() methods.
// This allows us to update the candidate images from a goroutine in the Evolver, and display them in the Server.
type SafeImage struct {
	img image.Image
	mux sync.RWMutex
}

func NewSafeImage(img image.Image) *SafeImage {
	return &SafeImage{img: img}
}

func (s *SafeImage) Update(img image.Image) {
	s.mux.Lock()
	s.img = img
	s.mux.Unlock()
}

func (s *SafeImage) Value() image.Image {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.img
}
