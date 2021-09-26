package game_test

import (
	"bytes"
	"errors"
	"image"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/game"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"image/color"
	_ "image/png"
)

const (
	// same as ImageWhaleAirTileSet
	imgRaw = "\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x000\x00\x00\x00\x10\b\x06\x00\x00\x00P\xae\xfc\xb1\x00\x00\x01\u007fIDATx\x9cb\xf9\xff\xff?\xc3P\x06, \xc2\xc4\xc1\xed\xff\x99\x03\xbb\x18\x91%@b\xb84aS;P\xfa\x19\x8d\xed]\xff\v\x8b\x880\xbc}\xf3\x06.\t\xd2,.\xcc\xcf\xf0\x87\x91\x15,\xc6\xce\xce\xce\xf0\xf3\xe7O8\x8d\xaev \xf5\xc3=\x80,\t\xd3\xcc\xf2\xff7\x033\a\x0f\x8afd\x003h \xf53!\xfb\x10\xa4\x01d\x18\f\xbc|\xfb\x11\xae\t\xd9\xe7\xbb\xd6T\xa2\x18\x84K?\b\xfc\xfd\xf1\x05%\x04\xa9\xad\x1f%\x06\u07bf\u007f\xcf\xc0\xc7\xc9\n\xf7=\xc8\x03\xc8\x06\x82<\x80\xceG\x0eAj\xeb\xaf-.c\x10\x97\xe2E\x04\xe8\xb3\xcf\f\xf5}\xfd(\xfa\xc1\x99\x18\xa4\x18\x14Ђ\x82\x82\xf0\x90\x00E\x1d\x03\xc3G\x86ƢB\xb0&\x10\x9dSÙ\r\xe0\xd3O\f\xc0\xa5\x1fd/\xb6\xa4\x83\f\xe01\x00\xd2\x04\xf29r(\x81\x1c\r\x02\xa0P\x00\xf9\x1e\x04\xd0C\x00\x16\x82\xb4\xd0O\b\xc0\xf3\x00,\xbaA\x02\xa0\x9c\r3\x18\xd9r\xe4\xa8\x04\xc9\xc3\xd4\xe2\xd3\x0fr,L\x1f\x88\x869\x9eX\xfd\xf8\x1c\x0eS\xcb\b\xaaȰ\x95\xc3\xc5RR\xff\x0f\xaa\xe90\xac^0\t.\x16\x9a\x90\xc7`\u007f\xeb\nC\xef\xb3g\x04\xcbq\x98~t@u\xfd \x0f`\xc3SUT\xfe\x93\">P\xfa\x99\xf0\xc6\xd5\x10\x00\xa3\x1e\x18h\x80\xd3\x03\xe6\xdbtI\x12\x1f(\xfd\x8cؚӏ?G\xc0\x05eyW0\x12\x12\x1fH\xfd\x80\x00\x00\x00\xff\xff\x8e\x18S\x12\xccё\xca\x00\x00\x00\x00IEND\xaeB`\x82"
)

const (
	tagTest = iota
	tagGround
	tagObject
)

func init() {
	rand.Seed(time.Now().UnixNano())
	if audio.CurrentContext() == nil {
		audio.NewContext(48000)
	}
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
	w := &ecs.World{}

	soundSystem := game.CreateSoundSystem(audio.CurrentContext())
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
	debugEnt := entity.CreateDebugInput()

	// Setup
	velocitySystem := game.CreateVelocitySystem(s)
	var velocityable *game.Velocityable
	w.AddSystemInterface(velocitySystem, velocityable, nil)

	var resolvable *game.Resolvable
	w.AddSystemInterface(game.CreateResolvSystem(s, debugEnt), resolvable, nil)

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
			Tags: []int{tagTest},
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
			Tags: []int{tagGround},
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
	assert.Equal(t, float64(14.5), entA.Postion.Y, "bounds y to stop collsion")
	assert.True(t, entA.Collisions.CollidingWith(tagGround))
	assert.True(t, entB.Collisions.CollidingWith(tagTest))

	entA.Vel.X = 10
	w.Update(1)
	assert.Equal(t, float64(25), entA.Postion.X, "bounds X to stop collsion")
	assert.True(t, entA.Collisions.CollidingWith(tagGround))
	assert.True(t, entB.Collisions.CollidingWith(tagTest))

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
			Tags: []int{tagGround},
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

func TestDumbVelocitySystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	// Setup
	dumbVelocitySystem := game.CreateDumbVelocitySystem()
	var dumbVelocityable *game.DumbVelocityable
	var exVelocityable *game.ExDumbVelocityable
	w.AddSystemInterface(dumbVelocitySystem, dumbVelocityable, exVelocityable)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.VelocityComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{
				X: 10,
				Y: 10,
			},
		},
	}
	w.AddEntity(ent)

	// 1.5 seconds passed
	ent.Vel.Y = 3
	w.Update(1.5)
	assert.Zero(t, ent.Vel.X, "vel X  should be set to 0")
	assert.Zero(t, ent.Vel.Y, "vel Y should be set to 0")
	assert.Equal(t, float64(3*1.5), ent.Postion.Y)

	w.RemoveEntity(ent.BasicEntity)

	assert.NotZero(t, dumbVelocitySystem.Priority())
}

func TestResolvSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	debugEnt := entity.CreateDebugInput()

	// Setup
	resolvSystem := game.CreateResolvSystem(s, debugEnt)

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
			Tags: []int{tagTest},
		},
		CollisionComponent: &components.CollisionComponent{
			Active:         true,
			CollisionShape: nil,
		},
	}

	object := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: 3,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{tagObject},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
	}
	w.AddEntity(object)

	ground := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.IdentityComponent
		*components.CollisionComponent
	}{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: 50,
				Y: 10,
			},
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{tagGround},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
	}
	w.AddEntity(ground)

	assert.False(t, s.Contains(ent.CollisionShape), "space should be empty")
	w.AddEntity(ent)
	assert.NotNil(t, ent.CollisionShape, "collsion shape should be init")
	assert.True(t, s.Contains(ent.CollisionShape), "space should contain shape")
	assert.Equal(t, tagTest, ent.CollisionShape.GetTags()[0], "tags should be copied to shape")

	testCases := []struct {
		name               string
		groundPostion      math.Vector2
		entPostion         math.Vector2
		objectPostion      math.Vector2
		entSpeed           math.Vector2
		entCollisionsCount []int
	}{
		{
			name:               "basic",
			groundPostion:      math.Vector2{X: 0, Y: 10},
			entPostion:         math.Vector2{X: 0, Y: 0},
			objectPostion:      math.Vector2{X: 14, Y: 0},
			entSpeed:           math.Vector2{X: 1, Y: 0},
			entCollisionsCount: []int{1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1},
		},
		{
			name:               "fast moving",
			groundPostion:      math.Vector2{X: 0, Y: 10},
			entPostion:         math.Vector2{X: 0, Y: 0},
			objectPostion:      math.Vector2{X: 18, Y: 0},
			entSpeed:           math.Vector2{X: 4, Y: 0},
			entCollisionsCount: []int{1, 1, 2, 2},
		},
	}

	for _, testCase := range testCases {
		ent.Postion = testCase.entPostion
		ground.Postion = testCase.groundPostion
		object.Postion = testCase.objectPostion

		assert.Equal(t, len(ent.Collisions), 0)
		assert.Equal(t, len(ground.Collisions), 0)

		for i, colCount := range testCase.entCollisionsCount {
			w.Update(0.1)
			ent.Postion = ent.Postion.Add(testCase.entSpeed)

			assert.Equalf(t, colCount, len(ent.Collisions), "case: %s loop: %d", testCase.name, i)
		}

		ent.Collisions = nil
		ground.Collisions = nil
	}

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

	for _, ent := range ents {
		w.RemoveEntity(ent.BasicEntity)
	}

	ents[0].Options.InvertX = true
	ents[0].Options.InvertY = true
	img := ebiten.NewImage(2, 1)
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	img.Set(1, 0, color.RGBA{G: 255, A: 255})
	ents[0].Image = img
	w.AddEntity(ents[0])
	renderQueue = nil
	imageRenderSystem.Render(&renderQueue)
	screen.Fill(color.Black)
	renderQueue[0].Draw(screen)
	assert.Equal(t, color.RGBA{G: 255, A: 255}, screen.At(0, 0))
	assert.Equal(t, color.RGBA{R: 255, A: 255}, screen.At(1, 0))

	w.Update(0.1)

	assert.Zero(t, imageRenderSystem.Priority())
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

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	font, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    160,
		DPI:     72 * 2,
		Hinting: font.HintingFull,
	})

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
			Font:  font,
			Color: color.Black,
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

	ent.TextComponent.Text = "updatedText"
	renderQueue = nil
	textRenderSystem.Render(&renderQueue)
	screen.Fill(color.White)
	renderQueue[0].Draw(screen)

	valid = false
	for y := screen.Bounds().Min.Y; y < screen.Bounds().Dy(); y++ {
		for x := screen.Bounds().Min.X; x < screen.Bounds().Dx(); x++ {
			if r, _, _, _ := screen.At(x, y).RGBA(); r != 255 {
				valid = true
			}
		}
	}
	assert.True(t, valid, "at least one pixel should not be white")

	w.RemoveEntity(ent.BasicEntity)

	w.Update(0.1)

	assert.NotZero(t, textRenderSystem.Priority())
}

func TestWrapInGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameScene := &game.MainGameScene{
		Rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:   s,
		World:   w,
		Gravity: 10,
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameScene)
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
	mainGameScene := &game.MainGameScene{
		Rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:          s,
		World:          w,
		Gravity:        10,
		ScrollingSpeed: math.Vector2{X: 50, Y: 0},
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameScene)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	scrollEnt := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.VelocityComponent
		*components.ScrollableComponent
	}{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		VelocityComponent:  &components.VelocityComponent{},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
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
	mainGameScene := &game.MainGameScene{
		Rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:          s,
		World:          w,
		Gravity:        10,
		ScrollingSpeed: math.Vector2{X: 50, Y: 0},
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameScene)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	var velocityable *game.Velocityable
	w.AddSystemInterface(game.CreateVelocitySystem(s), velocityable, nil)

	var resolvable *game.Resolvable
	w.AddSystemInterface(game.CreateResolvSystem(s, mainGameScene.InputEnt), resolvable, nil)

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
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{tagTest},
		},
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
			Tags: []int{tagGround},
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
	assert.True(t, fallingEnt.Collisions.CollidingWith(tagGround))

	w.RemoveEntity(fallingEnt.BasicEntity)
}

func TestBulletInGameRuleSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()
	mainGameScene := &game.MainGameScene{
		Rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Space:   s,
		World:   w,
		Gravity: 10,
		Level: &game.Level{
			// Disable building spawn
			StartX: 50000,
			Width:  500,
		},
	}

	gameRuleSystem := game.CreateGameRuleSystem(mainGameScene)
	var gameRuleable *game.GameRuleable
	w.AddSystemInterface(gameRuleSystem, gameRuleable, nil)

	var velocityable *game.Velocityable
	w.AddSystemInterface(game.CreateVelocitySystem(s), velocityable, nil)

	var resolvable *game.Resolvable
	w.AddSystemInterface(game.CreateResolvSystem(s, mainGameScene.InputEnt), resolvable, nil)

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

	bullet.Postion = math.Vector2{}
	bullet.CollisionShape = nil
	bullet.Speed.X = -10
	w.AddEntity(bullet)
	for bullet.Postion.X >= 0 {
		assert.True(t, s.Contains(bullet.CollisionShape))
		lastPostion := bullet.Postion
		w.Update(1)
		assert.Less(t, bullet.Postion.X, lastPostion.X)
	}
	w.Update(1)
}

type testDriver struct {
	axis                map[int]float64
	pressedKeys         map[ebiten.Key]int
	justPressedKeys     map[ebiten.Key]bool
	justReleasedKeys    map[ebiten.Key]bool
	pressedButtons      map[ebiten.GamepadButton]int
	justPressedButton   map[ebiten.GamepadButton]bool
	justReleasedButtons map[ebiten.GamepadButton]bool
}

func createTestDriver() *testDriver {
	return &testDriver{
		axis:                make(map[int]float64),
		pressedKeys:         make(map[ebiten.Key]int),
		justPressedKeys:     make(map[ebiten.Key]bool),
		justReleasedKeys:    make(map[ebiten.Key]bool),
		pressedButtons:      make(map[ebiten.GamepadButton]int),
		justPressedButton:   make(map[ebiten.GamepadButton]bool),
		justReleasedButtons: make(map[ebiten.GamepadButton]bool),
	}
}

