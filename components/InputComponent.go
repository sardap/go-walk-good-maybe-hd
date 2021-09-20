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

type InputKind int

const (
	// Menu Stuff
	InputKindSelect InputKind = iota
	// Movement
	InputKindMoveUp
	InputKindMoveDown
	InputKindMoveLeft
	InputKindMoveRight
	InputKindJump
	InputKindShoot
	// Misc
	InputKindChangeToGamepad
	InputKindChangeToKeyboard
	// Debug
	InputKindFastGameSpeed
	InputKindToggleCollsionOverlay
	InputKindLength
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
	IsKeyJustPressed(ebiten.Key) bool
	IsKeyJustReleased(ebiten.Key) bool
}

type EbitenKeyboardDriver struct {
}

func (EbitenKeyboardDriver) KeyPressDuration(k ebiten.Key) int {
	return inpututil.KeyPressDuration(k)
}

func (EbitenKeyboardDriver) IsKeyJustPressed(k ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(k)
}

func (EbitenKeyboardDriver) IsKeyJustReleased(k ebiten.Key) bool {
	return inpututil.IsKeyJustReleased(k)
}

type KeyboardInputType struct {
	Mapping map[InputKind]ebiten.Key
	Driver  KeyboardDriver
}

func DefaultKeyboardInputType() KeyboardInputType {
	return KeyboardInputType{
		Mapping: map[InputKind]ebiten.Key{
			InputKindMoveLeft:              ebiten.KeyLeft,
			InputKindMoveRight:             ebiten.KeyRight,
			InputKindMoveUp:                ebiten.KeyUp,
			InputKindMoveDown:              ebiten.KeyDown,
			InputKindChangeToGamepad:       ebiten.KeyG,
			InputKindChangeToKeyboard:      ebiten.KeyK,
			InputKindJump:                  ebiten.KeyZ,
			InputKindShoot:                 ebiten.KeyX,
			InputKindFastGameSpeed:         ebiten.KeyTab,
			InputKindToggleCollsionOverlay: ebiten.KeyO,
			InputKindSelect:                ebiten.KeyZ,
		},
		Driver: EbitenKeyboardDriver{},
	}
}

type GamepadDriver interface {
	GamepadAxis(ebiten.GamepadID, int) float64
	GamepadButtonPressDuration(ebiten.GamepadID, ebiten.GamepadButton) int
	IsGamepadButtonJustPressed(ebiten.GamepadID, ebiten.GamepadButton) bool
	IsGamepadButtonJustReleased(ebiten.GamepadID, ebiten.GamepadButton) bool
	Ready(g *GamepadInputType) bool
}

type EbitenGamepadDriver struct {
}

func (EbitenGamepadDriver) GamepadAxis(id ebiten.GamepadID, axis int) float64 {
	return ebiten.GamepadAxis(id, axis)
}

func (EbitenGamepadDriver) GamepadButtonPressDuration(id ebiten.GamepadID, btn ebiten.GamepadButton) int {
	return inpututil.GamepadButtonPressDuration(id, btn)
}

func (EbitenGamepadDriver) IsGamepadButtonJustPressed(id ebiten.GamepadID, btn ebiten.GamepadButton) bool {
	return inpututil.IsGamepadButtonJustPressed(id, btn)
}

func (EbitenGamepadDriver) IsGamepadButtonJustReleased(id ebiten.GamepadID, btn ebiten.GamepadButton) bool {
	return inpututil.IsGamepadButtonJustReleased(id, btn)
}

func (EbitenGamepadDriver) Ready(g *GamepadInputType) bool {
	return g.Id < 0 || inpututil.IsGamepadJustDisconnected(g.Id)
}

type GamepadInputType struct {
	Id        ebiten.GamepadID
	MoveAxisX int
	MoveAxisY int
	DeadZone  float64

	Mapping map[InputKind]ebiten.GamepadButton

	Driver GamepadDriver
}

func DefaultGamepadInputType() GamepadInputType {
	return GamepadInputType{
		Id:        ebiten.GamepadID(-1),
		MoveAxisX: 0,
		MoveAxisY: 1,
		DeadZone:  0.1,
		Mapping: map[InputKind]ebiten.GamepadButton{
			InputKindJump:   0,
			InputKindShoot:  1,
			InputKindSelect: 0,
		},
		Driver: EbitenGamepadDriver{},
	}
}

type InputComponent struct {
	InputMode InputMode
	Keyboard  KeyboardInputType
	Gamepad   GamepadInputType
}
