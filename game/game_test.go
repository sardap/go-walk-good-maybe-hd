package game_test

import (
	"bytes"
	"image"
	"math/rand"
	"testing"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
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
}

func TestVelocitySystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}
	s := resolv.NewSpace()

	// Setup
	velocitySystem := game.CreateVelocitySystem(s)

	var velocityable *game.Velocityable
	w.AddSystemInterface(velocitySystem, velocityable, nil)

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

	// asserts
	assert.False(t, s.Contains(entA.CollisionShape), "space should be empty")
	w.AddEntity(entA)
	assert.NotNil(t, entA.CollisionShape, "collsion shape should be init")
	assert.True(t, s.Contains(entA.CollisionShape), "space should contain shape")
	assert.Equal(t, "test", entA.CollisionShape.GetTags()[0], "tags should be copied to shape")

	w.RemoveEntity(entA.BasicEntity)
	assert.False(t, s.Contains(entA.CollisionShape), "shape should be removed from space")

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
}

func TestImageRenderSystem(t *testing.T) {
	t.Parallel()

	w := &ecs.World{}

	imageRenderSystem := game.CreateImageRenderSystem()

	var imageRenderable *game.ImageRenderable
	w.AddSystemInterface(imageRenderSystem, imageRenderable, nil)

	renderQueue := make(game.RenderCmds, 0)

	ents := []*struct {
		ecs.BasicEntity
		*components.TransformComponent
		*components.ImageComponent
		clr color.Color
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
			clr color.Color
		}{
			BasicEntity:        basics[i],
			TransformComponent: &components.TransformComponent{},
			ImageComponent: &components.ImageComponent{
				Active: true,
				Image:  img,
				Layer:  components.ImageLayer(i),
			},
			clr: clr,
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

	assert.Equal(t, len(renderQueue), len(ents))

	screen := ebiten.NewImage(500, 500)

	last := renderQueue[len(renderQueue)-1]
	for i := len(renderQueue) - 2; i > 0; i-- {
		currentEnt := ents[i]
		current := renderQueue[i]
		assert.Equal(t, int(currentEnt.Layer), current.GetLayer())
		assert.Greater(t, last.GetLayer(), current.GetLayer())
		current.Draw(screen)
		last = current
	}
}
