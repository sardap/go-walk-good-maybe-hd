package utility

import (
	"math/rand"
	"time"

	"github.com/sardap/walk-good-maybe-hd/math"
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

func ClampVec2(val, min, max math.Vector2) math.Vector2 {
	return math.Vector2{
		X: ClampFloat64(val.X, min.X, max.X),
		Y: ClampFloat64(val.Y, min.Y, max.Y),
	}
}

func RandRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandRangeFloat64(min, max int) float64 {
	return float64(rand.Intn(max-min)+min) + rand.Float64()
}
