package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/sardap/walk-good-maybe-hd/assets"
)

type Sound struct {
	Source     []byte
	SampleRate int
	SoundType  assets.SoundType
}

type SoundComponent struct {
	Sound   Sound
	Intro   time.Duration
	Active  bool
	Restart bool
	Loop    bool
	// This should be null on creation
	Player *audio.Player
}
