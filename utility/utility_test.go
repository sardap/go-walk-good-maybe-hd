package utility_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
	"github.com/stretchr/testify/assert"
)

func TestDeltaToDuration(t *testing.T) {
	t.Parallel()

	assert.Equal(t, time.Duration(500)*time.Millisecond, utility.DeltaToDuration(0.5))
}

func TestAbsInt(t *testing.T) {
	t.Parallel()

	assert.Equal(t, int(1), utility.AbsInt(-1))
	assert.Equal(t, int(1), utility.AbsInt(1))
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

	result = utility.WrapInt(-257, -256, 256)
	assert.Equal(t, result, int(255), "complete double wrap")
}

func TestWrapFloat64(t *testing.T) {
	t.Parallel()

	result := utility.WrapFloat64(1, 0, 10)
	assert.Equal(t, result, float64(1), "should have not wrapped")

	result = utility.WrapFloat64(-1, 0, 10)
	assert.Equal(t, result, float64(9), "should have wrapped min")

	result = utility.WrapFloat64(11, 0, 10)
	assert.Equal(t, result, float64(1), "should have wrapped max")

	result = utility.WrapFloat64(20, 0, 10)
	assert.Equal(t, result, float64(10), "complete double wrap")

	result = utility.WrapFloat64(-257, -256, 256)
	assert.Equal(t, result, float64(255), "complete double wrap")
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

func TestWrapVec2(t *testing.T) {
	t.Parallel()

	min := math.Vector2{
		X: 0, Y: -5,
	}
	max := math.Vector2{
		X: 20, Y: 25,
	}

	val := utility.WrapVec2(
		math.Vector2{
			X: 10, Y: 5,
		}, min, max,
	)
	assert.Equal(t, math.Vector2{X: 10, Y: 5}, val, "no change")
}

func TestRandRange(t *testing.T) {
	t.Parallel()

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100000; i++ {
		val := utility.RandRange(rand, 0, 100)
		assert.GreaterOrEqual(t, val, 0)
		assert.Less(t, val, 100)
		if t.Failed() {
			t.FailNow()
		}
	}
}

func TestRandRangeFloat64(t *testing.T) {
	t.Parallel()

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100000; i++ {
		val := utility.RandRangeFloat64(rand, 0, 100)
		assert.GreaterOrEqual(t, val, float64(0))
		assert.Less(t, val, float64(100))
		if t.Failed() {
			t.FailNow()
		}
	}
}

func TestContainsString(t *testing.T) {
	t.Parallel()

	ary := []string{"foo", "bar"}
	assert.True(t, utility.ContainsString(ary, "foo"))
	assert.True(t, utility.ContainsString(ary, "greg", "bar"))
	assert.False(t, utility.ContainsString(ary, "greg", "got"))
}

func TestContainsInt(t *testing.T) {
	t.Parallel()

	ary := []int{1, 2}
	assert.True(t, utility.ContainsInt(ary, 1))
	assert.True(t, utility.ContainsInt(ary, 1, 2))
	assert.False(t, utility.ContainsInt(ary, 4, 3))
}
