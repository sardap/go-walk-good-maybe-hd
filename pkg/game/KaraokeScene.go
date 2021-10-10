package game

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	gomath "math"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sardap/walk-good-maybe-hd/pkg/assets"
	"github.com/sardap/walk-good-maybe-hd/pkg/common"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
	"github.com/sardap/walk-good-maybe-hd/pkg/entity"
	"github.com/sardap/walk-good-maybe-hd/pkg/math"
	"github.com/sardap/walk-good-maybe-hd/pkg/utility"
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

const (
	karaBoundStep     = 40
	karaScoreSpinTime = 2*time.Second + 500*time.Millisecond
	uiY               = 200
	karaRightBound    = karaLeftBound + karaBoundStep*3
	karaLeftBound     = 260
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
	KaraokeScoreMiss    KaraokeScore = 0
	KaraokeScoreOkay    KaraokeScore = 10
	KaraokeScoreGood    KaraokeScore = 20
	KaraokeScoreGreat   KaraokeScore = 40
	KaraokeScorePerfect KaraokeScore = 80
)

func karaokeScore(hitTime, targetTime time.Duration) KaraokeScore {
	delta := time.Duration(gomath.Abs(float64(hitTime - targetTime)))

	switch {
	case delta < 250*time.Millisecond:
		return KaraokeScorePerfect
	case delta < 500*time.Millisecond:
		return KaraokeScoreGreat
	case delta < 1*time.Second:
		return KaraokeScoreGood
	}

	return KaraokeScoreMiss
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
	Session     *common.KaraokeSession
	inputLeeway time.Duration

	rand  *rand.Rand
	world *ecs.World

	inputEnt      *entity.InputEnt
	musicPlayer   *audio.Player
	karaokePlayer *entity.KaraokePlayer

	inputImg              *ebiten.Image
	currentImage          *ebiten.Image
	titleScreenImage      *ebiten.Image
	nextImage             *ebiten.Image
	backgroundFront       *entity.BasicImage
	backgroundFrontFadeIn time.Duration
	backgroundBack        *entity.BasicImage
	ui                    *entity.BasicImage
	targetMarker          *entity.BasicImage
	loadingBackgroundLock *sync.Mutex

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

	soundInfo  map[components.KaraokeSound]*karaokeInfo
	textImages map[KaraokeScore]*ebiten.Image
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
				Scale:   math.Vector2{X: 1, Y: 1},
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
			Options: components.DrawOptions{
				Scale: math.Vector2{X: 1, Y: 1},
			},
		},
	}
	k.world.AddEntity(k.backgroundBack)

	{
		k.ui = entity.CreateBasicImageEmpty()
		k.ui.Postion.Y = uiY
		k.ui.Layer = ImageLayerKaraokeUi

		img := ebiten.NewImage(windowWidth, 336)
		clr := parseHex("#545454")
		clr.A = 255 * 0.6
		img.Fill(clr)
		k.ui.Image = img
		k.ui.TransformComponent.Size.X = float64(img.Bounds().Dx())
		k.ui.TransformComponent.Size.Y = float64(img.Bounds().Dy())

		subRect := img.SubImage(image.Rect(0, 0, windowWidth, 5)).(*ebiten.Image)
		subRect.Fill(ClrBlack)

		subRect = img.SubImage(image.Rect(0, img.Bounds().Dy()-5, windowWidth, img.Bounds().Dy())).(*ebiten.Image)
		subRect.Fill(ClrBlack)

		k.world.AddEntity(k.ui)
	}

	img := ebiten.NewImage(karaRightBound-karaLeftBound, int(k.ui.TransformComponent.Size.Y-8))
	lightRed := parseHex("#f32f42")
	lightRed.A = 128
	img.Fill(lightRed)
	subRect := img.SubImage(image.Rect((karaBoundStep*3)/2-5, 0, (karaBoundStep*3)/2+10, img.Bounds().Dy())).(*ebiten.Image)
	aqua := parseHex("#2ff3e0")
	aqua.A = 128
	subRect.Fill(aqua)
	k.targetMarker = &entity.BasicImage{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Postion: math.Vector2{
				Y: k.ui.Postion.Y + 4,
				X: karaLeftBound,
			},
			Size: math.Vector2{
				X: float64(img.Bounds().Dx()),
				Y: float64(img.Bounds().Dy()),
			},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Layer:  ImageLayerKaraokeText,
			Options: components.DrawOptions{
				Scale: math.Vector2{X: 1, Y: 1},
			},
			Image: img,
		},
	}
	k.world.AddEntity(k.targetMarker)

	k.karaokePlayer = entity.CreateKaraokePlayer()
	k.karaokePlayer.ImageComponent.Layer = ImageLayerKaraokeObjects
	k.karaokePlayer.Postion.Y = windowHeight/2 + (k.karaokePlayer.TransformComponent.Size.Y * 0.3)
	k.karaokePlayer.ImageComponent.Options.Opacity = 1
	k.world.AddEntity(k.karaokePlayer)

	k.inputEnt = entity.CreateMenuInput()
	k.world.AddEntity(k.inputEnt)
}

