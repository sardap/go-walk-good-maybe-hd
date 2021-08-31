package utility_test

import (
	"testing"
	"time"

	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
	"github.com/stretchr/testify/assert"
)

func TestDeltaToDuration(t *testing.T) {
	assert.Equal(t, time.Duration(500)*time.Millisecond, utility.DeltaToDuration(0.5))
}

func TestWrapInt(t *testing.T) {
	t.Parallel()

	result := utility.WrapInt(1, 0, 10)
	assert.Equal(t, result, int(1), "should have not wrapped")

	result = utility.WrapInt(-1, 0, 10)
	assert.Equal(t, result, int(9), "should have wrapped min")

	result = utility.WrapInt(11, 0, 10)
	assert.Equal(t, result, int(1), "should have wrapped max")

	result = utility.WrapInt(20, 0, 10)
	assert.Equal(t, result, int(10), "complete double wrap")
}

func TestClampVec2(t *testing.T) {
	t.Parallel()

	min := math.Vector2{
		X: 0, Y: -5,
	}
	max := math.Vector2{
		X: 20, Y: 25,
	}

	val := utility.ClampVec2(
		math.Vector2{
			X: 10, Y: 5,
		}, min, max,
	)
	assert.Equal(t, math.Vector2{X: 10, Y: 5}, val, "no change")

	val = utility.ClampVec2(
		math.Vector2{
			X: 25, Y: 30,
		}, min, max,
	)
	assert.Equal(t, math.Vector2{X: 20, Y: 25}, val, "max clamp both values")

	val = utility.ClampVec2(
		math.Vector2{
			X: -5, Y: -10,
		}, min, max,
	)
	assert.Equal(t, math.Vector2{X: 0, Y: -5}, val, "min clamp both values")
}

func TestRandRange(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100000; i++ {
		val := utility.RandRange(0, 100)
		assert.GreaterOrEqual(t, val, 0)
		assert.Less(t, val, 100)
		if t.Failed() {
			t.FailNow()
		}
	}
}

func TestRandRangeFloat64(t *testing.T) {
	t.Parallel()

	for i := 0; i < 100000; i++ {
		val := utility.RandRangeFloat64(0, 100)
		assert.GreaterOrEqual(t, val, float64(0))
		assert.Less(t, val, float64(100))
		if t.Failed() {
			t.FailNow()
		}
	}
}
