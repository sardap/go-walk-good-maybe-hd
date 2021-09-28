package common

import (
	"time"
)

type KaraokeInput struct {
	StartTime  time.Duration
	Duration   time.Duration
	Sound      string
	XPostion   float64
	XSpeed     float64
	HitPostion float64
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