func KaraokeLoadImage(data []byte) image.Image {

	raw := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	base64.StdEncoding.Decode(raw, []byte(data))
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}

	return img
}

func (k *KaraokeScene) loadBackground() {
	k.loadingBackgroundLock.Lock()
	defer k.loadingBackgroundLock.Unlock()

	if k.nextImage == nil {
		k.nextImage = ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.Backgrounds[k.Session.BackgroundIdx].Image)))
	}

	if k.currentImage == nil {
		k.backgroundBack.Image = k.nextImage
		k.backgroundFront.Image = k.nextImage
		k.backgroundFrontFadeIn = 0
	} else {
		k.backgroundBack.Image = k.currentImage
		k.backgroundFront.Image = k.nextImage
		k.backgroundFront.Options.Opacity = 0.0001
		k.backgroundFrontFadeIn = time.Duration(k.Session.Backgrounds[k.Session.BackgroundIdx].FadeIn)
	}

	k.currentImage = k.nextImage

	if k.Session.BackgroundIdx+1 < len(k.Session.Backgrounds) {
		go func() {
			k.loadingBackgroundLock.Lock()
			defer k.loadingBackgroundLock.Unlock()
			k.nextImage = ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.Backgrounds[k.Session.BackgroundIdx+1].Image)))
		}()
	}
}

