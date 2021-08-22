package utility

import (
	"math/rand"
	"time"
)

func DeltaToDuration(dt float32) time.Duration {
	var i int = int(dt * float32(time.Second))
	return time.Duration(i)
}

func WrapInt(x, min, max int) int {
	if x >= max {
		return x + min - max
	} else if x < min {
		return x + max - min
	}
	return x
}

func ClampFloat64(x, min, max float64) float64 {
	if x > max {
		return max
	}

	if x < min {
		return min
	}

	return x
}

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min

}
