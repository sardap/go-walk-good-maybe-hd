package game_test

import (
	"bytes"
	"errors"
	"image"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/game"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/stretchr/testify/assert"

	"image/color"
	_ "image/png"
)

const (
	// same as ImageWhaleAirTileSet
	imgRaw = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x000\x00\x00\x00\x10\b\x06\x00\x00\x00P\xae\xfc\xb1\x00\x00\x01\u007fIDATx\x9cb\xf9\xff\xff?\xc3P\x06, \xc2\xc4\xc1\xed\xff\x99\x03\xbb\x18\x91%@b\xb84aS;P\xfa\x19\x8d\xed]\xff\v\x8b\x880\xbc}\xf3\x06.\t\xd2,.\xcc\xcf\xf0\x87\x91\x15,\xc6\xce\xce\xce\xf0\xf3\xe7O8\x8d\xaev \xf5\xc3=\x80,\t\xd3\xcc\xf2\xff7\x033\a\x0f\x8afd\x003h \xf53!\xfb\x10\xa4\x01d\x18\f\xbc|\xfb\x11\xae\t\xd9\xe7\xbb\xd6T\xa2\x18\x84K?\b\xfc\xfd\xf1\x05%\x04\xa9\xad\x1f%\x06\u07bf\u007f\xcf\xc0\xc7\xc9\n\xf7=\xc8\x03\xc8\x06\x82<\x80\xceG\x0eAj\xeb\xaf-.c\x10\x97\xe2E\x04\xe8\xb3\xcf\f\xf5}\xfd(\xfa\xc1\x99\x18\xa4\x18\x14Ђ\x82\x82\xf0\x90\x00E\x1d\x03\xc3G\x86ƢB\xb0&\x10\x9dSÙ\r\xe0\xd3O\f\xc0\xa5\x1fd/\xb6\xa4\x83\f\xe01\x00\xd2\x04\xf29r(\x81\x1c\r\x02\xa0P\x00\xf9\x1e\x04\xd0C\x00\x16\x82\xb4\xd0O\b\xc0\xf3\x00,\xbaA\x02\xa0\x9c\r3\x18\xd9r\xe4\xa8\x04\xc9\xc3\xd4\xe2\xd3\x0fr,L\x1f\x88\x869\x9eX\xfd\xf8\x1c\x0eS\xcb\b\xaaȰ\x95\xc3\xc5RR\xff\x0f\xaa\xe90\xac^0\t.\x16\x9a\x90\xc7`\u007f\xeb\nC\xef\xb3g\x04\xcbq\x98~t@u\xfd \x0f`\xc3SUT\xfe\x93\">P\xfa\x99\xf0\xc6\xd5\x10\x00\xa3\x1e\x18h\x80\xd3\x03\xe6\xdbtI\x12\x1f(\xfd\x8cؚӏ?G\xc0\x05eyW0\x12\x12\x1fH\xfd\x80\x00\x00\x00\xff\xff\x8e\x18S\x12\xccё\xca\x00\x00\x00\x00IEND\xaeB`\x82"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAnimeSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	// Setup
	animeSystem := game.CreateAnimeSystem()

	var animeable *game.Animeable
	w.AddSystemInterface(animeSystem, animeable, nil)

	img, _, err := image.Decode(bytes.NewBufferString(imgRaw))
	assert.NoError(t, err)
	eImg := ebiten.NewImageFromImage(img)

	ent := &struct {
		ecs.BasicEntity
		*components.AnimeComponent
		*components.TileImageComponent
	}{
		BasicEntity: ecs.NewBasic(),
		AnimeComponent: &components.AnimeComponent{
			FrameDuration:  50 * time.Millisecond,
			FrameRemaining: 50 * time.Millisecond,
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: components.CreateTileMap(1, 1, eImg, 16),
		},
	}
	ent.TileMap.SetTile(0, 0, 0)

	w.AddEntity(ent)

	// Asserts
	w.Update(0)
	assert.Zero(t, ent.Cycles, "no cycles with no time passing")

	// half anime cycle complete
	w.Update(float32(25*time.Millisecond) / float32(time.Second))
	assert.Zero(t, ent.Cycles, "no cycles with only 25mil passing")

	for i := 0; i < 3; i++ {
		// Next frame
		w.Update(float32(51*time.Millisecond) / float32(time.Second) * float32(i))
		assert.Equal(t, int16(i), ent.TileMap.Get(0, 0), "next frame after 50mil")
	}
	w.Update(float32(51*time.Millisecond) / float32(time.Second))
	assert.Equal(t, int(1), ent.Cycles, "complete cycle should be complete")
	assert.Zero(t, ent.TileMap.Get(0, 0), "frame should wrap")

	w.RemoveEntity(ent.BasicEntity)
	assert.NotZero(t, animeSystem.Priority())
}

func TestSoundSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	soundSystem := game.CreateSoundSystem()
	var soundable *game.Soundable
	w.AddSystemInterface(soundSystem, soundable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.SoundComponent
	}{
		BasicEntity: ecs.NewBasic(),
		SoundComponent: &components.SoundComponent{
			Sound: components.LoadSound(assets.SoundByJumpTwo),
		},
	}
	w.AddEntity(ent)

	w.Update(0)
	assert.Nil(t, ent.Player, "Sound com is not active so it should be nil")

	ent.Active = true
	w.Update(0)
	assert.NotNil(t, ent.Player, "player should be not nil since it's active")
	assert.True(t, ent.Player.IsPlaying())

	ent.Active = false
	w.Update(0)
	assert.True(t, !ent.Player.IsPlaying(), "player should stop when not active")

	ent.Sound = components.LoadSound(assets.MusicPdCity0)
	ent.Restart = true
	ent.Active = true
	w.Update(0)
	assert.True(t, ent.Player.IsPlaying(), "player should be playing when restarted")

	ent.Player.Pause()
	w.Update(0)
	assert.False(t, ent.Active)

	ent.Sound = components.LoadSound(assets.MusicPdCity0)
	ent.Loop = true
	ent.Restart = true
	ent.Active = true
	w.Update(0)
	assert.True(t, ent.Player.IsPlaying())

	w.RemoveEntity(ent.BasicEntity)
	assert.NotZero(t, soundSystem.Priority())
}

func TestVelocitySystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()

	// Setup
	velocitySystem := game.CreateVelocitySystem(s)
	var velocityable *game.Velocityable
	w.AddSystemInterface(velocitySystem, velocityable, nil)

	var resolvable *game.Resolvable
	w.AddSystemInterface(game.CreateResolvSystem(s), resolvable, nil)

	entA := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
		*components.VelocityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{"test"},
		},
		CollisionComponent: &components.CollisionComponent{
			Active:         true,
			CollisionShape: nil,
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
	}

	entB := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
		*components.VelocityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Postion: math.Vector2{
				X: 20,
				Y: 20,
			},
			Size: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{"ground"},
		},
		CollisionComponent: &components.CollisionComponent{
			Active:         true,
			CollisionShape: nil,
		},
		VelocityComponent: &components.VelocityComponent{},
	}

	w.AddEntity(entA)
	w.AddEntity(entB)

	// 1.5 seconds passed
	entA.Vel.Y = 3
	w.Update(1.5)
	assert.Zero(t, entA.Vel.X, "vel X  should be set to 0")
	assert.Zero(t, entA.Vel.Y, "vel Y should be set to 0")
	assert.Equal(t, float64(3*1.5), entA.Postion.Y)

	entA.Vel.Y = 10
	w.Update(1)
	assert.Equal(t, float64(9.5), entA.Postion.Y, "bounds y to stop collsion")
	assert.True(t, entA.Collisions.CollidingWith("ground"))
	assert.True(t, entB.Collisions.CollidingWith("test"))

	entA.Vel.X = 10
	w.Update(1)
	assert.Equal(t, float64(25), entA.Postion.X, "bounds X to stop collsion")
	assert.True(t, entA.Collisions.CollidingWith("ground"))
	assert.True(t, entB.Collisions.CollidingWith("test"))

	colShape := resolv.NewRectangle(0, 0, 10, 10)
	entC := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
		*components.VelocityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Postion: math.Vector2{
				X: 20,
				Y: 20,
			},
			Size: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{"ground"},
		},
		CollisionComponent: &components.CollisionComponent{
			Active:         true,
			CollisionShape: colShape,
		},
		VelocityComponent: &components.VelocityComponent{},
	}
	w.AddEntity(entC)
	assert.Equal(t, colShape, entC.CollisionShape, "if collission shape is provided it should be unchanged")
	assert.True(t, s.Contains(colShape), "it should be added to the space")
}

func TestResolvSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()

	// Setup
	resolvSystem := game.CreateResolvSystem(s)

	var resolvable *game.Resolvable
	w.AddSystemInterface(resolvSystem, resolvable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{"test"},
		},
		CollisionComponent: &components.CollisionComponent{
			Active:         true,
			CollisionShape: nil,
		},
	}

	// asserts
	assert.False(t, s.Contains(ent.CollisionShape), "space should be empty")
	w.AddEntity(ent)
	assert.NotNil(t, ent.CollisionShape, "collsion shape should be init")
	assert.True(t, s.Contains(ent.CollisionShape), "space should contain shape")
	assert.Equal(t, "test", ent.CollisionShape.GetTags()[0], "tags should be copied to shape")

	w.RemoveEntity(ent.BasicEntity)
	assert.False(t, s.Contains(ent.CollisionShape), "shape should be removed from space")

	w.AddEntity(ent)

	// don't need to test debug overlay so just do it here
	renderQueue := make(game.RenderCmds, 0)
	resolvSystem.OverlayEnabled = true
	resolvSystem.Render(&renderQueue)
	assert.Equal(t, 1, len(renderQueue))

	renderQueue = make(game.RenderCmds, 0)
	resolvSystem.OverlayEnabled = false
	resolvSystem.Render(&renderQueue)
	assert.Equal(t, 0, len(renderQueue))
}

func TestImageRenderSystem(t *testing.T) {

	w := &ecs.World{}

	imageRenderSystem := game.CreateImageRenderSystem()

	var imageRenderable *game.ImageRenderable
	w.AddSystemInterface(imageRenderSystem, imageRenderable, nil)

	renderQueue := make(game.RenderCmds, 0)

	ents := []*struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.ImageComponent
	}{}

	basics := ecs.NewBasics(10)
	for i := 0; i < 10; i++ {
		img := ebiten.NewImage(10, 10)
		clr := color.RGBA{byte(i * 10), 255, 255, 255}
		img.Fill(clr)

		ent := &struct {
			ecs.BasicEntity
			*components.TransformComponent
			*components.ImageComponent
		}{
			BasicEntity:        basics[i],
			TransformComponent: &components.TransformComponent{},
			ImageComponent: &components.ImageComponent{
				Active: true,
				Image:  img,
				Layer:  components.ImageLayer(i),
			},
		}

		ents = append(ents, ent)
	}

	shuffled := make([]int, len(ents))
	for i := range ents {
		shuffled[i] = i
	}
	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	for _, i := range shuffled {
		w.AddEntity(ents[i])
	}

	imageRenderSystem.Render(&renderQueue)
	renderQueue.Sort()

	assert.Equal(t, len(renderQueue), len(ents), "ents missing from the queue")

	screen := ebiten.NewImage(300, 300)
	last := renderQueue[len(renderQueue)-1]
	for i := len(renderQueue) - 2; i > 0; i-- {
		currentEnt := ents[i]
		current := renderQueue[i]
		assert.Equal(t, int(currentEnt.Layer), current.GetLayer())
		assert.Greater(t, last.GetLayer(), current.GetLayer())

		current.Draw(screen)
		expectedColor := color.RGBA{byte(i * 10), 255, 255, 255}
		assert.Equal(t, expectedColor, screen.At(0, 0), "incorrect colour at 0,0")

		last = current
	}
}

