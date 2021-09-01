package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
)

const (
	bottomImageLayer components.ImageLayer = iota
	playerImageLayer
	enemyLayer
	bulletImageLayer
	buildingForground
	uiImageLayer
	debugImageLayer
)

type Game struct {
	world        *ecs.World
	space        *resolv.Space
	lastTime     time.Time
	MainGameInfo *MainGameInfo
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
	world.AddSystemInterface(CreateGameRuleSystem(g.MainGameInfo, g.space), gameRuleable, nil)

	var velocityable *Velocityable
	world.AddSystemInterface(CreateVelocitySystem(g.space), velocityable, nil)

	var resolvable *Resolvable
	world.AddSystemInterface(CreateResolvSystem(g.space), resolvable, nil)
}

func (g *Game) startCityLevel() {
	g.MainGameInfo = &MainGameInfo{
		Gravity: startingGravity,
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

	g.MainGameInfo.Level = &Level{
		Width:  windowWidth,
		Height: windowHeight,
	}

	generateCityBuildings(g.MainGameInfo, g.world)
}

func (g *Game) Reset() {
	g.world = &ecs.World{}
	g.lastTime = time.Unix(0, 0)
	g.space = resolv.NewSpace()

	g.startCityLevel()
}

func CreateGame() *Game {
	result := &Game{}
	result.Reset()
	return result
}

func (g *Game) Update() error {
	if g.lastTime.Unix() == 0 {
		g.lastTime = time.Now()
	}

	dt := time.Since(g.lastTime)
	g.world.Update(float32(dt) / float32(time.Second))
	g.lastTime = time.Now()

	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.Reset()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	queue := RenderCmds{}

	screen.Fill(color.White)
	for _, system := range g.world.Systems() {
		if rendSys, ok := system.(RenderingSystem); ok {
			rendSys.Render(&queue)
		}
	}

	queue.Sort()

	for _, item := range queue {
		item.Draw(screen)
	}

	img := ebiten.NewImage(50, 50)
	ebitenutil.DebugPrint(img, fmt.Sprintf("%2.f", ebiten.CurrentFPS()))
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(10, 10)
	op.GeoM.Translate(float64(windowHeight-img.Bounds().Dx()), 0)
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}
