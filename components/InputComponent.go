package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type InputMode int

const (
	InputModeGamepad InputMode = iota
	InputModeKeyboard
)

func (i InputMode) String() string {
	switch i {
	case InputModeGamepad:
		return "gamepad"
	case InputModeKeyboard:
		return "keyboard"
	}

	panic("Unknown input type")
}

type KeyboardDriver interface {
	KeyPressDuration(ebiten.Key) int
}

type EbitenKeyboardDriver struct {
}

func (EbitenKeyboardDriver) KeyPressDuration(k ebiten.Key) int {
	return inpututil.KeyPressDuration(k)
}

type KeyboardInputType struct {
	KeyMoveUp    ebiten.Key
	KeyMoveDown  ebiten.Key
	KeyMoveLeft  ebiten.Key
	KeyMoveRight ebiten.Key
	// Main game stuff
	KeyJump  ebiten.Key
	KeyShoot ebiten.Key
	// Main game Debug input
	KeyScrollSpeedUp         ebiten.Key
	KeyToggleCollsionOverlay ebiten.Key

	Driver KeyboardDriver
}

func DefaultKeyboardInputType() KeyboardInputType {
	return KeyboardInputType{
		KeyMoveLeft:              ebiten.KeyLeft,
		KeyMoveRight:             ebiten.KeyRight,
		KeyMoveUp:                ebiten.KeyUp,
		KeyMoveDown:              ebiten.KeyDown,
		KeyJump:                  ebiten.KeyZ,
		KeyShoot:                 ebiten.KeyX,
		KeyScrollSpeedUp:         ebiten.KeyTab,
		KeyToggleCollsionOverlay: ebiten.KeyO,
		Driver:                   EbitenKeyboardDriver{},
	}
}

type GamepadDriver interface {
	GamepadAxis(ebiten.GamepadID, int) float64
	IsGamepadButtonJustPressed(ebiten.GamepadID, ebiten.GamepadButton) bool
}

type EbitenGamepadDriver struct {
}

func (EbitenGamepadDriver) GamepadAxis(id ebiten.GamepadID, axis int) float64 {
	return ebiten.GamepadAxis(id, axis)
}

func (EbitenGamepadDriver) IsGamepadButtonJustPressed(id ebiten.GamepadID, btn ebiten.GamepadButton) bool {
	return inpututil.IsGamepadButtonJustPressed(id, btn)

}

type GamepadInputType struct {
	Id          ebiten.GamepadID
	MoveAxisX   int
	MoveAxisY   int
	DeadZone    float64
	ButtonJump  ebiten.GamepadButton
	ButtonShoot ebiten.GamepadButton

	Driver GamepadDriver
}

func DefaultGamepadInputType() GamepadInputType {
	return GamepadInputType{
		Id:          ebiten.GamepadID(-1),
		MoveAxisX:   0,
		MoveAxisY:   1,
		DeadZone:    0.1,
		ButtonJump:  0,
		ButtonShoot: 1,
		Driver:      EbitenGamepadDriver{},
	}
}

type InputComponent struct {
	InputMode InputMode
	Keyboard  KeyboardInputType
	Gamepad   GamepadInputType
}
