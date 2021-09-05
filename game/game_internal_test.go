package game

import (
	"fmt"
	gomath "math"
	"math/rand"
	"testing"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
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

	garbage := rand.Int63()
	fmt.Printf("%d:%f\n", garbage, rand.New(&fakeRand{seq: []int64{garbage}}).Float64())

	testCases := []struct {
		name string
		seq  []int64
	}{
		// 0.058
		{name: "building 0", seq: []int64{543898249864404905}},
		// 0.276
		{name: "building 1", seq: []int64{2552056708243655789}},
		// 0.480
		{name: "building 2", seq: []int64{4430758520341445551}},
		// 0.608
		{name: "building 3", seq: []int64{5610831761584960326}},
		// 0.940
		{name: "building 4", seq: []int64{8678356114173921525}},
	}

	for _, testcase := range testCases {
		info.Rand = rand.New(&fakeRand{seq: testcase.seq})

		LevelBlock := createRandomLevelBlock(info.Rand, ecs.NewBasic())
		assert.Equal(t, LevelBlock.TileMap.TileWidth*LevelBlock.TileMap.TileXNum, int(LevelBlock.Size.X))
		for _, tile := range LevelBlock.TileMap.Map {
			assert.GreaterOrEqual(t, tile, int16(0))
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
