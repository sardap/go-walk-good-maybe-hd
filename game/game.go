package game

import (
	"container/heap"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
)

const scaleMultiplier = 16
const gameWidth = 240 * scaleMultiplier
const gameHeight = 160 * scaleMultiplier

var (
	gGame *Game
)

const (
	bottomImageLayer components.ImageLayer = iota
	middleImageLayer
	uiImageLayer
	debugImageLayer
)

type Game struct {
	world    *ecs.World
	lastTime time.Time
}

func addSystems(world *ecs.World) {
	var collisionable *Collisionable
	world.AddSystemInterface(CreateCollisionSystem(), collisionable, nil)

	var animeable *Animeable
	world.AddSystemInterface(CreateAnimeSystem(), animeable, nil)

	var renderable *ImageRenderable
	world.AddSystemInterface(CreateImageRenderSystem(), renderable, nil)

	var textRenderable *TextRenderable
	world.AddSystemInterface(CreateTextRenderSystem(), textRenderable, nil)

	var inputable *Inputable
	world.AddSystemInterface(CreateInputSystem(), inputable, nil)

	var soundable *Soundable
	world.AddSystemInterface(CreateSoundSystem(), soundable, nil)

	var gameRuleable *GameRuleable
	world.AddSystemInterface(CreateGameRuleSystem(), gameRuleable, nil)

	var Velocityable *Velocityable
	world.AddSystemInterface(CreateVelocitySystem(), Velocityable, nil)
}

func (g *Game) startCityLevel() {
	g.world.AddEntity(entity.CreateCityMusic())

	cityBackground := entity.CreateCityBackground()
	cityBackground.GeoM.Scale(scaleMultiplier, scaleMultiplier)
	cityBackground.ImageComponent.Layer = bottomImageLayer
	g.world.AddEntity(cityBackground)

	cityBackground = entity.CreateCityBackground()
	cityBackground.GeoM.Scale(scaleMultiplier, scaleMultiplier)
	_, w, _, _ := bounds(cityBackground)
	cityBackground.GeoM.Translate(w, 0)
	cityBackground.ImageComponent.Layer = bottomImageLayer
	g.world.AddEntity(cityBackground)

	player := entity.CreatePlayer()
	player.ImageComponent.Layer = middleImageLayer
	player.GeoM.Scale(scaleMultiplier, scaleMultiplier)
	g.world.AddEntity(player)

	testBox := entity.CreateTestBox()
	testBox.ImageComponent.Layer = uiImageLayer
	testBox.Translate(500, 500)
	g.world.AddEntity(testBox)
}

func CreateGame() *Game {
	world := &ecs.World{}
	addSystems(world)

	result := &Game{
		world:    world,
		lastTime: time.Now(),
	}

	result.startCityLevel()

	gGame = result

	return result
}

func (g *Game) Update() error {
	dt := time.Since(g.lastTime)
	g.world.Update(float32(dt / time.Millisecond))
	g.lastTime = time.Now()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	queue := &RenderCmds{}

	screen.Fill(color.White)
	for _, system := range g.world.Systems() {
		if rendSys, ok := system.(RenderingSystem); ok {
			rendSys.Render(queue)
		}
	}

	heap.Init(queue)

	for queue.Len() > 0 {
		item := heap.Pop(queue).(*RenderCmd)
		screen.DrawImage(item.Image, item.Options)
	}

	img := ebiten.NewImage(50, 50)
	ebitenutil.DebugPrint(img, fmt.Sprintf("%2.f", ebiten.CurrentFPS()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 10)
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gameWidth, gameHeight
}
