package game

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	gomath "math"
	"math/rand"
	"strconv"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/icza/gox/imagex/colorx"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ImageLayerKaraokeBack components.ImageLayer = iota
	ImageLayerKaraokeFront
	ImageLayerKaraokeObjects
	ImageLayerKaraokeUi
	ImageLayerKaraokeText
)

type DurationMil time.Duration

func (d *DurationMil) UnmarshalJSON(b []byte) error {
	var v time.Duration
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*d = DurationMil(v * time.Millisecond)

	return nil
}

type KaraokeInput struct {
	StartTime  DurationMil             `json:"start_time"`
	Duration   DurationMil             `json:"duration"`
	Sound      components.KaraokeSound `json:"sound"`
	xPostion   float64
	xSpeed     float64
	hitPostion float64
}

func (k *KaraokeInput) Y() float64 {
	switch k.Sound {
	case components.KaraokeSoundA:
		return 550
	case components.KaraokeSoundB:
		return 450
	case components.KaraokeSoundX:
		return 350
	case components.KaraokeSoundY:
		return 250
	}

	return 400
}

const (
	karaBoundStep     = 100
	karaCenter        = 850
	karaLeftBound     = karaCenter - karaBoundStep
	karaRightBound    = karaCenter + karaBoundStep
	karaScoreSpinTime = 2*time.Second + 500*time.Millisecond
)

type KaraokeScore int

func (k KaraokeScore) String() string {
	switch {
	case k < 25:
		return "Okay"
	case k < 50:
		return "Good"
	case k < 75:
		return "Great"
	}

	return "Perfect"
}

const (
	KaraokeScoreOkay    KaraokeScore = 10
	KaraokeScoreGood    KaraokeScore = 20
	KaraokeScoreGreat   KaraokeScore = 40
	KaraokeScorePerfect KaraokeScore = 80
)

func Score(x float64) KaraokeScore {
	if x == 0 {
		return 0
	}

	delta := gomath.Abs((x + 50) - (windowWidth - karaCenter))

	switch {
	case delta < 25:
		return KaraokeScoreOkay
	case delta < 50:
		return KaraokeScoreGood
	case delta < 75:
		return KaraokeScoreGreat
	}

	return KaraokeScorePerfect
}

type KaraokeBackground struct {
	Duration DurationMil `json:"duration"`
	FadeIn   DurationMil `json:"fade_in"`
	Image    string      `json:"image"`
}

type KaraokeSession struct {
	Inputs        []*KaraokeInput      `json:"inputs"`
	Backgrounds   []*KaraokeBackground `json:"backgrounds"`
	Sounds        map[string]string    `json:"sounds"`
	Music         string               `json:"music"`
	SampleRate    int                  `json:"sampleRate"`
	backgroundIdx int
}

type KaraokeState int

const (
	KaraokeStateStarting KaraokeState = iota
	KaraokeStateSinging
	KaraokeStateComplete
)

type karaokeInfo struct {
	sound *entity.KaraokeInputSound
	image *ebiten.Image
	input components.InputKind
}

type KaraokeScene struct {
	Session     *KaraokeSession
	inputLeeway time.Duration

	rand  *rand.Rand
	world *ecs.World

	inputEnt      *entity.InputEnt
	musicEnt      *entity.SoundPlayer
	karaokePlayer *entity.KaraokePlayer

	currentImage          *ebiten.Image
	nextImage             *ebiten.Image
	backgroundFront       *entity.BasicImage
	backgroundFrontFadeIn time.Duration
	backgroundBack        *entity.BasicImage
	ui                    *entity.BasicImage

	scorePlayer      *entity.SoundPlayer
	scoreTitleFont   font.Face
	scoreFont        font.Face
	scoreColors      []color.Color
	nextColorUpdate  time.Duration
	activeScoreColor color.Color
	textScreen       *ebiten.Image
	scoreOpacity     float64

	comboFont font.Face

	timeElapsed       time.Duration
	backgroundElapsed time.Duration
	state             KaraokeState

	soundInfo map[components.KaraokeSound]*karaokeInfo
}

