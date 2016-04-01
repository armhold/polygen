package polygen

import "math/rand"

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
