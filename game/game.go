package game

import (
	"fmt"
	"image/color"
	"math/rand"
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
	ImageLayerBottom components.ImageLayer = iota
	ImageLayerCityLayer
	ImageLayercityFogLayer
	ImageLayerplayer
	ImagelayerEnemyLayer
	ImagelayerToken
	ImageLayerbullet
	ImageLayerbuildingForground
	ImageLayerUi
	ImageLayerDebug
)

type Game struct {
	world    *ecs.World
	lastTime time.Time
	Info     *Info
}

func (g *Game) addSystems() {
	world := g.world

	mainGameInfo := g.Info.MainGameInfo

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

	mainGameInfo.InputEnt = entity.CreateDebugInput()
	world.AddEntity(mainGameInfo.InputEnt)

	var soundable *Soundable
	world.AddSystemInterface(CreateSoundSystem(), soundable, nil)

	var gameRuleable *GameRuleable
	world.AddSystemInterface(CreateGameRuleSystem(g.Info), gameRuleable, nil)

	var velocityable *Velocityable
	world.AddSystemInterface(CreateVelocitySystem(g.Info.Space), velocityable, nil)

	var dumbVelocityable *DumbVelocityable
	var exVelocityable *ExDumbVelocityable
	world.AddSystemInterface(CreateDumbVelocitySystem(), dumbVelocityable, exVelocityable)

	var resolvable *Resolvable
	world.AddSystemInterface(CreateResolvSystem(mainGameInfo, g.Info.Space), resolvable, nil)

	var playerable *Playerable
	world.AddSystemInterface(CreatePlayerSystem(mainGameInfo), playerable, nil)

	var lifeable *Lifeable
	world.AddSystemInterface(CreateLifeSystem(), lifeable, nil)

	world.AddSystemInterface(CreateMainGameUiSystem(), gameRuleable, nil)

	var enemyBiscuitable *EnemyBiscuitable
	world.AddSystemInterface(CreateEnemyBiscuitSystem(g.Info.Space), enemyBiscuitable, nil)

}

func (g *Game) startCityLevel() {
	g.Info.MainGameInfo = &MainGameInfo{
		Gravity: startingGravity,
	}

	g.addSystems()

	g.world.AddEntity(entity.CreateCityMusic())

	cityBackground := entity.CreateCityBackground()
	cityBackground.ImageComponent.Layer = ImageLayerCityLayer
	g.world.AddEntity(cityBackground)

	cityBackground = entity.CreateCityBackground()
	cityBackground.ImageComponent.Layer = ImageLayerCityLayer
	cityBackground.Postion.X = cityBackground.TransformComponent.Size.X
	g.world.AddEntity(cityBackground)

	citySky := entity.CreateCitySkyBackground()
	citySky.Layer = ImageLayerBottom
	g.world.AddEntity(citySky)

	cityFog := entity.CreateCityFogBackground()
	cityFog.ImageComponent.Layer = ImageLayercityFogLayer
	g.world.AddEntity(cityFog)

	cityFog = entity.CreateCityFogBackground()
	cityFog.ImageComponent.Layer = ImageLayercityFogLayer
	cityFog.Postion.X = cityFog.TransformComponent.Size.X
	g.world.AddEntity(cityFog)

	player := entity.CreatePlayer()
	player.TileImageComponent.Layer = ImageLayerplayer
	player.MaxHp = 300
	player.HP = player.MaxHp
	player.JumpPower = startingPlayerJumpPower
	player.AirHorzSpeedModifier = startingPlayerAirHorzMod
	g.world.AddEntity(player)

	leftKillBox := entity.CreateKillBox()
	leftKillBox.Postion.X = -500 - leftKillBox.Size.X
	g.world.AddEntity(leftKillBox)

	rightKillBox := entity.CreateKillBox()
	rightKillBox.Postion.X = windowWidth + 1000 + rightKillBox.Size.X
	g.world.AddEntity(rightKillBox)

	bottomKillBox := entity.CreateKillBox()
	bottomKillBox.Size.X = windowWidth
	bottomKillBox.Postion.Y = windowHeight + 100
	g.world.AddEntity(bottomKillBox)

	g.Info.MainGameInfo.Level = &Level{
		Width:  windowWidth,
		Height: windowHeight,
	}

	generateCityBuildings(g.Info, g.world)
}

func (g *Game) Reset() {
	g.world = &ecs.World{}
	g.lastTime = time.Unix(0, 0)
	g.Info.Space = resolv.NewSpace()

	g.startCityLevel()
}

func CreateGame() *Game {
	result := &Game{
		Info: &Info{
			Rand: rand.New(rand.NewSource(time.Now().Unix())),
		},
	}
	result.Reset()
	return result
}

func (g *Game) Update() error {
	if g.lastTime.Unix() == 0 {
		g.lastTime = time.Now()
	}

	dt := time.Since(g.lastTime)
	if g.Info.MainGameInfo.InputEnt.FastGameSpeed {
		dt *= 20
		g.Info.MainGameInfo.InputEnt.FastGameSpeed = false
	}
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
