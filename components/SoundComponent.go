package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type SoundType int

const (
	SoundTypeMp3 SoundType = iota
)

type Sound struct {
	Source    []byte
	SoundType SoundType
}

type SoundComponent struct {
	Sound  Sound
	Intro  time.Duration
	Active bool
	Loop   bool
	// This should be null on creation
	Player *audio.Player
}
