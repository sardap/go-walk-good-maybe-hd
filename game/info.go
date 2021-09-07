package game

import (
	"math/rand"

	"github.com/SolarLune/resolv"
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

type Info struct {
	Rand         *rand.Rand
	Space        *resolv.Space
	MainGameInfo *MainGameInfo
}
