package global

import (
	"fmt"
)

func Version() string {
	return fmt.Sprintf("ezorm v%d.%d.%d", vMajor, vMinor, vPatch)
}

const (
	vMajor = 0
	vMinor = 0
	vPatch = 2
)
