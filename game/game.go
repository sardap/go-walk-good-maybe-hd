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
	"github.com/sardap/walk-good-maybe-hd/utility"
)

const (
	scaleMultiplier   = 16
	gameWidth         = 240 * scaleMultiplier
	gameHeight        = 160 * scaleMultiplier
	xStartScrollSpeed = -10.5
)

var (
	gGame *Game
)

const (
	bottomImageLayer components.ImageLayer = iota
	playerImageLayer
	buildingForground
	uiImageLayer
	debugImageLayer
)

type Game struct {
	world    *ecs.World
	lastTime time.Time
}

func (g *Game) addSystems() {
	world := g.world

	var collisionable *Collisionable
	world.AddSystemInterface(CreateCollisionSystem(), collisionable, nil)

	var animeable *Animeable
	world.AddSystemInterface(CreateAnimeSystem(), animeable, nil)

	var renderable *ImageRenderable
	world.AddSystemInterface(CreateImageRenderSystem(), renderable, nil)

	var tileImageRenderable *TileImageRenderable
	world.AddSystemInterface(CreateTileImageRenderSystem(), tileImageRenderable, nil)

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
	mainGameInfo = &MainGameInfo{}

	g.addSystems()

	g.world.AddEntity(entity.CreateCityMusic())

	cityBackground := entity.CreateCityBackground()
	cityBackground.ImageComponent.Layer = bottomImageLayer
	g.world.AddEntity(cityBackground)

	cityBackground = entity.CreateCityBackground()
	cityBackground.ImageComponent.Layer = bottomImageLayer
	cityBackground.Postion.X = cityBackground.TransformComponent.Size.X
	g.world.AddEntity(cityBackground)

	player := entity.CreatePlayer()
	player.ImageComponent.Layer = playerImageLayer
	g.world.AddEntity(player)

	testBox := entity.CreateTestBox()
	testBox.ImageComponent.Layer = uiImageLayer
	testBox.TransformComponent.Postion.X = 500
	testBox.TransformComponent.Postion.Y = 500
	g.world.AddEntity(testBox)

	ents := ecs.NewBasics(5)
	startX := float64(0)
	for _, ent := range ents {
		block := createBuilding0(ent)
		trans := block.GetTransformComponent()
		block.GetTransformComponent().Postion.Y = gameHeight/scaleMultiplier - trans.Size.Y
		block.GetTransformComponent().Postion.X = startX
		startX += block.GetTransformComponent().Size.X + float64(utility.RandRange(30, 50))
		g.world.AddEntity(block)
	}
}

func CreateGame() *Game {
	world := &ecs.World{}

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
	g.world.Update(float32(dt) / float32(time.Second))
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
		item := heap.Pop(queue).(RenderCmd)
		item.Draw(screen)
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