func TestTileImageRenderSystem(t *testing.T) {

	w := &ecs.World{}

	tileImageRenderSystem := game.CreateTileImageRenderSystem()

	var tileImageRenderable *game.TileImageRenderable
	w.AddSystemInterface(tileImageRenderSystem, tileImageRenderable, nil)

	ents := []*struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.TileImageComponent
	}{}

	basics := ecs.NewBasics(10)
	for i := 0; i < 10; i++ {
		img := ebiten.NewImage(3, 1)
		value := byte((i + 1) * 10)
		img.Set(0, 0, color.RGBA{
			R: value,
			A: 255,
		})
		img.Set(1, 0, color.RGBA{
			G: value + 1,
			A: 255,
		})
		img.Set(2, 0, color.RGBA{
			B: value + 2,
			A: 255,
		})

		tileMap := components.CreateTileMap(2, 2, img, 1)

		ent := &struct {
			ecs.BasicEntity
			*components.TransformComponent
			*components.TileImageComponent
		}{
			BasicEntity:        basics[i],
			TransformComponent: &components.TransformComponent{},
			TileImageComponent: &components.TileImageComponent{
				Active:  true,
				TileMap: tileMap,
				Layer:   components.ImageLayer(i),
			},
		}
		tileMap.SetTile(0, 0, 0)
		tileMap.SetTile(1, 0, 1)
		tileMap.SetTile(0, 1, 2)

		ents = append(ents, ent)
		w.AddEntity(ents[i])
	}

	renderQueue := make(game.RenderCmds, 0)
	tileImageRenderSystem.Render(&renderQueue)

	renderQueue.Sort()
	assert.Equal(t, len(renderQueue), len(ents), "missing ents from render queue")

	screen := ebiten.NewImage(2, 2)
	for i := len(renderQueue) - 1; i > 0; i-- {
		screen.Fill(color.Black)
		current := renderQueue[i]

		current.Draw(screen)
		expectedValue := byte((i + 1) * 10)

		assert.Equal(t, color.RGBA{R: expectedValue, A: 255}, screen.At(0, 0), "incorrect colour at 0,0")
		assert.Equal(t, color.RGBA{G: expectedValue + 1, A: 255}, screen.At(1, 0), "incorrect colour at 1,0")
		assert.Equal(t, color.RGBA{B: expectedValue + 2, A: 255}, screen.At(0, 1), "incorrect colour at 0,1")
		assert.Equal(t, color.RGBA{A: 255}, screen.At(1, 1), "incorrect colour at 0,1")
	}

	for _, ent := range ents {
		w.RemoveEntity(ent.BasicEntity)
	}

	ent := ents[0]
	w.AddEntity(ent)

	ent.TileMap.Options.InvertX = true
	ent.TileMap.Options.InvertY = true

	renderQueue = make(game.RenderCmds, 0)
	tileImageRenderSystem.Render(&renderQueue)
	renderQueue.Sort()

	screen.Fill(color.Black)
	renderQueue[0].Draw(screen)

	// I really don't think this works right
	expectedValue := byte(10)
	assert.Equal(t, color.RGBA{R: expectedValue, A: 255}, screen.At(0, 0), "incorrect colour at 0,0")
	assert.Equal(t, color.RGBA{G: expectedValue + 1, A: 255}, screen.At(1, 0), "incorrect colour at 1,0")
	assert.Equal(t, color.RGBA{B: expectedValue + 2, A: 255}, screen.At(0, 1), "incorrect colour at 0,1")
	assert.Equal(t, color.RGBA{A: 255}, screen.At(1, 1), "incorrect colour at 0,1")

	assert.NotZero(t, tileImageRenderSystem.Priority())
	tileImageRenderSystem.Update(0)
}

func TestTextRenderSystem(t *testing.T) {

	w := &ecs.World{}

	textRenderSystem := game.CreateTextRenderSystem()
	var textRenderable *game.TextRenderable
	w.AddSystemInterface(textRenderSystem, textRenderable, nil)

	renderQueue := make(game.RenderCmds, 0)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.TextComponent
	}{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		TextComponent: &components.TextComponent{
			Text:  "why i'm I chasing test coverage",
			Layer: 0,
		},
	}
	w.AddEntity(ent)

	textRenderSystem.Render(&renderQueue)
	renderQueue.Sort()

	assert.Equal(t, 1, len(renderQueue))

	screen := ebiten.NewImage(300, 300)
	screen.Fill(color.White)
	renderQueue[0].Draw(screen)

	valid := false
	for y := screen.Bounds().Min.Y; y < screen.Bounds().Dy(); y++ {
		for x := screen.Bounds().Min.X; x < screen.Bounds().Dx(); x++ {
			if r, _, _, _ := screen.At(x, y).RGBA(); r != 255 {
				valid = true
			}
		}
	}
	assert.True(t, valid, "at least one pixel should not be white")

	w.RemoveEntity(ent.BasicEntity)

	assert.NotZero(t, textRenderSystem.Priority())
}

func TestWrapInGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &game.MainGameInfo{
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameInfo, s)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	wrapEnt := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.WrapComponent
	}{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		WrapComponent: &components.WrapComponent{
			Min: math.Vector2{X: 0, Y: 0},
			Max: math.Vector2{X: 50, Y: 0},
		},
	}
	w.AddEntity(wrapEnt)

	wrapTestCases := []struct {
		Postion  math.Vector2
		Expected math.Vector2
	}{
		{
			Postion:  math.Vector2{},
			Expected: math.Vector2{},
		},
		{
			Postion:  math.Vector2{X: 51},
			Expected: math.Vector2{X: 1},
		},
	}

	for _, testCase := range wrapTestCases {
		wrapEnt.Postion = testCase.Postion
		w.Update(1)
		assert.Equal(t, testCase.Expected, wrapEnt.Postion, "No movement no change")
	}

	w.RemoveEntity(wrapEnt.BasicEntity)
}

func TestScrollInGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &game.MainGameInfo{
		ScrollingSpeed: math.Vector2{X: 50, Y: 0},
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameInfo, s)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	scrollEnt := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.VelocityComponent
		*components.ScrollableComponent
	}{
		BasicEntity:         ecs.NewBasic(),
		TransformComponent:  &components.TransformComponent{},
		VelocityComponent:   &components.VelocityComponent{},
		ScrollableComponent: &components.ScrollableComponent{},
	}
	w.AddEntity(scrollEnt)

	scrollTestCases := []struct {
		Vel         math.Vector2
		ExpectedVel math.Vector2
	}{
		{
			Vel:         math.Vector2{},
			ExpectedVel: math.Vector2{X: 50},
		},
	}

	for _, testCase := range scrollTestCases {
		scrollEnt.Vel = testCase.Vel
		w.Update(1)
		assert.Equal(t, testCase.ExpectedVel, scrollEnt.Vel, "No vel no change")
	}

	w.RemoveEntity(scrollEnt.BasicEntity)
}

func TestGravityInGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &game.MainGameInfo{
		ScrollingSpeed: math.Vector2{X: 50, Y: 0},
		Gravity:        10,
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameInfo, s)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	var velocityable *game.Velocityable
	w.AddSystemInterface(game.CreateVelocitySystem(s), velocityable, nil)

	var resolvable *game.Resolvable
	w.AddSystemInterface(game.CreateResolvSystem(s), resolvable, nil)

	fallingEnt := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.CollisionComponent
		*components.GravityComponent
		*components.IdentityComponent
		*components.VelocityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: 1, Y: 1},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		GravityComponent:  &components.GravityComponent{},
		VelocityComponent: &components.VelocityComponent{},
		IdentityComponent: &components.IdentityComponent{},
	}
	w.AddEntity(fallingEnt)

	ground := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Postion: math.Vector2{
				X: 0,
				Y: 100,
			},
			Size: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{"ground"},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
	}
	w.AddEntity(ground)

	for i := 0; i < 9; i++ {
		w.Update(1)
		assert.Equal(t, float64((i+1)*10), fallingEnt.Postion.Y, "no postion change")
	}

	w.Update(1)
	assert.True(t, fallingEnt.Collisions.CollidingWith(entity.TagGround))

	w.RemoveEntity(fallingEnt.BasicEntity)
}

func TestBulletInGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameInfo := &game.MainGameInfo{
		Gravity: 10,
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
			Width:  500,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameInfo, s)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	var velocityable *game.Velocityable
	w.AddSystemInterface(game.CreateVelocitySystem(s), velocityable, nil)

	var resolvable *game.Resolvable
	w.AddSystemInterface(game.CreateResolvSystem(s), resolvable, nil)

	bullet := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.BulletComponent
		*components.CollisionComponent
		*components.IdentityComponent
		*components.VelocityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: 1, Y: 1},
		},
		BulletComponent: &components.BulletComponent{
			Speed: math.Vector2{X: 10},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		VelocityComponent: &components.VelocityComponent{},
		IdentityComponent: &components.IdentityComponent{},
	}
	w.AddEntity(bullet)
	for bullet.Postion.X <= mainGameInfo.Level.Width {
		assert.True(t, s.Contains(bullet.CollisionShape))
		w.Update(1)
	}
	w.Update(1)
	assert.False(t, s.Contains(bullet.CollisionShape), "should be removed off screen")

	bullet.Postion = math.Vector2{}
	bullet.CollisionShape = nil
	bullet.Speed.X = -10
	w.AddEntity(bullet)
	for bullet.Postion.X >= 0 {
		assert.True(t, s.Contains(bullet.CollisionShape))
		w.Update(1)
	}
	w.Update(1)
	assert.False(t, s.Contains(bullet.CollisionShape), "should be removed off screen")
}

type testGame struct {
	m    *testing.M
	code int
}

var (
	errRegularTermination = errors.New("regular termination")
)

func (g *testGame) Update() error {
	g.code = g.m.Run()
	return errRegularTermination
}

func (*testGame) Draw(screen *ebiten.Image) {
}

func (*testGame) Layout(int, int) (int, int) {
	return 300, 300
}

func TestMain(m *testing.M) {
	g := &testGame{
		m: m,
	}
	if err := ebiten.RunGame(g); err != nil && err != errRegularTermination {
		panic(err)
	}
	os.Exit(g.code)
}
