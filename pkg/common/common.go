package common

import (
	"time"
)

type KaraokeInput struct {
	TargetHitTime time.Duration
	Sound         string
	HitTime       time.Duration
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
	MusicDuration    time.Duration
	SampleRate       int
	BackgroundIdx    int
}

type KaraokeIndex struct {
	KaraokeGames []string `toml:"karaoke_games"`
}
