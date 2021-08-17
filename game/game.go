package game

import (
	"container/heap"
	"image/color"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
)

const scaleMutiplier = 16
const gameWidth = 240 * scaleMutiplier
const gameHeight = 160 * scaleMutiplier

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

	var Velocityable *Velocityable
	world.AddSystemInterface(CreateVelocitySystem(), Velocityable, nil)

	var gameRuleable *GameRuleable
	world.AddSystemInterface(CreateGameRuleSystem(), gameRuleable, nil)
}

func (g *Game) startCityLevel() {
	g.world.AddEntity(entity.CreateCityMusic())

	cityBackground := entity.CreateCityBackground()
	cityBackground.GeoM.Scale(scaleMutiplier, scaleMutiplier)
	cityBackground.ImageComponent.Layer = bottomImageLayer
	g.world.AddEntity(cityBackground)

	player := entity.CreatePlayer()
	player.ImageComponent.Layer = middleImageLayer
	player.GeoM.Scale(scaleMutiplier, scaleMutiplier)
	g.world.AddEntity(player)

	testBox := entity.CreateTestBox()
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return gameWidth, gameHeight
}