func (t *testDriver) Ready(g *components.GamepadInputType) bool {
	return false
}

func (t *testDriver) GamepadAxis(id ebiten.GamepadID, axis int) float64 {
	return t.axis[axis]
}

func (t *testDriver) IsGamepadButtonJustPressed(id ebiten.GamepadID, btn ebiten.GamepadButton) bool {
	return t.justPressedButton[btn]
}

func (t *testDriver) KeyPressDuration(k ebiten.Key) int {
	return t.pressedKeys[k]
}

func (t *testDriver) IsKeyJustPressed(key ebiten.Key) bool {
	return t.justPressedKeys[key]
}

func (t *testDriver) IsKeyJustReleased(key ebiten.Key) bool {
	return t.justReleasedKeys[key]
}

func (t *testDriver) GamepadButtonPressDuration(_ ebiten.GamepadID, btn ebiten.GamepadButton) int {
	return t.pressedButtons[btn]
}

func (t *testDriver) IsGamepadButtonJustReleased(_ ebiten.GamepadID, btn ebiten.GamepadButton) bool {
	return t.justReleasedButtons[btn]
}

func TestInputSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	inputSystem := game.CreateInputSystem()
	var inputable *game.Inputable
	w.AddSystemInterface(inputSystem, inputable, nil)

	driver := createTestDriver()

	ent := &struct {
		ecs.BasicEntity
		*components.MovementComponent
		*components.InputComponent
	}{
		BasicEntity:       ecs.NewBasic(),
		MovementComponent: components.CreateMovementComponent(),
		InputComponent: &components.InputComponent{
			InputMode: components.InputModeGamepad,
			Gamepad:   components.DefaultGamepadInputType(),
			Keyboard:  components.DefaultKeyboardInputType(),
		},
	}
	ent.Keyboard.Driver = driver
	ent.Gamepad.Driver = driver
	w.AddEntity(ent)

	xAxisTestCases := []struct {
		xAxis  int
		target *int
		name   string
	}{
		{xAxis: 1, target: &ent.MovementComponent.PressedDuration[components.InputKindMoveRight], name: "Move right"},
		{xAxis: -1, target: &ent.MovementComponent.PressedDuration[components.InputKindMoveLeft], name: "Move left"},
	}
	for _, testCase := range xAxisTestCases {
		driver.axis[ent.Gamepad.MoveAxisX] = float64(testCase.xAxis)
		assert.Zero(t, *testCase.target, testCase.name)
		w.Update(1)
		assert.NotZero(t, *testCase.target, testCase.name)
	}

	ent.MovementComponent = components.CreateMovementComponent()

	// Button pressed
	for kind, btn := range ent.Gamepad.Mapping {
		assert.Zero(t, ent.PressedDuration[kind])
		assert.False(t, ent.JustPressed[kind])
		assert.False(t, ent.JustReleased[kind])

		driver.pressedButtons[btn] = 1
		driver.justPressedButton[btn] = true
		driver.justReleasedButtons[btn] = true

		ent.InputComponent.InputMode = components.InputModeGamepad
		w.Update(1)

		assert.NotZero(t, ent.PressedDuration[kind])
		assert.True(t, ent.JustPressed[kind])
		assert.True(t, ent.JustReleased[kind])

		driver.pressedButtons[btn] = 0
		driver.justPressedButton[btn] = false
		driver.justReleasedButtons[btn] = false

		ent.MovementComponent = components.CreateMovementComponent()
	}

	ent.InputMode = components.InputModeKeyboard

	// Key
	for kind, key := range ent.Keyboard.Mapping {
		if kind == 11 || kind == 12 {
			continue
		}

		assert.Zero(t, ent.PressedDuration[kind])
		assert.False(t, ent.JustPressed[kind])
		assert.False(t, ent.JustReleased[kind])

		driver.pressedKeys[key] = 1
		driver.justPressedKeys[key] = true
		driver.justReleasedKeys[key] = true

		ent.InputComponent.InputMode = components.InputModeKeyboard
		w.Update(1)

		assert.NotZerof(t, ent.PressedDuration[kind], "%v", kind)
		assert.Truef(t, ent.JustPressed[kind], "%v", kind)
		assert.Truef(t, ent.JustReleased[kind], "%v", kind)

		driver.pressedKeys[key] = 0
		driver.justPressedKeys[key] = false
		driver.justReleasedKeys[key] = false

		ent.MovementComponent = components.CreateMovementComponent()
	}

	w.RemoveEntity(ent.BasicEntity)

	assert.NotZero(t, inputSystem.Priority())
}

func TestLifeSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	lifeSystem := game.CreateLifeSystem()
	var lifeable *game.Lifeable
	w.AddSystemInterface(lifeSystem, lifeable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.LifeComponent
	}{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		LifeComponent:      &components.LifeComponent{},
	}

	testCases := []struct {
		StartingHP        float64
		InvincibilityTime time.Duration
		Damage            []float64
		ExpectedHp        []float64
	}{
		{StartingHP: 100, Damage: []float64{99, 1}, ExpectedHp: []float64{1, 0}},
		{StartingHP: 100, InvincibilityTime: 100 * time.Millisecond, Damage: []float64{99, 1, 1}, ExpectedHp: []float64{1, 1, 0}},
	}

	for _, testCase := range testCases {
		ent.HP = testCase.StartingHP
		ent.InvincibilityTime = testCase.InvincibilityTime
		w.AddEntity(ent)

		for i, damage := range testCase.Damage {
			ent.DamageEvents = append(ent.DamageEvents, &components.DamageEvent{
				Damage: damage,
			})
			w.Update(float32(100*time.Millisecond) / float32(time.Second))

			assert.Equalf(t, testCase.ExpectedHp[i], ent.HP, "loop %d %v", i+1, testCase)
		}

		w.RemoveEntity(ent.BasicEntity)
	}
}

func TestConstantSpeedSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	speedSystem := game.CreateConstantSpeedSystem()
	var speedable *game.ConstantSpeedable
	w.AddSystemInterface(speedSystem, speedable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.ConstantSpeedComponent
	}{
		BasicEntity:            ecs.NewBasic(),
		TransformComponent:     &components.TransformComponent{},
		ConstantSpeedComponent: &components.ConstantSpeedComponent{},
	}
	w.AddEntity(ent)

	testCases := []struct {
		StartingPostion math.Vector2
		Speed           math.Vector2
		ExpectedPostion math.Vector2
	}{
		{StartingPostion: math.Vector2{X: 1, Y: 0}, Speed: math.Vector2{X: 1, Y: -1}, ExpectedPostion: math.Vector2{X: 2, Y: -1}},
	}
	for _, testCase := range testCases {
		ent.Postion = testCase.ExpectedPostion
		ent.Speed = testCase.Speed

		w.Update(1)

		assert.Equal(t, testCase.ExpectedPostion, ent.Postion)
	}

}

func TestDestoryBoundSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	destorySystem := game.CreateDestoryBoundSystem()
	var destoryBoundable *game.DestoryBoundable
	w.AddSystemInterface(destorySystem, destoryBoundable, nil)

	ent := &struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.DestoryBoundComponent
	}{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		DestoryBoundComponent: &components.DestoryBoundComponent{
			Min: math.Vector2{},
			Max: math.Vector2{},
		},
	}

	testCases := []struct {
		Postion  math.Vector2
		BoundMin math.Vector2
		BoundMax math.Vector2
		Exists   bool
	}{
		{Postion: math.Vector2{}, BoundMin: math.Vector2{X: 10, Y: 10}, BoundMax: math.Vector2{X: 100, Y: 100}, Exists: false},
		{Postion: math.Vector2{X: 10, Y: 10}, BoundMin: math.Vector2{}, BoundMax: math.Vector2{X: 100, Y: 100}, Exists: true},
	}
	for _, testCase := range testCases {
		ent.Postion = testCase.Postion
		ent.DestoryBoundComponent.Max = testCase.BoundMax
		ent.DestoryBoundComponent.Min = testCase.BoundMin
		w.AddEntity(ent)

		w.Update(1)

		assert.Equal(t, destorySystem.ContainsEnt(ent.ID()), testCase.Exists)

		w.RemoveEntity(ent.BasicEntity)
	}
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
