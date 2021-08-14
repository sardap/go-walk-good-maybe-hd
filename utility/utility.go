package utility

import (
	"time"
)

func DeltaToDuration(dt float32) time.Duration {
	var i int = int(dt)
	return time.Duration(i) * time.Millisecond
}

func WrapInt(x, min, max int) int {
	if x >= max {
		return x + min - max
	} else if x < min {
		return x + max - min
	}
	return x
}
