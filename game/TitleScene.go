package game

import (
	"bytes"
	"image"
	"image/color"
	"io"
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

type MenuItem struct {
	TargetScene Scene
	Text        *ebiten.Image
}

type TitleScene struct {
	titleText          *ebiten.Image
	beach              *ebiten.Image
	beachWater         *ebiten.Image
	beachWaterAnimeIdx int
	waterAnimeTimer    time.Duration
	city               *ebiten.Image
	citySky            *ebiten.Image
	cityFog            *ebiten.Image
	xOffset            float64
	xFogOffset         float64
	world              *ecs.World
	inputEnt           *entity.InputEnt
	player             *audio.Player
	menuItems          []MenuItem
	selectedIdx        int
	// Arrow's
	selectionArrowCooldown time.Duration
	selectionActiveArrow   *ebiten.Image
	whiteArrow             *ebiten.Image
	redArrow               *ebiten.Image
}

func (s *TitleScene) Start(game *Game) {
	img, _ := assets.LoadEbitenImage(assets.ImageTitleSceneText)
	s.titleText = img

	img, _ = assets.LoadEbitenImage(assets.ImageTitleSceneBeach)
	s.beach = img

	img, _ = assets.LoadEbitenImage(assets.ImageTitleSceneBeachWaterTileSet)
	s.beachWater = img
	s.beachWaterAnimeIdx = 0
	s.waterAnimeTimer = 0

	img, _ = assets.LoadEbitenImage(assets.ImageBackgroundCity)
	s.city = img

	img, _ = assets.LoadEbitenImage(assets.ImageCityFog)
	s.cityFog = img

	img, _ = assets.LoadEbitenImage(assets.ImageSkyCity)
	s.citySky = img

	s.xOffset = 0
	s.xFogOffset = 0

	img, _ = assets.LoadEbitenImageColorSwap(
		assets.ImageTitleSceneArrow,
		map[color.RGBA]color.RGBA{
			swapColor: {R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
		},
	)
	s.whiteArrow = img

	img, _ = assets.LoadEbitenImageColorSwap(
		assets.ImageTitleSceneArrow,
		map[color.RGBA]color.RGBA{
			swapColor: {R: 0xFF, G: 0x00, B: 0x00, A: 0xFF},
		},
	)
	s.redArrow = img

	s.selectionActiveArrow = s.whiteArrow
	s.selectionArrowCooldown = 0

	img, _ = assets.LoadEbitenImage(assets.ImageTitleSceneGameText)
	s.menuItems = []MenuItem{
		{
			TargetScene: &MainGameScene{},
			Text:        img,
		},
		{
			TargetScene: &TitleScene{},
			Text:        img,
		},
	}
	s.selectedIdx = 0

	s.world = &ecs.World{}

	var inputable *Inputable
	s.world.AddSystemInterface(CreateInputSystem(), inputable, nil)

	s.inputEnt = entity.CreateMenuInput()
	s.world.AddEntity(s.inputEnt)

	sound := components.LoadSound(assets.MusicPdTitleScreen)
	buffer := bytes.NewReader(sound.Source)
	var stream io.Reader
	stream, _ = mp3.DecodeWithSampleRate(sound.SampleRate, buffer)
	mp3Stream := stream.(*mp3.Stream)
	stream = audio.NewInfiniteLoop(mp3Stream, mp3Stream.Length())
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

	s.waterAnimeTimer += dt
	if s.waterAnimeTimer > 200*time.Millisecond {
		count := s.beachWater.Bounds().Dx() / assets.ImageTitleSceneBeachWaterTileSet.FrameWidth
		s.beachWaterAnimeIdx = utility.WrapInt(s.beachWaterAnimeIdx+1, 0, count)
		s.waterAnimeTimer = 0
	}

	if s.selectionArrowCooldown > 0 {
		s.selectionArrowCooldown -= dt
	} else {
		s.selectionActiveArrow = s.whiteArrow
	}

	// Input's
	if s.inputEnt.InputJustPressed(components.InputKindSelect) {
		defer game.ChangeScene(s.menuItems[s.selectedIdx].TargetScene)
	}

	if s.inputEnt.InputJustPressed(components.InputKindMoveUp) {
		s.selectedIdx = utility.WrapInt(s.selectedIdx-1, 0, len(s.menuItems))
		s.selectionActiveArrow = s.redArrow
		s.selectionArrowCooldown = 125 * time.Millisecond
	}

	if s.inputEnt.InputJustPressed(components.InputKindMoveDown) {
		s.selectedIdx = utility.WrapInt(s.selectedIdx+1, 0, len(s.menuItems))
		s.selectionActiveArrow = s.redArrow
		s.selectionArrowCooldown = 125 * time.Millisecond
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

	xStart := s.beachWaterAnimeIdx * assets.ImageTitleSceneBeachWaterTileSet.FrameWidth
	beachWater := s.beachWater.SubImage(image.Rectangle{
		Min: image.Pt(xStart, 0),
		Max: image.Pt(int(xStart+assets.ImageTitleSceneBeachWaterTileSet.FrameWidth), windowHeight),
	}).(*ebiten.Image)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(beachWater.Bounds().Dx())-160, 107*10)
	beachWater = beachWater.SubImage(image.Rectangle{
		Min: image.Pt(beachWater.Bounds().Min.X+int(s.xOffset), 0),
		Max: image.Pt(beachWater.Bounds().Min.X+int(s.xOffset+windowWidth/2), windowHeight),
	}).(*ebiten.Image)
	screen.DrawImage(beachWater, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(s.beach.Bounds().Dx())-160, 0)
	screen.DrawImage(s.beach.SubImage(subRect).(*ebiten.Image), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(windowWidth/2-float64(s.titleText.Bounds().Dx()/2), 0)
	screen.DrawImage(s.titleText, op)

	textXStart := float64(windowWidth/2 - 150)
	yStart := float64(900)
	for _, item := range s.menuItems {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(textXStart, yStart)
		screen.DrawImage(item.Text, op)
		yStart += 130
	}

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(textXStart-float64(s.selectionActiveArrow.Bounds().Dx())-10, 900+float64(s.selectedIdx*130))
	screen.DrawImage(s.selectionActiveArrow, op)
}
