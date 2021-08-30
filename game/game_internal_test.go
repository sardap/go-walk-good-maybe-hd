package game

import (
	"math"
	"testing"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/stretchr/testify/assert"
)

func TestCityLevelGenerate(t *testing.T) {
	t.Parallel()

	mainGameInfo = &MainGameInfo{
		level: &Level{},
	}

	w := &ecs.World{}
	s := resolv.NewSpace()

	var velocityable *Velocityable
	w.AddSystemInterface(CreateVelocitySystem(s), velocityable, nil)

	generateCityBuildings(w)

	ground := s.FilterByTags(entity.TagGround)
	for i := 0; i < ground.Length()-1; i++ {
		left := ground.Get(i).(*resolv.Rectangle)
		right := ground.Get(i + 1).(*resolv.Rectangle)
		dist := math.Abs((left.X + left.W) - (right.X + right.W))
		assert.Less(t, dist, float64(100), "distance between buildings must be jumpable")
	}
}
