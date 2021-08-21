package game

import "github.com/sardap/walk-good-maybe-hd/math"

var (
	mainGameInfo *MainGameInfo
)

type gameState int

const (
	gameStateStarting gameState = iota
	gameStateScrolling
)

type MainGameInfo struct {
	scrollingSpeed math.Vector2
	state          gameState
}