func (k *KaraokeScene) addSystems(audioCtx *audio.Context) {
	var soundable *Soundable
	k.world.AddSystemInterface(CreateSoundSystem(audioCtx), soundable, nil)

	var renderable *ImageRenderable
	k.world.AddSystemInterface(CreateImageRenderSystem(), renderable, nil)

	var inputable *Inputable
	k.world.AddSystemInterface(CreateInputSystem(), inputable, nil)

	var constantSpeedable *ConstantSpeedable
	k.world.AddSystemInterface(CreateConstantSpeedSystem(), constantSpeedable, nil)

	var dumbVelocityable *DumbVelocityable
	k.world.AddSystemInterface(CreateDumbVelocitySystem(), dumbVelocityable, nil)

	var textRenderable *TextRenderable
	k.world.AddSystemInterface(CreateTextRenderSystem(), textRenderable, nil)

	var destoryBoundable *DestoryBoundable
	k.world.AddSystemInterface(CreateDestoryBoundSystem(), destoryBoundable, nil)
}

func (k *KaraokeScene) addEnts() {
	rawMusic, _ := base64.StdEncoding.DecodeString(k.Session.Music)
	k.musicEnt = &entity.SoundPlayer{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		SoundComponent: &components.SoundComponent{
			Sound: &components.Sound{
				Source:     rawMusic,
				SampleRate: k.Session.SampleRate,
				SoundType:  assets.SoundTypeMp3,
			},
		},
	}
	k.musicEnt.SoundComponent.Active = false
	k.world.AddEntity(k.musicEnt)

	k.backgroundFront = &entity.BasicImage{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageUiLifeAmountTileSet.FrameWidth),
				Y: float64(assets.ImageUiLifeAmountTileSet.FrameWidth),
			},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Layer:  ImageLayerKaraokeFront,
			Options: components.DrawOptions{
				Opacity: 0,
			},
		},
	}
	k.world.AddEntity(k.backgroundFront)

	k.backgroundBack = &entity.BasicImage{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageUiLifeAmountTileSet.FrameWidth),
				Y: float64(assets.ImageUiLifeAmountTileSet.FrameWidth),
			},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Layer:  ImageLayerKaraokeBack,
		},
	}
	k.world.AddEntity(k.backgroundBack)

	k.ui = entity.CreateBasicImage(assets.ImageKaraokeBackground)
	k.ui.Layer = ImageLayerKaraokeUi
	k.world.AddEntity(k.ui)

	k.karaokePlayer = entity.CreateKaraokePlayer()
	k.karaokePlayer.ImageComponent.Layer = ImageLayerKaraokeObjects
	k.karaokePlayer.Postion.Y = windowHeight / 2
	k.karaokePlayer.ImageComponent.Options.Opacity = 0.01
	k.world.AddEntity(k.karaokePlayer)

	k.inputEnt = entity.CreateMenuInput()
	k.world.AddEntity(k.inputEnt)
}

func loadImage(encoded string) image.Image {
	decoded, _ := base64.StdEncoding.DecodeString(encoded)
	raw := bytes.NewBuffer(decoded)
	img, _, err := image.Decode(raw)
	if err != nil {
		panic(err)
	}

	return img
}

func (k *KaraokeScene) loadBackground() {
	if k.nextImage == nil {
		k.nextImage = ebiten.NewImageFromImage(loadImage(k.Session.Backgrounds[k.Session.backgroundIdx].Image))
	}

	if k.currentImage == nil {
		k.backgroundBack.Image = k.nextImage
		k.backgroundFront.Image = k.nextImage
		k.backgroundFrontFadeIn = 0
	} else {
		k.backgroundBack.Image = k.currentImage
		k.backgroundFront.Image = k.nextImage
		k.backgroundFront.Options.Opacity = 0.0001
		k.backgroundFrontFadeIn = time.Duration(k.Session.Backgrounds[k.Session.backgroundIdx].FadeIn)
	}

	k.currentImage = k.nextImage

	if k.Session.backgroundIdx+1 < len(k.Session.Backgrounds) {
		k.nextImage = ebiten.NewImageFromImage(loadImage(k.Session.Backgrounds[k.Session.backgroundIdx+1].Image))
	}
}

func parseHex(hex string) color.Color {
	result, _ := colorx.ParseHexColor(hex)
	return result
}

