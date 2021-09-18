package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sardap/walk-good-maybe-hd/components"
)

const (
	ImageLayerBottom components.ImageLayer = iota
	ImageLayerCityLayer
	ImageLayercityFogLayer
	ImageLayerObjects
	ImagelayerToken
	ImageLayerbullet
	ImageLayerbuildingForground
	ImageLayerUi
	ImageLayerDebug
)

type Scene interface {
	Start(*Info)
	End(*Info)
	Update(dt time.Duration, info *Info)
	Draw(screen *ebiten.Image)
}

type Game struct {
	lastTime time.Time
	Info     *Info
	current  Scene
}

func (g *Game) ChangeScene(newScene Scene) {
	if g.current != nil {
		g.current.End(g.Info)
	}
	g.current = newScene
	g.current.Start(g.Info)
	g.lastTime = time.Unix(0, 0)
}

func CreateGame() *Game {
	result := &Game{}
	result.ChangeScene(&MainGameScene{})

	return result
}

func (g *Game) Update() error {
	if g.lastTime.Unix() == 0 {
		g.lastTime = time.Now()
	}

	dt := time.Since(g.lastTime)
	g.current.Update(dt, g.Info)
	g.lastTime = time.Now()

	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.ChangeScene(g.current)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.White)
	g.current.Draw(screen)

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
