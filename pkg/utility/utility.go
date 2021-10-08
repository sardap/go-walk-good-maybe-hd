package utility

import (
	"math/rand"
	"time"

	"github.com/sardap/walk-good-maybe-hd/pkg/math"
)

func DeltaToDuration(dt float32) time.Duration {
	var i int = int(dt * float32(time.Second))
	return time.Duration(i)
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func WrapInt(x, min, max int) (result int) {
	if x >= max {
		result = x + min - max

	} else if x < min {
		result = x + max - min
	} else {
		result = x
	}

	if result > max || result < min {
		result = WrapInt(result, min, max)
	}
	return result
}

func WrapFloat64(x, min, max float64) float64 {
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

func WrapVec2(val, min, max math.Vector2) math.Vector2 {
	return math.Vector2{
		X: WrapFloat64(val.X, min.X, max.X),
		Y: WrapFloat64(val.Y, min.Y, max.Y),
	}
}

func RandRange(rand *rand.Rand, min, max int) int {
	return rand.Intn(max-min) + min
}

func RandRangeFloat64(rand *rand.Rand, min, max int) float64 {
	return float64(rand.Intn(max-min)+min) + rand.Float64()
}

func ContainsString(ary []string, tags ...string) bool {
	for _, otherTag := range ary {
		for _, sTag := range tags {
			if otherTag == sTag {
				return true
			}
		}
	}

	return false
}

func ContainsInt(ary []int, tags ...int) bool {
	for _, otherTag := range ary {
		for _, sTag := range tags {
			if otherTag == sTag {
				return true
			}
		}
	}

	return false
}