func (k *KaraokeScene) Start(game *Game) {
	if k.Session == nil || len(k.Session.Backgrounds) <= 0 || len(k.Session.Inputs) <= 0 {
		panic("Must set Session, at least one background must be set and one input")
	}

	k.world = &ecs.World{}
	k.rand = rand.New(rand.NewSource(time.Now().Unix()))
	k.state = KaraokeStateStarting

	k.soundInfo = map[components.KaraokeSound]*karaokeInfo{
		components.KaraokeSoundA: {
			input: components.InputKindKaraokeA,
			sound: entity.CreateKaraokeInputSound(),
		},
		components.KaraokeSoundB: {
			input: components.InputKindKaraokeB,
			sound: entity.CreateKaraokeInputSound(),
		},
		components.KaraokeSoundX: {
			input: components.InputKindKaraokeX,
			sound: entity.CreateKaraokeInputSound(),
		},
		components.KaraokeSoundY: {
			input: components.InputKindKaraokeY,
			sound: entity.CreateKaraokeInputSound(),
		},
	}

	for key, _ := range k.soundInfo {
		raw, _ := base64.StdEncoding.DecodeString(k.Session.Sounds[string(key)])
		k.soundInfo[key].sound.Sound = &components.Sound{
			Source:     []byte(raw),
			SampleRate: 44100,
			Volume:     1,
			SoundType:  assets.SoundTypeWav,
		}
	}

	k.Session.backgroundIdx = 0
	k.inputLeeway = 100 * time.Millisecond
	k.backgroundElapsed = 0

	k.addSystems(game.audioCtx)
	k.addEnts()

	// Make it do it dynamically based on binding
	switch k.inputEnt.InputMode {
	case components.InputModeGamepad:
		img, _ := assets.LoadEbitenImage(assets.ImageIconXboxSeriesXA)
		k.soundInfo[components.KaraokeSoundA].image = img

		img, _ = assets.LoadEbitenImage(assets.ImageIconXboxSeriesXB)
		k.soundInfo[components.KaraokeSoundB].image = img

		img, _ = assets.LoadEbitenImage(assets.ImageIconXboxSeriesXX)
		k.soundInfo[components.KaraokeSoundX].image = img

		img, _ = assets.LoadEbitenImage(assets.ImageIconXboxSeriesXY)
		k.soundInfo[components.KaraokeSoundY].image = img
	case components.InputModeKeyboard:

		img, _ := assets.LoadEbitenImage(assets.ImageIconKeyboardKeyW)
		k.soundInfo[components.KaraokeSoundA].image = img

		img, _ = assets.LoadEbitenImage(assets.ImageIconKeyboardKeyD)
		k.soundInfo[components.KaraokeSoundB].image = img

		img, _ = assets.LoadEbitenImage(assets.ImageIconKeyboardKeyA)
		k.soundInfo[components.KaraokeSoundX].image = img

		img, _ = assets.LoadEbitenImage(assets.ImageIconKeyboardKeyS)
		k.soundInfo[components.KaraokeSoundY].image = img
	}

	for _, info := range k.soundInfo {
		k.world.AddEntity(info.sound)
	}

	k.loadBackground()

	k.scorePlayer = entity.CreateSoundPlayer(assets.SoundBit8CoinOne)
	k.world.AddEntity(k.scorePlayer)

	k.scoreOpacity = 0

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	k.scoreTitleFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    160,
		DPI:     72 * 2,
		Hinting: font.HintingFull,
	})

	k.scoreFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    80,
		DPI:     72 * 2,
		Hinting: font.HintingFull,
	})

	k.scoreFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    80,
		DPI:     72 * 2,
		Hinting: font.HintingFull,
	})

	k.comboFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    80,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	k.scoreColors = []color.Color{
		parseHex("#ff0000"), parseHex("#ff8000"),
		parseHex("#ffff00"), parseHex("#80ff00"),
		parseHex("#00ff00"), parseHex("#00ff80"),
		parseHex("#00ffff"), parseHex("#0080ff"),
		parseHex("#0000ff"), parseHex("#8000ff"),
		parseHex("#ff00ff"), parseHex("#ff0080"),
	}

	k.activeScoreColor = k.scoreColors[0]
	k.nextColorUpdate = 0

	k.textScreen = ebiten.NewImage(windowWidth, windowHeight)
}

