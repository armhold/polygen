package polygen

import "math/rand"

// RandomInt returns a random integer that is >= min, but < max
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func NextBool() bool {
	i := rand.Intn(2)
	if i == 0 {
		return false
	}

	return true
}
