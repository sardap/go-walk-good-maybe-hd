package components

import (
	"time"
)

type KaraokeSound string

const (
	KaraokeSoundA = "A"
	KaraokeSoundB = "B"
	KaraokeSoundX = "X"
	KaraokeSoundY = "Y"
)

type KaraokeImageComponent struct {
	Start     time.Duration
	SoundType KaraokeSound
}
