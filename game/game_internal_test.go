package game

import (
	gomath "math"
	"testing"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/stretchr/testify/assert"
)

func TestCityLevelGenerate(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		Level: &Level{},
	}

	var velocityable *Velocityable
	w.AddSystemInterface(CreateVelocitySystem(s), velocityable, nil)

	generateCityBuildings(mainGameInfo, w)

	ground := s.FilterByTags(entity.TagGround)
	for i := 0; i < ground.Length()-1; i++ {
		left := ground.Get(i).(*resolv.Rectangle)
		right := ground.Get(i + 1).(*resolv.Rectangle)
		dist := gomath.Abs((left.X + left.W) - (right.X + right.W))
		assert.Less(t, dist, float64(150), "distance between buildings must be jumpable")
	}
}

func TestBuildingsRemoveGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		ScrollingSpeed: math.Vector2{X: -1, Y: 0},
		Level: &Level{
			// Disable building spawn Yes we want this
			StartX: 50000,
		},
	}

	gameRuleSystem := CreateGameRuleSystem(mainGameInfo, s)
	var gameRuleable *GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	// Buildings need to move
	var resolveable *Resolvable
	w.AddSystemInterface(CreateResolvSystem(mainGameInfo, s), resolveable, nil)
	var velocityable *Velocityable
	w.AddSystemInterface(CreateVelocitySystem(s), velocityable, nil)

	block := &LevelBlock{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Postion: math.Vector2{
				X: 10,
			},
		},
		VelocityComponent:   &components.VelocityComponent{},
		TileImageComponent:  &components.TileImageComponent{},
		CollisionComponent:  &components.CollisionComponent{},
		ScrollableComponent: &components.ScrollableComponent{},
		IdentityComponent:   &components.IdentityComponent{},
	}

	w.AddEntity(block)

	scrollTestCases := []struct {
		loops int
	}{
		{
			loops: 12,
		},
	}

	for _, testCase := range scrollTestCases {
		for i := 0; i < testCase.loops-1; i++ {
			w.Update(1)
			_, ok := gameRuleSystem.ents[block.ID()]
			assert.Truef(t, ok, "Should still have building on loop %d", i)
		}
		w.Update(1)
		_, ok := gameRuleSystem.ents[block.ID()]
		assert.False(t, ok, "building should be removed")
	}

	w.RemoveEntity(block.BasicEntity)
}

func TestPlayerSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		ScrollingSpeed: math.Vector2{X: -1, Y: 0},
		Level: &Level{
			Width: 500,
		},
	}

	playerSystem := CreatePlayerSystem(mainGameInfo)
	var playerable *Playerable
	w.AddSystemInterface(playerSystem, playerable, nil)

	var gameRuleable *GameRuleable
	w.AddSystemInterface(CreateGameRuleSystem(mainGameInfo, s), gameRuleable, nil)
	var resolveable *Resolvable
	w.AddSystemInterface(CreateResolvSystem(mainGameInfo, s), resolveable, nil)
	var velocityable *Velocityable
	w.AddSystemInterface(CreateVelocitySystem(s), velocityable, nil)

	player := entity.CreatePlayer()
	player.Postion.X = 5
	w.AddEntity(player)

	player.MoveRight = true
	lastPostion := player.Postion
	w.Update(0.1)
	assert.Less(t, lastPostion.X, player.Postion.X, "player should move right")

	player.MoveLeft = true
	lastPostion = player.Postion
	w.Update(0.1)
	assert.Greater(t, lastPostion.X, player.Postion.X, "player should move left")

	w.Update(1)

	// Loop until hit ground
	for !player.Collisions.CollidingWith(entity.TagGround) {
		w.Update(0.1)
	}
	// Change state
	w.Update(0.1)

	// Jump
	player.MoveUp = true
	w.Update(0.1)
	assert.False(t, player.MoveUp, "movement should be reset every update")
	player.Cycles = 1
	w.Update(0.1)
	assert.Equal(t, player.Cycles, 0, "cycles should get changed on anime change")
	lastPostion = player.Postion
	w.Update(0.1)
	assert.Less(t, lastPostion.Y, player.Postion.Y)

	// Player jumping
	for !player.Collisions.CollidingWith(entity.TagGround) {
		w.Update(0.1)
	}

	player.MoveRight = true
	lastPostion = player.Postion
	w.Update(0.9)
	assert.Less(t, lastPostion.X, player.Postion.X, "player should move right")

	player.MoveLeft = true
	lastPostion = player.Postion
	w.Update(0.01)
	assert.Greater(t, lastPostion.X, player.Postion.X, "player should move left")

	// Check bullet is flipped
	player.Shoot = true
	// This is very sensitive since a bullet can be destroyed if the player is too far left
	w.Update(0.001)
	bullet, ok := s.FilterByTags(entity.TagBullet).Get(0).GetData().(*entity.Bullet)
	assert.True(t, ok, "should be bullet attached to data")
	assert.True(t, bullet.Options.InvertX, "bullet should be flipped")
	assert.Less(t, bullet.Speed.X, float64(0), "bullet should be moving left")

	w.RemoveEntity(player.BasicEntity)
}
