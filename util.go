package polygen

import "math/rand"

// RandomInt uses the default random Source to return a random integer that is >= min, but < max
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandomBool uses the default random Source to return either true or false.
func RandomBool() bool {
	i := rand.Intn(2)
	if i == 0 {
		return false
	}

	return true
}
