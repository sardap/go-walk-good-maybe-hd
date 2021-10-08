package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/sardap/walk-good-maybe-hd/pkg/assets"
)

type Sound struct {
	Source     []byte
	SampleRate int
	SoundType  assets.SoundType
	Volume     float64
}

type SoundComponent struct {
	Sound   *Sound
	Intro   time.Duration
	Active  bool
	Restart bool
	Loop    bool
	// This should be null on creation
	Player *audio.Player
}

func (s *SoundComponent) ChangeSound(newSound *Sound) {
	if s.Player != nil {
		s.Player.Pause()
		s.Player.Close()
	}

	s.Sound = newSound
}
