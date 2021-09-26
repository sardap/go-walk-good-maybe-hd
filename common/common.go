// package common

// import (
// 	"time"

// 	"github.com/sardap/walk-good-maybe-hd/components"
// )

// type KaraokeInput struct {
// 	StartTime  time.Time               `json:"start_time"`
// 	Duration   time.Time               `json:"duration"`
// 	Sound      components.KaraokeSound `json:"sound"`
// 	xPostion   float64
// 	xSpeed     float64
// 	hitPostion float64
// }

// func (k *KaraokeInput) Y() float64 {
// 	switch k.Sound {
// 	case components.KaraokeSoundA:
// 		return 550
// 	case components.KaraokeSoundB:
// 		return 450
// 	case components.KaraokeSoundX:
// 		return 350
// 	case components.KaraokeSoundY:
// 		return 250
// 	}

// 	return 400
// }

// const (
// 	karaBoundStep     = 100
// 	karaCenter        = 850
// 	karaLeftBound     = karaCenter - karaBoundStep
// 	karaRightBound    = karaCenter + karaBoundStep
// 	karaScoreSpinTime = 2*time.Second + 500*time.Millisecond
// )

// type KaraokeScore int

// func (k KaraokeScore) String() string {
// 	switch {
// 	case k < 25:
// 		return "Okay"
// 	case k < 50:
// 		return "Good"
// 	case k < 75:
// 		return "Great"
// 	}

// 	return "Perfect"
// }

// const (
// 	KaraokeScoreOkay    KaraokeScore = 10
// 	KaraokeScoreGood    KaraokeScore = 20
// 	KaraokeScoreGreat   KaraokeScore = 40
// 	KaraokeScorePerfect KaraokeScore = 80
// )

// func Score(x float64) KaraokeScore {
// 	if x == 0 {
// 		return 0
// 	}

// 	delta := gomath.Abs((x + 50) - (windowWidth - karaCenter))

// 	switch {
// 	case delta < 25:
// 		return KaraokeScoreOkay
// 	case delta < 50:
// 		return KaraokeScoreGood
// 	case delta < 75:
// 		return KaraokeScoreGreat
// 	}

// 	return KaraokeScorePerfect
// }

// type KaraokeBackground struct {
// 	Duration DurationMil `json:"duration"`
// 	FadeIn   DurationMil `json:"fade_in"`
// 	Image    string      `json:"image"`
// }

// type KaraokeSession struct {
// 	Inputs        []*KaraokeInput      `json:"inputs"`
// 	Backgrounds   []*KaraokeBackground `json:"backgrounds"`
// 	Sounds        map[string]string    `json:"sounds"`
// 	Music         string               `json:"music"`
// 	SampleRate    int                  `json:"sampleRate"`
// 	backgroundIdx int
// }