func (k *KaraokeScene) Start(game *Game) {
	if k.Session == nil || len(k.Session.Backgrounds) <= 0 || len(k.Session.Inputs) <= 0 {
		panic("Must set Session, at least one background must be set and one input")
	}

	k.loadingBackgroundLock = &sync.Mutex{}

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

	k.textImages = map[KaraokeScore]*ebiten.Image{
		KaraokeScoreOkay:    ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.TextImages["okay"]))),
		KaraokeScoreGood:    ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.TextImages["good"]))),
		KaraokeScoreGreat:   ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.TextImages["great"]))),
		KaraokeScorePerfect: ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.TextImages["perfect"]))),
	}

	for key := range k.soundInfo {
		raw, _ := base64.StdEncoding.DecodeString(k.Session.Sounds[string(key)])
		k.soundInfo[key].sound.Sound = &components.Sound{
			Source:     []byte(raw),
			SampleRate: 44100,
			Volume:     1,
			SoundType:  assets.SoundTypeWav,
		}
	}

	k.Session.BackgroundIdx = 0
	k.inputLeeway = 100 * time.Millisecond
	k.backgroundElapsed = 0

	k.addSystems(game.audioCtx)
	k.addEnts()

	k.titleScreenImage = ebiten.NewImageFromImage(KaraokeLoadImage([]byte(k.Session.TitleScreenImage)))

	{
		buf := bytes.NewBuffer([]byte(k.Session.Music))
		base64.StdEncoding.Decode(buf.Bytes(), []byte(k.Session.Music))
		stream, err := mp3.DecodeWithSampleRate(48000, bytes.NewReader(buf.Bytes()))
		if err != nil {
			panic(err)
		}
		k.musicPlayer, _ = audio.NewPlayer(game.audioCtx, stream)

		img := ebiten.NewImage(windowWidth, windowHeight)

		k.soundInfo[components.KaraokeSoundA].image = loadIconImage(k.inputEnt.InputComponent, components.InputKindKaraokeA)
		k.soundInfo[components.KaraokeSoundB].image = loadIconImage(k.inputEnt.InputComponent, components.InputKindKaraokeB)
		k.soundInfo[components.KaraokeSoundX].image = loadIconImage(k.inputEnt.InputComponent, components.InputKindKaraokeX)
		k.soundInfo[components.KaraokeSoundY].image = loadIconImage(k.inputEnt.InputComponent, components.InputKindKaraokeY)

		k.inputImg = img
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
	k.karaokePlayer = nil

	k.musicPlayer.Pause()

	if k.scorePlayer != nil && k.scorePlayer.Player != nil {
		k.scorePlayer.Player.Pause()
		k.scorePlayer.Player.Close()
	}

	if k.currentImage != nil {
		k.currentImage.Dispose()
	}

	if k.titleScreenImage != nil {
		k.titleScreenImage.Dispose()
	}

	if k.nextImage != nil {
		k.nextImage.Dispose()
	}

	if k.textScreen != nil {
		k.textScreen.Dispose()
	}

	if k.inputImg != nil {
		k.inputImg.Dispose()
	}

	for _, image := range k.textImages {
		image.Dispose()
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
		const fadeInTimeBackground = 3 * time.Second
		if k.timeElapsed > fadeInTimeBackground {
			k.timeElapsed = 0
			k.state = KaraokeStateSinging
			k.musicPlayer.Play()
		}
	case KaraokeStateSinging:
		k.timeElapsed = k.musicPlayer.Current()

		if k.backgroundFront.Options.Opacity > 0 && k.backgroundFront.Options.Opacity < 1 {
			step := (1 / (float64(k.backgroundFrontFadeIn) / float64(time.Second))) * (float64(dt) / float64(time.Second))
			k.backgroundFront.Options.Opacity = k.backgroundFront.Options.Opacity + step
			if k.backgroundFront.Options.Opacity > 1 {
				k.backgroundFront.Options.Opacity = 1
			}
		}

		targetBackground := k.Session.Backgrounds[k.Session.BackgroundIdx]
		if k.Session.BackgroundIdx+1 < len(k.Session.Backgrounds) &&
			k.timeElapsed > k.backgroundElapsed+time.Duration(targetBackground.Duration) {
			k.backgroundElapsed = k.timeElapsed
			k.Session.BackgroundIdx++
			k.loadBackground()
		}

		complete := false
		{
			var selectedInput *common.KaraokeInput
			bestScore := KaraokeScore(0)

			if !k.musicPlayer.IsPlaying() {
				complete = true
			}

			for _, input := range k.Session.Inputs {
				inputKind := k.soundInfo[components.KaraokeSound(input.Sound)].input
				if input.HitTime == 0 && k.inputEnt.InputJustPressed(inputKind) {
					if score := karaokeScore(k.timeElapsed, input.TargetHitTime); score != KaraokeScoreMiss && score > bestScore {
						selectedInput = input
						bestScore = score
					}
				}

				if score := karaokeScore(k.timeElapsed, input.TargetHitTime); score == KaraokeScoreGreat {
					fmt.Printf("fuck")
				}
			}

			if selectedInput != nil {
				selectedInput.HitTime = k.timeElapsed
				k.soundInfo[components.KaraokeSound(selectedInput.Sound)].sound.Active = true
				k.soundInfo[components.KaraokeSound(selectedInput.Sound)].sound.Restart = true

				textEnt := entity.CreateFloatingTimedImage()
				textEnt.DestoryBoundComponent.Min = math.Vector2{
					X: -500,
					Y: -500,
				}
				textEnt.DestoryBoundComponent.Max = math.Vector2{
					X: windowWidth + 500,
					Y: windowHeight + 500,
				}
				textEnt.ConstantSpeedComponent.Speed.X = 800
				textEnt.Image = k.textImages[bestScore]
				textEnt.Postion.X = k.karaokePlayer.Postion.X + k.karaokePlayer.TransformComponent.Size.X
				textEnt.Postion.Y = k.karaokePlayer.Postion.Y + (k.karaokePlayer.TransformComponent.Size.Y * 0.3)
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
		const centerPlayerTime = float64(500*time.Millisecond) / float64(time.Second)

		k.ui.Options.Opacity = utility.ClampFloat64(float64(uiFadeOutTime-k.timeElapsed)/float64(uiFadeOutTime), 0, 1)
		k.targetMarker.Options.Opacity = utility.ClampFloat64(float64(uiFadeOutTime-k.timeElapsed)/float64(uiFadeOutTime), 0, 1)
		if k.ui.Options.Opacity <= 0 {
			defer k.world.RemoveEntity(k.ui.BasicEntity)
			defer k.world.RemoveEntity(k.targetMarker.BasicEntity)
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

func karaokeInputY(k *common.KaraokeInput) float64 {
	const buffer = 12
	const size = 80
	switch k.Sound {
	case components.KaraokeSoundA:
		return uiY + (size * 3) + buffer
	case components.KaraokeSoundB:
		return uiY + (size * 2) + buffer
	case components.KaraokeSoundX:
		return uiY + (size * 1) + buffer
	case components.KaraokeSoundY:
		return uiY + (size * 0) + buffer
	}

	return 400
}

func (k *KaraokeScene) calcScore() (score KaraokeScore) {
	score = 1

	for _, input := range k.Session.Inputs {
		score += karaokeScore(input.HitTime, input.TargetHitTime)
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
	case KaraokeStateStarting:
		screen.DrawImage(k.titleScreenImage, &ebiten.DrawImageOptions{})
	case KaraokeStateSinging:
		k.inputImg.Fill(color.Transparent)

		for _, input := range k.Session.Inputs {
			x := (input.TargetHitTime-k.musicPlayer.Current())/time.Millisecond + 40
			if x > windowWidth {
				continue
			}
			op := ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), karaokeInputY(input))

			if input.HitTime > 0 {
				op.ColorM.Scale(-1, -1, -1, 1)
				op.ColorM.Translate(1, 1, 1, 0)
			}

			k.inputImg.DrawImage(k.soundInfo[components.KaraokeSound(input.Sound)].image, &op)
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(0, 0)
		screen.DrawImage(k.inputImg, op)
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
				score = utility.RandRange(k.rand, int(k.calcScore())/6, int(k.calcScore())*6)
			} else {
				score = int(k.calcScore())
			}

			scoreText := strconv.Itoa(score)
			b := text.BoundString(k.scoreFont, scoreText)
			x := windowWidth/2 - b.Dx()/2
			text.Draw(k.textScreen, scoreText, k.scoreFont, x, yStart, k.activeScoreColor)
		}

		screen.DrawImage(k.textScreen, op)
	}

}
