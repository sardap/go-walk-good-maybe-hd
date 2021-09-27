package common

import (
	"time"

	"github.com/sardap/walk-good-maybe-hd/components"
)

type KaraokeInput struct {
	StartTime  time.Duration           `toml:"start_time"`
	Duration   time.Duration           `toml:"duration"`
	Sound      components.KaraokeSound `toml:"sound"`
	XPostion   float64                 `toml:"x_postion,omitempty"`
	XSpeed     float64                 `toml:"x_speed,omitempty"`
	HitPostion float64                 `toml:"x_hit_postion,omitempty"`
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

type KaraokeBackground struct {
	Duration time.Duration `json:"duration"`
	FadeIn   time.Duration `json:"fade_in"`
	Image    string        `json:"image"`
}

type KaraokeSession struct {
	Inputs        []*KaraokeInput      `json:"inputs"`
	Backgrounds   []*KaraokeBackground `json:"backgrounds"`
	Sounds        map[string]string    `json:"sounds"`
	Music         string               `json:"music"`
	SampleRate    int                  `json:"sampleRate"`
	BackgroundIdx int
}
