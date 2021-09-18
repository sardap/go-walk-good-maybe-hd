package game

import (
	"bytes"
	"image"
	"image/color"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type TitleScene struct {
	titleText  *ebiten.Image
	beach      *ebiten.Image
	city       *ebiten.Image
	citySky    *ebiten.Image
	cityFog    *ebiten.Image
	xOffset    float64
	xFogOffset float64
	world      *ecs.World
	inputEnt   *entity.InputEnt
	player     *audio.Player
}

func (s *TitleScene) Start(game *Game) {
	img, _ := assets.LoadEbitenImage(assets.ImageTitleSceneText)
	s.titleText = img

	img, _ = assets.LoadEbitenImage(assets.ImageTitleSceneBeach)
	s.beach = img

	img, _ = assets.LoadEbitenImage(assets.ImageBackgroundCity)
	s.city = img

	img, _ = assets.LoadEbitenImage(assets.ImageCityFog)
	s.cityFog = img

	img, _ = assets.LoadEbitenImage(assets.ImageSkyCity)
	s.citySky = img

	s.xOffset = 0
	s.xFogOffset = 0

	s.world = &ecs.World{}

	var inputable *Inputable
	s.world.AddSystemInterface(CreateInputSystem(), inputable, nil)

	s.inputEnt = entity.CreateMenuInput()
	s.world.AddEntity(s.inputEnt)

	sound := components.LoadSound(assets.MusicPdTitleScreen)
	buffer := bytes.NewReader(sound.Source)
	stream, _ := mp3.DecodeWithSampleRate(sound.SampleRate, buffer)
	s.player, _ = audio.NewPlayer(game.audioCtx, stream)
	s.player.Play()
}

func (s *TitleScene) End(*Game) {
	s.world = nil
	s.inputEnt = nil
	s.player.Pause()
	s.player = nil
}

func (s *TitleScene) Update(dt time.Duration, game *Game) {
	s.xOffset = utility.WrapFloat64(s.xOffset+1, 0, float64(s.city.Bounds().Dx()/2))
	s.xFogOffset = utility.WrapFloat64(s.xFogOffset+0.5, 0, float64(s.city.Bounds().Dx()/2))

	s.world.Update(float32(dt) / float32(time.Second))

	if s.inputEnt.MovementComponent.Select {
		game.ChangeScene(&MainGameScene{})
	}
}

func (s *TitleScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{B: 255, A: 255})

	var op *ebiten.DrawImageOptions

	subRect := image.Rectangle{
		Min: image.Pt(int(s.xOffset), 0),
		Max: image.Pt(int(s.xOffset+windowWidth/2), windowHeight),
	}

	op = &ebiten.DrawImageOptions{}
	screen.DrawImage(s.citySky.SubImage(subRect).(*ebiten.Image), op)
	screen.DrawImage(s.city.SubImage(subRect).(*ebiten.Image), op)
	screen.DrawImage(s.cityFog.SubImage(image.Rectangle{
		Min: image.Pt(int(s.xFogOffset), 0),
		Max: image.Pt(int(s.xFogOffset)+windowWidth/2, windowHeight),
	}).(*ebiten.Image), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(s.beach.Bounds().Dx())-160, 0)
	screen.DrawImage(s.beach.SubImage(subRect).(*ebiten.Image), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(windowWidth/2-float64(s.titleText.Bounds().Dx()/2), 0)
	screen.DrawImage(s.titleText, op)
}
