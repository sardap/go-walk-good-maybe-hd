package game

import (
	"fmt"
	"math"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
)

type Inputable interface {
	ecs.BasicFace
	components.MovementFace
	components.InputFace
}

type InputSystem struct {
	ents    map[uint64]Inputable
	infoEnt *entity.InputInfo
}

func CreateInputSystem() *InputSystem {
	return &InputSystem{}
}

func (s *InputSystem) Priority() int {
	return int(systemPriorityInputSystem)
}

func (s *InputSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Inputable)

	s.infoEnt = entity.CreateInputInfo()
	s.infoEnt.Text = ""
	s.infoEnt.Postion.X = 300
	s.infoEnt.Postion.Y = 10
	world.AddEntity(s.infoEnt)
}

func (s *InputSystem) setInputMode(com *components.InputComponent, mode components.InputMode) {
	s.infoEnt.TextComponent.Text = fmt.Sprintf("Current:%s Change with K or G", mode.String())
	com.InputMode = mode
}

func (s *InputSystem) processGamepad(ent Inputable) {
	inputCom := ent.GetInputComponent()
	driver := inputCom.Gamepad.Driver
	gamepad := inputCom.Gamepad

	if driver.Ready(&inputCom.Gamepad) {
		if len(inpututil.JustConnectedGamepadIDs()) <= 0 {
			s.setInputMode(ent.GetInputComponent(), components.InputModeKeyboard)
			return
		}
		gamepad.Id = inpututil.JustConnectedGamepadIDs()[0]
	}

	id := gamepad.Id
	// maxButton := ebiten.GamepadButton(ebiten.GamepadButtonNum(id))
	// for b := ebiten.GamepadButton(id); b < maxButton; b++ {
	// 	// Log button events.
	// 	if inpututil.IsGamepadButtonJustPressed(id, b) {
	// 		log.Printf("button pressed: id: %d, button: %d", id, b)
	// 	}
	// 	if inpututil.IsGamepadButtonJustReleased(id, b) {
	// 		log.Printf("button released: id: %d, button: %d", id, b)
	// 	}
	// }

	vx := driver.GamepadAxis(id, gamepad.MoveAxisX)
	move := ent.GetMovementComponent()
	if math.Abs(vx) > float64(gamepad.MoveAxisX) {
		if vx > 0 {
			move.PressedDuration[components.InputKindMoveRight]++
		} else {
			if move.PressedDuration[components.InputKindMoveRight] > 0 {
				move.JustReleased[components.InputKindMoveRight] = true
			} else {
				move.JustReleased[components.InputKindMoveLeft] = false
			}
			move.PressedDuration[components.InputKindMoveRight] = 0
			move.JustPressed[components.InputKindMoveRight] = false
		}
		if vx < 0 {
			move.PressedDuration[components.InputKindMoveLeft]++
		} else {
			if move.PressedDuration[components.InputKindMoveLeft] > 0 {
				move.JustReleased[components.InputKindMoveLeft] = true
			} else {
				move.JustReleased[components.InputKindMoveLeft] = false
			}
			move.PressedDuration[components.InputKindMoveLeft] = 0
		}
	}

	for kind, btn := range gamepad.Mapping {
		move.PressedDuration[kind] = driver.GamepadButtonPressDuration(gamepad.Id, btn)
		move.JustPressed[kind] = driver.IsGamepadButtonJustPressed(gamepad.Id, btn)
		move.JustReleased[kind] = driver.IsGamepadButtonJustReleased(gamepad.Id, btn)
	}
}

func (s *InputSystem) processKeyboard(ent Inputable) {
	keyboard := ent.GetInputComponent().Keyboard
	driver := keyboard.Driver

	move := ent.GetMovementComponent()

	for kind, key := range keyboard.Mapping {
		move.PressedDuration[kind] = driver.KeyPressDuration(key)
		move.JustPressed[kind] = driver.IsKeyJustPressed(key)
		move.JustReleased[kind] = driver.IsKeyJustReleased(key)
	}
}

func (s *InputSystem) Update(dt float32) {
	for _, ent := range s.ents {
		// keyboard := ent.GetInputComponent().Keyboard
		// gamepad := ent.GetInputComponent().Gamepad

		// if gamepad.Driver != nil && keyboard.Driver.IsKeyJustReleased(keyboard.Mapping[components.InputKindChangeToGamepad]) {
		// 	ent.GetMovementComponent().PressedDuration[components.InputKindChangeToGamepad] = 1
		// } else if keyboard.Driver.KeyPressDuration(keyboard.Mapping[components.InputKindChangeToKeyboard]) {
		// 	ent.GetMovementComponent().PressedDuration.ChangeToKeyboard = true
		// }

		switch ent.GetInputComponent().InputMode {
		case components.InputModeGamepad:
			s.processGamepad(ent)
		case components.InputModeKeyboard:
			s.processKeyboard(ent)
		}

	}

}

func (s *InputSystem) Add(r Inputable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *InputSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *InputSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Inputable))
}
