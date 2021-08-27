package game

import (
	"container/heap"
	"fmt"
	"image/color"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
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
	space    *resolv.Space
	lastTime time.Time
}

func (g *Game) addSystems() {
	world := g.world

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
	world.AddSystemInterface(CreateGameRuleSystem(g.space), gameRuleable, nil)

	var Velocityable *Velocityable
	world.AddSystemInterface(CreateVelocitySystem(g.space), Velocityable, nil)
}

func (g *Game) startCityLevel() {
	mainGameInfo = &MainGameInfo{
		gravity: startingGravity,
	}

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
	player.TileImageComponent.Layer = playerImageLayer
	g.world.AddEntity(player)

	testBox := entity.CreateTestBox()
	testBox.ImageComponent.Layer = uiImageLayer
	testBox.TransformComponent.Postion.X = 500
	testBox.TransformComponent.Postion.Y = 500
	g.world.AddEntity(testBox)

	mainGameInfo.level = &Level{}

	generateBuildings(g.world)
}

func CreateGame() *Game {
	world := &ecs.World{}

	result := &Game{
		world:    world,
		lastTime: time.Unix(0, 0),
		space:    resolv.NewSpace(),
	}

	result.startCityLevel()

	return result
}

func (g *Game) Update() error {
	if g.lastTime.Unix() == 0 {
		g.lastTime = time.Now()
	}

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
