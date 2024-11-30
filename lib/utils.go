package lib

import (
	"math/rand/v2"
)

func RandomDigitGenerate(min int, max int) int {
	return rand.IntN(max-min) + min
}
