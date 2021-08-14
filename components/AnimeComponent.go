package components

import "time"

type AnimeComponent struct {
	FrameWidth     int
	FrameHeight    int
	CurrentFrame   int
	FrameDuration  time.Duration
	FrameRemaining time.Duration
}