func (k *KaraokeScene) End(*Game) {
	k.Session = nil

	k.world = nil
	k.rand = nil
	k.timeElapsed = 0
	k.inputEnt = nil
	if k.musicEnt.Player != nil {
		k.musicEnt.Player.Pause()
	}
	k.musicEnt = nil
	k.karaokePlayer = nil

	if k.scorePlayer != nil && k.scorePlayer.Player != nil {
		k.scorePlayer.Player.Pause()
		k.scorePlayer.Player.Close()
	}

	if k.textScreen != nil {
		k.textScreen.Dispose()
	}
}

func (k *KaraokeScene) Update(dt time.Duration, _ *Game) {
	if k.inputEnt.InputPressedDuration(components.InputKindFastGameSpeed) > 0 {
		dt *= 5
	}

	k.world.Update(float32(dt) / float32(time.Second))
	k.timeElapsed += dt

	dtSecond := float64(dt) / float64(time.Second)

	switch k.state {
	case KaraokeStateStarting:
		const fadeInTime = 0 * time.Second
		const fadeInTimeBackground = 0 * time.Second

		k.karaokePlayer.ImageComponent.Options.Opacity = utility.ClampFloat64(1-(float64(fadeInTime-k.timeElapsed)/float64(fadeInTime)), 0, 1)
		k.backgroundFront.ImageComponent.Options.Opacity = utility.ClampFloat64(1-(float64(fadeInTimeBackground-k.timeElapsed)/float64(fadeInTimeBackground)), 0, 1)
		if k.timeElapsed > fadeInTimeBackground {
			k.timeElapsed = 0
			k.state = KaraokeStateSinging
			k.musicEnt.SoundComponent.Active = true
		}
	case KaraokeStateSinging:
		if k.backgroundFront.Options.Opacity > 0 && k.backgroundFront.Options.Opacity < 1 {
			step := (1 / (float64(k.backgroundFrontFadeIn) / float64(time.Second))) * (float64(dt) / float64(time.Second))
			k.backgroundFront.Options.Opacity = k.backgroundFront.Options.Opacity + step
			if k.backgroundFront.Options.Opacity > 1 {
				k.backgroundFront.Options.Opacity = 1
			}
		}

		targetBackground := k.Session.Backgrounds[k.Session.backgroundIdx]
		if k.Session.backgroundIdx+1 < len(k.Session.Backgrounds) &&
			k.timeElapsed > k.backgroundElapsed+time.Duration(targetBackground.Duration) {
			k.backgroundElapsed = k.timeElapsed
			k.Session.backgroundIdx++
			k.loadBackground()
		}

		complete := false
		{
			var selectedInput *KaraokeInput
			bestScore := KaraokeScore(0)

			for i, input := range k.Session.Inputs {

				if k.timeElapsed < time.Duration(input.StartTime) || input.xPostion-100 > windowWidth {
					if input.xPostion-100 > windowWidth && i == len(k.Session.Inputs)-1 {
						complete = true
					}
					continue
				}

				if input.xSpeed == 0 {
					input.xSpeed = 1650 / (float64(input.Duration) / float64(time.Second))
				}

				input.xPostion += input.xSpeed * (float64(dt) / float64(time.Second))

				inputKind := k.soundInfo[input.Sound].input
				x := windowWidth - input.xPostion + 50
				if input.hitPostion <= 0 && x > karaLeftBound && x < karaRightBound && k.inputEnt.InputJustPressed(inputKind) {
					if score := Score(input.xPostion); score > bestScore {
						selectedInput = input
						bestScore = score
					}
				}
			}

			if selectedInput != nil {
				selectedInput.hitPostion = selectedInput.xPostion
				k.soundInfo[selectedInput.Sound].sound.Active = true
				k.soundInfo[selectedInput.Sound].sound.Restart = true

				textEnt := entity.CreateFloatingText()
				textEnt.DestoryBoundComponent.Max = math.Vector2{
					X: windowWidth,
					Y: windowHeight,
				}
				textEnt.Color = parseHex("#A020F0")
				textEnt.ConstantSpeedComponent.Speed.Y = -1000
				textEnt.TextComponent.Text = Score(selectedInput.xPostion).String()
				textEnt.TextComponent.Font = k.comboFont
				b := text.BoundString(textEnt.Font, textEnt.Text)
				textEnt.Postion.X = windowWidth - selectedInput.xPostion - float64(b.Dx()/2)
				textEnt.Postion.Y = selectedInput.Y()
				textEnt.Layer = ImageLayerKaraokeText
				defer k.world.AddEntity(textEnt)
			}
		}

		if complete {
			k.state = KaraokeStateComplete
			k.timeElapsed = 0
			k.scorePlayer.Active = true
			k.scorePlayer.Loop = true
		}
	case KaraokeStateComplete:
		const uiFadeOutTime = 1 * time.Second
		const scoreFadeInTime = 1*time.Second + 250*time.Millisecond
		const centerPlayerTime = float64(200*time.Millisecond) / float64(time.Second)

		k.ui.Options.Opacity = utility.ClampFloat64(float64(uiFadeOutTime-k.timeElapsed)/float64(uiFadeOutTime), 0, 1)
		if k.ui.Options.Opacity <= 0 {
			defer k.world.RemoveEntity(k.ui.BasicEntity)
		}

		maxX := (windowWidth/2 - k.karaokePlayer.TransformComponent.Size.X/2)
		k.karaokePlayer.Postion.X = utility.ClampFloat64(
			k.karaokePlayer.Postion.X+(maxX/float64(centerPlayerTime)*dtSecond),
			0,
			maxX,
		)

		scoreOpacityStep := (1 / (float64(scoreFadeInTime) / float64(time.Second))) * (float64(dt) / float64(time.Second))
		k.scoreOpacity = utility.ClampFloat64(
			k.scoreOpacity+scoreOpacityStep,
			0, 1,
		)

		if k.timeElapsed > karaScoreSpinTime {
			if k.activeScoreColor != color.White {
				k.scorePlayer.SoundComponent.ChangeSound(components.LoadSound(assets.SoundBit8CoinOneRepeated))
				k.scorePlayer.Restart = true
				k.scorePlayer.Loop = false
			}
			k.activeScoreColor = color.White
		} else if k.timeElapsed > k.nextColorUpdate {
			k.nextColorUpdate = k.timeElapsed + 250*time.Millisecond
			k.activeScoreColor = k.scoreColors[k.rand.Intn(len(k.scoreColors))]
		}

	}
}

