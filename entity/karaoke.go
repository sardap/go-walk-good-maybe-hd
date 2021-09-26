package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type KaraokePlayer struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
}

func CreateKaraokePlayer() *KaraokePlayer {
	img, _ := assets.LoadEbitenImage(assets.ImageWhaleSinging)

	return &KaraokePlayer{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(img.Bounds().Dx()),
				Y: float64(img.Bounds().Dy()),
			},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
	}
}

type KaraokeInputSound struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.SoundComponent
}

func CreateKaraokeInputSound(soundAsset interface{}) *KaraokeInputSound {
	return &KaraokeInputSound{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		SoundComponent: &components.SoundComponent{
			Sound: components.LoadSound(soundAsset),
		},
	}
}

type KaraokeInputIcon struct {
	ecs.BasicEntity
	*components.KaraokeImageComponent
}

func CreateKaraokeInputIcon(soundType components.KaraokeSound, startTime time.Duration) *KaraokeInputIcon {
	return &KaraokeInputIcon{
		BasicEntity: ecs.NewBasic(),
		KaraokeImageComponent: &components.KaraokeImageComponent{
			SoundType: soundType,
			Start:     startTime,
		},
	}
}

func CreateKaraokeInputIconA(startTime time.Duration) *KaraokeInputIcon {
	return CreateKaraokeInputIcon(components.KaraokeSoundA, startTime)
}

func CreateKaraokeInputIconB(startTime time.Duration) *KaraokeInputIcon {
	return CreateKaraokeInputIcon(components.KaraokeSoundB, startTime)
}

func CreateKaraokeInputIconX(startTime time.Duration) *KaraokeInputIcon {
	return CreateKaraokeInputIcon(components.KaraokeSoundX, startTime)
}

func CreateKaraokeInputIconY(startTime time.Duration) *KaraokeInputIcon {
	return CreateKaraokeInputIcon(components.KaraokeSoundY, startTime)
}
