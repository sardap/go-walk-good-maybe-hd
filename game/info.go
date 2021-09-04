package game

import (
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type gameState int

const (
	gameStateStarting gameState = iota
	gameStateScrolling
)

type MainGameInfo struct {
	ScrollingSpeed math.Vector2
	Gravity        float64
	State          gameState
	Level          *Level
	InputEnt       *entity.DebugInput
}
