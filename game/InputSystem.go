package game

import (
	"log"
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sardap/walk-good-maybe-hd/components"
)

const (
	moveAxisX = 2
	moveAxisY = 3
	deadZone  = 0.1
)

type inputMode int

const (
	inputModeGamepad = iota
	inputModeKeyboard
)

type InputSystem struct {
	ents          map[uint64]Inputable
	playerGamepad ebiten.GamepadID
	inputMode     inputMode
}

func CreateInputSystem() *InputSystem {
	return &InputSystem{}
}

func (s *InputSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Inputable)
}

func (s *InputSystem) processGamepad() {
	if inpututil.IsGamepadJustDisconnected(s.playerGamepad) {
		if len(inpututil.JustConnectedGamepadIDs()) <= 0 {
			s.inputMode = inputModeKeyboard
			return
		}
		s.playerGamepad = inpututil.JustConnectedGamepadIDs()[0]
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyK) {
		s.inputMode = inputModeKeyboard
		return
	}

	id := s.playerGamepad
	maxButton := ebiten.GamepadButton(ebiten.GamepadButtonNum(id))
	for b := ebiten.GamepadButton(id); b < maxButton; b++ {
		// Log button events.
		if inpututil.IsGamepadButtonJustPressed(id, b) {
			log.Printf("button pressed: id: %d, button: %d", id, b)
		}
		if inpututil.IsGamepadButtonJustReleased(id, b) {
			log.Printf("button released: id: %d, button: %d", id, b)
		}
	}

	vx := ebiten.GamepadAxis(id, moveAxisX)
	vy := ebiten.GamepadAxis(id, moveAxisY)
	for _, ent := range s.ents {
		move := ent.GetMovementComponent()
		if math.Abs(vx) > deadZone {
			if vx > 0 {
				move.MoveRight = true
			}
			if vx < 0 {
				move.MoveLeft = true
			}
		}

		if math.Abs(vy) > deadZone {
			if vy > 0 {
				move.MoveDown = true
			}
			if vy < 0 {
				move.MoveUp = true
			}
		}
	}
}

func (s *InputSystem) processKeyboard() {
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		s.inputMode = inputModeGamepad
		return
	}

	for _, ent := range s.ents {
		move := ent.GetMovementComponent()
		if inpututil.KeyPressDuration(ebiten.KeyLeft) > 0 {
			move.MoveLeft = true
		}

		if inpututil.KeyPressDuration(ebiten.KeyRight) > 0 {
			move.MoveRight = true
		}

		if inpututil.KeyPressDuration(ebiten.KeyUp) > 0 {
			move.MoveUp = true
		}

		if inpututil.KeyPressDuration(ebiten.KeyDown) > 0 {
			move.MoveDown = true
		}
	}
}

func (s *InputSystem) Update(dt float32) {
	switch s.inputMode {
	case inputModeGamepad:
		s.processGamepad()
	case inputModeKeyboard:
		s.processKeyboard()
	}
}

func (s *InputSystem) Add(r Inputable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *InputSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type Inputable interface {
	ecs.BasicFace
	components.MovementFace
}

func (s *InputSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Inputable))
}
