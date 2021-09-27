package common

import (
	"time"

	"github.com/sardap/walk-good-maybe-hd/components"
)

type KaraokeInput struct {
	StartTime  time.Duration
	Duration   time.Duration
	Sound      components.KaraokeSound
	XPostion   float64
	XSpeed     float64
	HitPostion float64
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

type KaraokeBackground struct {
	Duration time.Duration
	FadeIn   time.Duration
	Image    string
}

type KaraokeSession struct {
	Inputs           []*KaraokeInput
	TitleScreenImage string
	Backgrounds      []*KaraokeBackground
	Sounds           map[string]string
	TextImages       map[string]string
	Music            string
	SampleRate       int
	BackgroundIdx    int
}
