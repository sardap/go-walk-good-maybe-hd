package game

import (
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
	"github.com/sardap/walk-good-maybe-hd/pkg/entity"
	"github.com/sardap/walk-good-maybe-hd/pkg/math"
	"github.com/sardap/walk-good-maybe-hd/pkg/utility"
)

type MainGameScene struct {
	Rand           *rand.Rand
	Space          *resolv.Space
	World          *ecs.World
	ScrollingSpeed math.Vector2
	Gravity        float64
	State          gameState
	Level          *Level
	InputEnt       *entity.InputEnt
	TimeElapsed    time.Duration
}

func (m *MainGameScene) addSystems(audioCtx *audio.Context) {
	var animeable *Animeable
	m.World.AddSystemInterface(CreateAnimeSystem(), animeable, nil)

	var renderable *ImageRenderable
	m.World.AddSystemInterface(CreateImageRenderSystem(), renderable, nil)

	var tileImageRenderable *TileImageRenderable
	m.World.AddSystemInterface(CreateTileImageRenderSystem(), tileImageRenderable, nil)

	var textRenderable *TextRenderable
	m.World.AddSystemInterface(CreateTextRenderSystem(), textRenderable, nil)

	var inputable *Inputable
	m.World.AddSystemInterface(CreateInputSystem(), inputable, nil)

	m.InputEnt = entity.CreateDebugInput()
	m.World.AddEntity(m.InputEnt)

	var soundable *Soundable
	m.World.AddSystemInterface(CreateSoundSystem(audioCtx), soundable, nil)

	var gameRuleable *GameRuleable
	m.World.AddSystemInterface(CreateGameRuleSystem(m), gameRuleable, nil)

	var velocityable *Velocityable
	m.World.AddSystemInterface(CreateVelocitySystem(m.Space), velocityable, nil)

	var dumbVelocityable *DumbVelocityable
	var exVelocityable *ExDumbVelocityable
	m.World.AddSystemInterface(CreateDumbVelocitySystem(), dumbVelocityable, exVelocityable)

	var resolvable *Resolvable
	m.World.AddSystemInterface(CreateResolvSystem(m.Space, m.InputEnt), resolvable, nil)

	var playerable *Playerable
	m.World.AddSystemInterface(CreatePlayerSystem(m), playerable, nil)

	var lifeable *Lifeable
	m.World.AddSystemInterface(CreateLifeSystem(), lifeable, nil)

	m.World.AddSystemInterface(CreateMainGameUiSystem(), gameRuleable, nil)

	var enemyBiscuitable *EnemyBiscuitable
	m.World.AddSystemInterface(CreateEnemyBiscuitSystem(m.Space), enemyBiscuitable, nil)
}

func (m *MainGameScene) addEnts() {
	m.World.AddEntity(entity.CreateCityMusic())

	cityBackground := entity.CreateCityBackground()
	cityBackground.ImageComponent.Layer = ImageLayerCityLayer
	m.World.AddEntity(cityBackground)

	cityBackground = entity.CreateCityBackground()
	cityBackground.ImageComponent.Layer = ImageLayerCityLayer
	cityBackground.Postion.X = cityBackground.TransformComponent.Size.X
	m.World.AddEntity(cityBackground)

	citySky := entity.CreateCitySkyBackground()
	citySky.Layer = ImageLayerBottom
	m.World.AddEntity(citySky)

	cityFog := entity.CreateCityFogBackground()
	cityFog.ImageComponent.Layer = ImageLayercityFogLayer
	m.World.AddEntity(cityFog)

	cityFog = entity.CreateCityFogBackground()
	cityFog.ImageComponent.Layer = ImageLayercityFogLayer
	cityFog.Postion.X = cityFog.TransformComponent.Size.X
	m.World.AddEntity(cityFog)

	player := entity.CreatePlayer()
	player.TileImageComponent.Layer = ImageLayerObjects
	player.MaxHp = 300
	player.HP = player.MaxHp
	player.JumpPower = startingPlayerJumpPower
	player.AirHorzSpeedModifier = startingPlayerAirHorzMod
	m.World.AddEntity(player)

	leftKillBox := entity.CreateKillBox()
	leftKillBox.Postion.X = -500 - leftKillBox.Size.X
	m.World.AddEntity(leftKillBox)

	rightKillBox := entity.CreateKillBox()
	rightKillBox.Postion.X = windowWidth + 1000 + rightKillBox.Size.X
	m.World.AddEntity(rightKillBox)

	bottomKillBox := entity.CreateKillBox()
	bottomKillBox.Size.X = windowWidth
	bottomKillBox.Postion.Y = windowHeight + 100
	m.World.AddEntity(bottomKillBox)
}

func (m *MainGameScene) Start(game *Game) {
	m.World = &ecs.World{}
	m.Space = resolv.NewSpace()
	m.ScrollingSpeed = math.Vector2{}
	m.Gravity = startingGravity
	m.Rand = rand.New(rand.NewSource(time.Now().Unix()))
	m.State = gameStateStarting
	m.Level = &Level{
		Width:  windowWidth,
		Height: windowHeight,
	}

	m.addSystems(game.audioCtx)
	m.addEnts()
}

func (m *MainGameScene) End(*Game) {
	for _, system := range m.World.Systems() {
		if soundSystem, ok := system.(*SoundSystem); ok {
			for _, ent := range soundSystem.ents {
				soundSystem.Remove(*ent.GetBasicEntity())
			}
		}
	}

	m.World = nil
	m.Space = nil
	m.ScrollingSpeed = math.Vector2{}
	m.Rand = nil
	m.Gravity = 0
	m.State = gameStateStarting
	m.TimeElapsed = 0
	m.Level = nil
	m.InputEnt = nil
}

func (m *MainGameScene) Update(dt time.Duration, _ *Game) {
	if m.InputEnt.InputPressed(components.InputKindFastGameSpeed) {
		dt *= 20
	}

	m.World.Update(float32(dt) / float32(time.Second))
	m.TimeElapsed += dt
}

func (m *MainGameScene) Draw(screen *ebiten.Image) {
	queue := RenderCmds{}
	for _, system := range m.World.Systems() {
		if render, ok := system.(RenderingSystem); ok {
			render.Render(&queue)
		}
	}

	queue.Sort()
	for _, item := range queue {
		item.Draw(screen)
	}
}

func (m *MainGameScene) GenerateCityBuildings() {
	x := m.Level.StartX
	for x < m.Level.Width {
		ent := ecs.NewBasic()
		levelBlock := createRandomLevelBlock(m.Rand, ent)
		trans := levelBlock.GetTransformComponent()
		levelBlock.GetTransformComponent().Postion.Y = m.Level.Height - trans.Size.Y
		levelBlock.GetTransformComponent().Postion.X = x
		x += levelBlock.Size.X + float64(utility.RandRange(m.Rand, minSpaceBetweenBuildings, minSpaceBetweenBuildings+20*scaleMultiplier))
		m.World.AddEntity(levelBlock)
		populateLevelBlock(m.Rand, m.World, levelBlock)
	}
	m.Level.StartX = x
}
