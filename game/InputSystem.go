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

	if driver.Ready(&inputCom.Gamepad) {
		if len(inpututil.JustConnectedGamepadIDs()) <= 0 {
			s.setInputMode(ent.GetInputComponent(), components.InputModeKeyboard)
			return
		}
		inputCom.Gamepad.Id = inpututil.JustConnectedGamepadIDs()[0]
	}

	id := inputCom.Gamepad.Id
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

	vx := driver.GamepadAxis(id, inputCom.Gamepad.MoveAxisX)
	move := ent.GetMovementComponent()
	if math.Abs(vx) > float64(inputCom.Gamepad.MoveAxisX) {
		if vx > 0 {
			move.MoveRight = true
		}
		if vx < 0 {
			move.MoveLeft = true
		}
	}

	if driver.IsGamepadButtonJustPressed(id, inputCom.Gamepad.ButtonJump) {
		move.MoveUp = true
	}
	if driver.IsGamepadButtonJustPressed(id, inputCom.Gamepad.ButtonShoot) {
		move.Shoot = true
	}

}

func (s *InputSystem) processKeyboard(ent Inputable) {
	keyboard := ent.GetInputComponent().Keyboard
	driver := keyboard.Driver

	move := ent.GetMovementComponent()
	if driver.KeyPressDuration(keyboard.KeyMoveLeft) > 0 {
		move.MoveLeft = true
	}

	if driver.KeyPressDuration(keyboard.KeyMoveRight) > 0 {
		move.MoveRight = true
	}

	if driver.KeyPressDuration(keyboard.KeyMoveUp) > 0 || driver.KeyPressDuration(keyboard.KeyJump) > 0 {
		move.MoveUp = true
	}

	if driver.KeyPressDuration(keyboard.KeyMoveDown) > 0 {
		move.MoveDown = true
	}

	if driver.KeyPressDuration(keyboard.KeyShoot) > 0 {
		move.Shoot = true
	}

	if driver.KeyPressDuration(keyboard.KeyFastGameMode) > 0 {
		move.FastGameSpeed = true
	}

	if driver.KeyPressDuration(keyboard.KeyToggleCollsionOverlay) > 0 {
		move.ToggleCollsionOverlay = true
	}

	if driver.KeyPressDuration(keyboard.KeySelect) > 0 {
		move.Select = true
	}
}

func (s *InputSystem) Update(dt float32) {
	for _, ent := range s.ents {
		keyboard := ent.GetInputComponent().Keyboard
		gamepad := ent.GetInputComponent().Gamepad

		if gamepad.Driver != nil && keyboard.Driver.KeyPressDuration(keyboard.KeyChangeToGamepad) > 0 {
			ent.GetMovementComponent().ChangeToGamepad = true
		} else if keyboard.Driver.KeyPressDuration(keyboard.KeyChangeToKeyboard) > 0 {
			ent.GetMovementComponent().ChangeToKeyboard = true
		}

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
