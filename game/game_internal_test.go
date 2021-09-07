package game

import (
	gomath "math"
	"math/rand"
	"testing"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/stretchr/testify/assert"
)

type fakeRand struct {
	seq []int64
	top int
}

func (f *fakeRand) Int63() int64 {
	if f.top >= len(f.seq) {
		return rand.Int63()
	}
	result := f.seq[f.top]
	f.top++
	return result
}

func (*fakeRand) Seed(seed int64) {
}

func findSeed(min, max float64) (seed int64, num float64) {
	for {
		seed = rand.Int63()
		num = rand.New(&fakeRand{seq: []int64{seed}}).Float64()
		if num >= min && num <= max {
			return
		}
	}
}

func TestCreateAllBuildings(t *testing.T) {
	t.Parallel()

	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		Level: &Level{},
	}

	info := &Info{
		MainGameInfo: mainGameInfo,
		Space:        s,
	}

	// Building 5 is non standard shape
	testCases := []struct {
		name string
		seq  []int64
	}{
		{name: "building 0"},
		{name: "building 1"},
		{name: "building 2"},
		{name: "building 3"},
		{name: "building 4"},
		{name: "building 5"},
	}

	for i, testcase := range testCases {
		step := 1.0 / float64(len(testCases))
		max := step * float64(i+1)
		seed, _ := findSeed(max-step, max)

		info.Rand = rand.New(&fakeRand{seq: []int64{seed}})

		LevelBlock := createRandomLevelBlock(info.Rand, ecs.NewBasic())
		expected := LevelBlock.TileMap.TileWidth * LevelBlock.TileMap.TileXNum
		assert.Equalf(t, expected, int(LevelBlock.Size.X), "invalid tileWidth for %s", testcase.name)

		for _, tile := range LevelBlock.TileMap.Map {
			assert.GreaterOrEqualf(t, tile, int16(0), "invalid greater or equal %s", testcase.name)
		}
	}
}

func TestCityLevelGenerate(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		Level: &Level{},
	}
	info := &Info{
		MainGameInfo: mainGameInfo,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:        s,
	}

	var velocityable *Velocityable
	w.AddSystemInterface(CreateVelocitySystem(s), velocityable, nil)

	generateCityBuildings(info, w)

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
	info := &Info{
		MainGameInfo: mainGameInfo,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:        s,
	}

	gameRuleSystem := CreateGameRuleSystem(info)
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
		VelocityComponent:  &components.VelocityComponent{},
		TileImageComponent: &components.TileImageComponent{},
		CollisionComponent: &components.CollisionComponent{},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		IdentityComponent: &components.IdentityComponent{},
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
	info := &Info{
		MainGameInfo: mainGameInfo,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:        s,
	}

	playerSystem := CreatePlayerSystem(mainGameInfo)
	var playerable *Playerable
	w.AddSystemInterface(playerSystem, playerable, nil)

	var gameRuleable *GameRuleable
	w.AddSystemInterface(CreateGameRuleSystem(info), gameRuleable, nil)
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

func TestDestoryOnAnimeableGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		Gravity: 10,
		Level: &Level{
			// Disable building spawn
			StartX: 50000,
			Width:  500,
		},
	}
	info := &Info{
		MainGameInfo: mainGameInfo,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:        s,
	}

	gameRuleSystem := CreateGameRuleSystem(info)
	var gameRuleable *GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)
	var animeable *Animeable
	w.AddSystemInterface(CreateAnimeSystem(), animeable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.AnimeComponent
		*components.DestoryOnAnimeComponent
		*components.TileImageComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: 1, Y: 1},
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration:  1 * time.Second,
			FrameRemaining: 1 * time.Second,
		},
		DestoryOnAnimeComponent: &components.DestoryOnAnimeComponent{
			CyclesTilDeath: 1,
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: components.CreateTileMap(1, 1, ebiten.NewImage(1, 1), 1),
		},
	}
	w.AddEntity(ent)
	_, ok := gameRuleSystem.ents[ent.ID()]
	assert.True(t, ok, "entity should be removed once anime completed")

	for ent.Cycles <= 0 {
		w.Update(1)
	}
	w.Update(1)

	_, ok = gameRuleSystem.ents[ent.ID()]
	assert.False(t, ok, "entity should be removed once anime completed")
}

func TestEnemyBiscuitGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		Gravity: 10,
		Level: &Level{
			// Disable building spawn
			StartX: 50000,
			Width:  500,
		},
	}
	info := &Info{
		MainGameInfo: mainGameInfo,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:        s,
	}

	gameRuleSystem := CreateGameRuleSystem(info)
	var gameRuleable *GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)
	var animeable *Animeable
	w.AddSystemInterface(CreateAnimeSystem(), animeable, nil)
	var resolvable *Resolvable
	w.AddSystemInterface(CreateResolvSystem(mainGameInfo, s), resolvable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.BiscuitEnemyComponent
		*components.CollisionComponent
		*components.IdentityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: 1, Y: 1},
		},
		BiscuitEnemyComponent: &components.BiscuitEnemyComponent{},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{entity.TagEnemy},
		},
	}
	w.AddEntity(ent)

	w.Update(0.1)
}

func TestUfoBiscuitGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &MainGameInfo{
		Gravity: 10,
		Level: &Level{
			// Disable building spawn
			StartX: 50000,
			Width:  500,
		},
	}
	info := &Info{
		MainGameInfo: mainGameInfo,
		Rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:        s,
	}

	gameRuleSystem := CreateGameRuleSystem(info)
	var gameRuleable *GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)
	var animeable *Animeable
	w.AddSystemInterface(CreateAnimeSystem(), animeable, nil)
	var resolvable *Resolvable
	w.AddSystemInterface(CreateResolvSystem(mainGameInfo, s), resolvable, nil)

	ufo := entity.CreateUfoBiscuitEnemy()
	w.AddEntity(ufo)

	w.Update(0.1)
	_, ok := gameRuleSystem.ents[ufo.ID()]
	assert.True(t, ok, "ufo should exist still")

	ufo.Collisions = append(ufo.Collisions, &components.CollisionEvent{
		Tags: []string{entity.TagBullet},
	})
	w.Update(0.1)
	_, ok = gameRuleSystem.ents[ufo.ID()]
	assert.False(t, ok, "ufo should no longer exist")
}
