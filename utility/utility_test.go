package utility_test

import (
	"testing"

	"github.com/sardap/walk-good-maybe-hd/utility"
	"github.com/stretchr/testify/assert"
)

func TestWrapInt(t *testing.T) {
	result := utility.WrapInt(1, 0, 10)
	assert.Equal(t, result, int(1), "should have not wrapped")

	result = utility.WrapInt(-1, 0, 10)
	assert.Equal(t, result, int(9), "should have wrapped min")

	result = utility.WrapInt(11, 0, 10)
	assert.Equal(t, result, int(1), "should have wrapped max")

	result = utility.WrapInt(20, 0, 10)
	assert.Equal(t, result, int(10), "complete double wrap")
}
