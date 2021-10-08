package game

import (
	"bytes"
	"image"
	"image/color"
	"io"
	"log"
	"sync"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sardap/walk-good-maybe-hd/pkg/assets"
	"github.com/sardap/walk-good-maybe-hd/pkg/common"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
	"github.com/sardap/walk-good-maybe-hd/pkg/entity"
	"github.com/sardap/walk-good-maybe-hd/pkg/utility"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type MenuItem interface {
	Action(*Game, *TitleScene)
	GetIcon() *ebiten.Image
}

type SceneMenuItem struct {
	TargetScene Scene
	Text        *ebiten.Image
}

func (s *SceneMenuItem) GetIcon() *ebiten.Image {
	return s.Text
}

func (s *SceneMenuItem) Action(g *Game, t *TitleScene) {
	g.ChangeScene(s.TargetScene)
}

type KaraokeMenuItem struct {
	KaraokeScene Scene
	Text         *ebiten.Image
}

func (k *KaraokeMenuItem) GetIcon() *ebiten.Image {
	return k.Text
}

func (k *KaraokeMenuItem) Action(_ *Game, s *TitleScene) {
	s.karaokeLoadingLock.Lock()
	defer s.karaokeLoadingLock.Unlock()
	s.selectedIdx = 0
	s.state = TitleSceneStateKaraoke
}

type TitleSceneState int

const (
	TitleSceneStateMainMenu TitleSceneState = iota
	TitleSceneStateKaraoke
)

type TitleScene struct {
	state              TitleSceneState
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
	// Karaoke
	karaokeLoadingLock *sync.Mutex
	karaokeIdx         *common.KaraokeIndex

	font font.Face
}

func (s *TitleScene) Start(game *Game) {
	s.state = TitleSceneStateMainMenu

	s.karaokeLoadingLock = &sync.Mutex{}

	s.karaokeLoadingLock.Lock()
	go func() {
		defer s.karaokeLoadingLock.Unlock()
		s.karaokeIdx = assets.LoadKaraokeIndex()

		tt, err := opentype.Parse(assets.FontSignikaRegular)
		if err != nil {
			log.Fatal(err)
		}
		s.font, _ = opentype.NewFace(tt, &opentype.FaceOptions{
			Size:    80,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	}()

	img, _ := assets.LoadEbitenImageAsset(assets.ImageTitleSceneText)
	s.titleText = img

	img, _ = assets.LoadEbitenImageAsset(assets.ImageTitleSceneBeach)
	s.beach = img

	img, _ = assets.LoadEbitenImageAsset(assets.ImageTitleSceneBeachWaterTileSet)
	s.beachWater = img
	s.beachWaterAnimeIdx = 0
	s.waterAnimeTimer = 0

	img, _ = assets.LoadEbitenImageAsset(assets.ImageBackgroundCity)
	s.city = img

	img, _ = assets.LoadEbitenImageAsset(assets.ImageCityFog)
	s.cityFog = img

	img, _ = assets.LoadEbitenImageAsset(assets.ImageSkyCity)
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

	img, _ = assets.LoadEbitenImageAsset(assets.ImageTitleSceneGameText)
	s.menuItems = []MenuItem{
		&SceneMenuItem{
			TargetScene: &MainGameScene{},
			Text:        img,
		},
		&SceneMenuItem{
			TargetScene: &TitleScene{},
			Text:        img,
		},
		&KaraokeMenuItem{
			Text: img,
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

func (s *TitleScene) Update(dt time.Duration, g *Game) {
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

	switch s.state {
	case TitleSceneStateMainMenu:
		// Input's
		if s.inputEnt.InputJustPressed(components.InputKindSelect) {
			defer s.menuItems[s.selectedIdx].Action(g, s)
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
	case TitleSceneStateKaraoke:
		// Input's
		if s.inputEnt.InputJustPressed(components.InputKindSelect) {
			defer func() {
				selectedSession := s.karaokeIdx.KaraokeGames[s.selectedIdx]
				g.ChangeScene(&KaraokeScene{
					Session: assets.LoadKaraokeSession(selectedSession),
				})
			}()
		}

		if s.inputEnt.InputJustPressed(components.InputKindMoveUp) {
			s.selectedIdx = utility.WrapInt(s.selectedIdx-1, 0, len(s.karaokeIdx.KaraokeGames))
			s.selectionActiveArrow = s.redArrow
			s.selectionArrowCooldown = 125 * time.Millisecond
		}

		if s.inputEnt.InputJustPressed(components.InputKindMoveDown) {
			s.selectedIdx = utility.WrapInt(s.selectedIdx+1, 0, len(s.karaokeIdx.KaraokeGames))
			s.selectionActiveArrow = s.redArrow
			s.selectionArrowCooldown = 125 * time.Millisecond
		}
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
	op.ColorM.Reset()
	op.GeoM.Reset()

	xStart := s.beachWaterAnimeIdx * assets.ImageTitleSceneBeachWaterTileSet.FrameWidth
	beachWater := s.beachWater.SubImage(image.Rectangle{
		Min: image.Pt(xStart, 0),
		Max: image.Pt(int(xStart+assets.ImageTitleSceneBeachWaterTileSet.FrameWidth), windowHeight),
	}).(*ebiten.Image)

	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(beachWater.Bounds().Dx())+125, 107*10)
	beachWater = beachWater.SubImage(image.Rectangle{
		Min: image.Pt(beachWater.Bounds().Min.X+int(s.xOffset), 0),
		Max: image.Pt(beachWater.Bounds().Min.X+int(s.xOffset+windowWidth/2), windowHeight),
	}).(*ebiten.Image)
	screen.DrawImage(beachWater, op)
	op.ColorM.Reset()
	op.GeoM.Reset()

	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(s.beach.Bounds().Dx())+125, 0)
	screen.DrawImage(s.beach.SubImage(subRect).(*ebiten.Image), op)
	op.ColorM.Reset()
	op.GeoM.Reset()

	op.GeoM.Translate(windowWidth/2-float64(s.titleText.Bounds().Dx()/2), 0)
	screen.DrawImage(s.titleText, op)
	op.ColorM.Reset()
	op.GeoM.Reset()

	switch s.state {
	case TitleSceneStateMainMenu:
		{
			textXStart := float64(windowWidth/2 - 150)
			yStart := float64(900)
			for _, item := range s.menuItems {
				op.GeoM.Translate(textXStart, yStart)
				screen.DrawImage(item.GetIcon(), op)
				op.ColorM.Reset()
				op.GeoM.Reset()
				yStart += 130
			}

			op.GeoM.Translate(textXStart-float64(s.selectionActiveArrow.Bounds().Dx())-10, 900+float64(s.selectedIdx*130))
			screen.DrawImage(s.selectionActiveArrow, op)
			op.ColorM.Reset()
			op.GeoM.Reset()
		}
	case TitleSceneStateKaraoke:
		{
			textXStart := float64(windowWidth/2 - 150)
			yStart := float64(900)
			for _, karaokeGame := range s.karaokeIdx.KaraokeGames {
				op.GeoM.Translate(textXStart, yStart)
				text.Draw(screen, karaokeGame, s.font, int(textXStart), int(yStart), color.White)
				op.ColorM.Reset()
				op.GeoM.Reset()
				yStart += 130
			}

			op.GeoM.Translate(textXStart-float64(s.selectionActiveArrow.Bounds().Dx())-10, 850+float64(s.selectedIdx*130))
			screen.DrawImage(s.selectionActiveArrow, op)
			op.ColorM.Reset()
			op.GeoM.Reset()
		}
	}
}
