package polygen

import (
	"fmt"
	"math/rand"
	"path"
	"path/filepath"
	"strings"
)

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

// try to make it painless for users to automatically get a checkpoint file.
// We want the file to be tied both the source file name and the polygon count.
// If an explicit cpArg is given, we just use that.
func DeriveCheckpointFile(sourceFile, cpArg string, polyCount int) string {
	if cpArg != "" {
		return cpArg
	}

	basename := path.Base(sourceFile)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	return fmt.Sprintf("%s-%d-checkpoint.tmp", name, polyCount)
}
