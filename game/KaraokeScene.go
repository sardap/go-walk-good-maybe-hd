package game

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type DurationMil time.Duration

type KaraokeInput struct {
	Duration DurationMil             `json:"duration"`
	Sound    components.KaraokeSound `json:"sound"`
}

func (d *DurationMil) UnmarshalJSON(b []byte) error {
	var v time.Duration
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	*d = DurationMil(v * time.Millisecond)

	return nil
}

type KaraokeSession struct {
	Inputs []KaraokeInput `json:"inputs"`
}

type KaraokeState int

const (
	KaraokeStateStarting KaraokeState = iota
	KaraokeStateSinging
)

type KaraokeScene struct {
	Session  *KaraokeSession
	eventIdx int

	rand  *rand.Rand
	world *ecs.World

	inputEnt      *entity.InputEnt
	musicEnt      *entity.SoundPlayer
	karaokePlayer *entity.KaraokePlayer
	background    *entity.BasicImage

	timeElapsed time.Duration
	baseElapsed time.Duration
	state       KaraokeState

	karaokeSounds map[components.KaraokeSound]*entity.KaraokeInputSound
	inputImages   map[components.KaraokeSound]*ebiten.Image
}

func (k *KaraokeScene) addSystems(audioCtx *audio.Context) {
	var soundable *Soundable
	k.world.AddSystemInterface(CreateSoundSystem(audioCtx), soundable, nil)

	var renderable *ImageRenderable
	k.world.AddSystemInterface(CreateImageRenderSystem(), renderable, nil)

	var inputable *Inputable
	k.world.AddSystemInterface(CreateInputSystem(), inputable, nil)
}

func (k *KaraokeScene) addEnts() {
	k.musicEnt = entity.CreateSoundPlayer(assets.MusicPdRockBackground)
	k.musicEnt.SoundComponent.Active = false
	k.world.AddEntity(k.musicEnt)

	k.background = entity.CreateBasicImage(assets.ImageKaraokeBackground)
	k.background.ImageComponent.Layer = ImageLayerBottom
	k.world.AddEntity(k.background)

	k.karaokePlayer = entity.CreateKaraokePlayer()
	k.karaokePlayer.ImageComponent.Layer = ImageLayerObjects
	k.karaokePlayer.Postion.Y = windowHeight / 2
	k.karaokePlayer.ImageComponent.Options.Opacity = 0.01
	k.world.AddEntity(k.karaokePlayer)
}

func (k *KaraokeScene) Start(game *Game) {
	k.world = &ecs.World{}
	k.rand = rand.New(rand.NewSource(time.Now().Unix()))
	k.state = KaraokeStateStarting

	k.karaokeSounds = make(map[components.KaraokeSound]*entity.KaraokeInputSound)

	k.karaokeSounds[components.KaraokeSoundA] = entity.CreateKaraokeInputSound(assets.SoundSpaceVoiceC5)
	k.world.AddEntity(k.karaokeSounds[components.KaraokeSoundA])

	k.karaokeSounds[components.KaraokeSoundB] = entity.CreateKaraokeInputSound(assets.SoundSpaceVoiceC5)
	k.world.AddEntity(k.karaokeSounds[components.KaraokeSoundB])

	k.karaokeSounds[components.KaraokeSoundX] = entity.CreateKaraokeInputSound(assets.SoundSpaceVoiceC5)
	k.world.AddEntity(k.karaokeSounds[components.KaraokeSoundX])

	k.karaokeSounds[components.KaraokeSoundY] = entity.CreateKaraokeInputSound(assets.SoundSpaceVoiceC5)
	k.world.AddEntity(k.karaokeSounds[components.KaraokeSoundY])

	k.inputImages = make(map[components.KaraokeSound]*ebiten.Image)

	img, _ := assets.LoadEbitenImage(assets.ImageIconXboxSeriesXA)
	k.inputImages[components.KaraokeSoundA] = img

	img, _ = assets.LoadEbitenImage(assets.ImageIconXboxSeriesXB)
	k.inputImages[components.KaraokeSoundB] = img

	img, _ = assets.LoadEbitenImage(assets.ImageIconXboxSeriesXX)
	k.inputImages[components.KaraokeSoundX] = img

	img, _ = assets.LoadEbitenImage(assets.ImageIconXboxSeriesXY)
	k.inputImages[components.KaraokeSoundY] = img

	k.eventIdx = 0

	k.addSystems(game.audioCtx)
	k.addEnts()
}

func (k *KaraokeScene) End(*Game) {
	k.world = nil
	k.rand = nil
	k.timeElapsed = 0
	k.inputEnt = nil
	if k.musicEnt.Player != nil {
		k.musicEnt.Player.Pause()
	}
	k.musicEnt = nil
	k.karaokePlayer = nil
}

func (k *KaraokeScene) Update(dt time.Duration, _ *Game) {
	k.world.Update(float32(dt) / float32(time.Second))
	k.timeElapsed += dt

	const fadeInTime = 0 * time.Second
	const fadeInTimeBackground = 0 * time.Second
	switch k.state {
	case KaraokeStateStarting:
		k.karaokePlayer.ImageComponent.Options.Opacity = utility.ClampFloat64(1-(float64(fadeInTime-k.timeElapsed)/float64(fadeInTime)), 0, 1)
		k.background.ImageComponent.Options.Opacity = utility.ClampFloat64(1-(float64(fadeInTimeBackground-k.timeElapsed)/float64(fadeInTimeBackground)), 0, 1)
		if k.timeElapsed > fadeInTimeBackground {
			k.timeElapsed = 0
			k.state = KaraokeStateSinging
			k.musicEnt.SoundComponent.Active = true
		}
	case KaraokeStateSinging:
		if k.eventIdx > len(k.Session.Inputs)-1 {
			break
		}

		if k.timeElapsed > k.baseElapsed+time.Duration(k.Session.Inputs[k.eventIdx].Duration) {
			k.eventIdx++
			k.baseElapsed = k.timeElapsed
		}
	}
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

	if k.state != KaraokeStateSinging {
		return
	}

	op := ebiten.DrawImageOptions{}
	timeBase := time.Duration(0)
	for i := k.eventIdx; i < k.eventIdx+1; i++ {
		op.ColorM.Reset()
		op.GeoM.Reset()

		event := &k.Session.Inputs[i]
		timeBase += time.Duration(event.Duration)

		percent := float64(k.timeElapsed-k.baseElapsed) / float64(timeBase)

		x := float64(windowWidth - 1650*percent)

		op.GeoM.Translate(x, 400)
		screen.DrawImage(k.inputImages[event.Sound], &op)
	}

}
