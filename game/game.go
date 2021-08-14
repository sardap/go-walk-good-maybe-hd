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
	g.world.Update(float32(dt * time.Second))
	g.lastTime = time.Now()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, system := range g.world.Systems() {
		if rendSys, ok := system.(systems.RenderingSystem); ok {
			rendSys.Render(screen)
		}
	}

	// op := &ebiten.DrawImageOptions{}
	// op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	// op.GeoM.Translate(screenWidth/2, screenHeight/2)
	// i := (g.count / 5) % frameNum
	// sx, sy := frameOX+i*frameWidth, frameOY
	// screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 3840, 2560
}