func (k *KaraokeScene) CalcScore() (score KaraokeScore) {
	score = 1

	for _, input := range k.Session.Inputs {
		score += Score(input.hitPostion)
	}

	return
}

func (k *KaraokeScene) Draw(screen *ebiten.Image) {
	queue := RenderCmds{}
	for _, system := range k.world.Systems() {
		if render, ok := system.(RenderingSystem); ok {
			render.Render(&queue)
		}
	}

	queue.Sort()
	for _, item := range queue {
		item.Draw(screen)
	}

	switch k.state {
	case KaraokeStateSinging:
		op := ebiten.DrawImageOptions{}
		for _, input := range k.Session.Inputs {

			if input.xPostion <= 0 || input.xPostion-100 > windowWidth {
				continue
			}

			op.ColorM.Reset()
			op.GeoM.Reset()

			op.GeoM.Translate(windowWidth-input.xPostion, input.Y())

			if input.hitPostion > 0 {
				op.ColorM.Scale(-1, -1, -1, 1)
				op.ColorM.Translate(1, 1, 1, 0)
			}

			screen.DrawImage(k.soundInfo[input.Sound].image, &op)
		}
	case KaraokeStateComplete:
		yStart := 300
		k.textScreen.Fill(color.Transparent)

		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(1, 1, 1, k.scoreOpacity)

		{
			const scoreTitleText = "SCORE"
			b := text.BoundString(k.scoreTitleFont, scoreTitleText)
			x := windowWidth/2 - b.Dx()/2
			text.Draw(k.textScreen, scoreTitleText, k.scoreTitleFont, x, yStart, color.Black)
			yStart += b.Dy() + 30
		}
		{
			var score int
			if k.timeElapsed < karaScoreSpinTime {
				score = utility.RandRange(k.rand, int(k.CalcScore())/6, int(k.CalcScore())*6)
			} else {
				score = int(k.CalcScore())
			}

			scoreText := strconv.Itoa(score)
			b := text.BoundString(k.scoreFont, scoreText)
			x := windowWidth/2 - b.Dx()/2
			text.Draw(k.textScreen, scoreText, k.scoreFont, x, yStart, k.activeScoreColor)
		}

		screen.DrawImage(k.textScreen, op)
	}

}
