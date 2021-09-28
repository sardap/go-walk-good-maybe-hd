package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sardap/walk-good-maybe-hd/assets"
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
	Start(*Game)
	End(*Game)
	Update(dt time.Duration, game *Game)
	Draw(screen *ebiten.Image)
}

type Game struct {
	audioCtx *audio.Context
	lastTime time.Time
	Info     *Info
	current  Scene
}

func (g *Game) ChangeScene(newScene Scene) {
	if g.current != nil {
		g.current.End(g)
	}
	g.current = newScene
	g.current.Start(g)
	g.lastTime = time.Unix(0, 0)
}

func CreateGame() *Game {
	result := &Game{
		audioCtx: audio.NewContext(48000),
	}

	return result
}

func (g *Game) Update() error {
	if g.current == nil {
		g.ChangeScene(&TitleScene{})
	}
	if g.lastTime.Unix() == 0 {
		g.lastTime = time.Now()
	}

	dt := time.Since(g.lastTime)
	g.current.Update(dt, g)
	g.lastTime = time.Now()

	if inpututil.IsKeyJustReleased(ebiten.KeyR) {
		g.ChangeScene(&TitleScene{})
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
	op.GeoM.Translate(0, 0)
	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func loadIconImage(inputCom *components.InputComponent, kind components.InputKind) *ebiten.Image {
	var result *ebiten.Image

	switch inputCom.InputMode {
	case components.InputModeGamepad:
		gamepad := inputCom.Gamepad
		result, _ = assets.LoadEbitenImageRaw([]byte(assets.IconXboxSeries[int(gamepad.Mapping[kind])]))
	case components.InputModeKeyboard:
		keyboard := inputCom.Keyboard
		result, _ = assets.LoadEbitenImageRaw([]byte(assets.IconKeyboardDark[keyboard.Mapping[kind]]))
	}

	return result
}
