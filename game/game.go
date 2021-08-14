package game

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/systems"
)

type Game struct {
	world    *ecs.World
	lastTime time.Time
}

func addSystems(world *ecs.World) {
	var animeable *systems.Animeable
	world.AddSystemInterface(systems.CreateAnimeSystem(), animeable, nil)
	var renderable *systems.ImageRenderable
	world.AddSystemInterface(systems.CreateRenderSystem(), renderable, nil)
}

func CreateGame() *Game {
	world := &ecs.World{}
	addSystems(world)

	world.AddEntity(entity.Createplayer())

	return &Game{
		world:    world,
		lastTime: time.Now(),
	}
}

func (g *Game) Update() error {
	dt := time.Since(g.lastTime)
	g.world.Update(float32(dt / time.Millisecond))
	g.lastTime = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, system := range g.world.Systems() {
		if rendSys, ok := system.(systems.RenderingSystem); ok {
			rendSys.Render(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 3840, 2560
}
