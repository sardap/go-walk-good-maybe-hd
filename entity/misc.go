package entity

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type SoundPlayer struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.SoundComponent
}

func CreateSoundPlayer(asset interface{}) *SoundPlayer {
	return &SoundPlayer{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		SoundComponent: &components.SoundComponent{
			Sound: components.LoadSound(asset),
		},
	}
}

type InputInfo struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.TextComponent
}

func CreateInputInfo() *InputInfo {
	return &InputInfo{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		TextComponent:      &components.TextComponent{},
	}
}

type TestBox struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
	*components.CollisionComponent
}

func CreateTestBox() *TestBox {
	rect := ebiten.NewImage(20, 50)
	ebitenutil.DrawRect(rect, 0, 0, 20, 50, color.RGBA{0, 0, 0, 255})

	w, h := rect.Size()

	result := &TestBox{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: float64(w), Y: float64(h)},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ImageComponent: &components.ImageComponent{
			Image: rect,
		},
	}

	return result
}

type DebugInput struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.MovementComponent
	*components.InputComponent
}

func CreateDebugInput() *DebugInput {
	result := &DebugInput{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		MovementComponent:  &components.MovementComponent{},
		InputComponent: &components.InputComponent{
			InputMode: components.InputModeKeyboard,
			Keyboard:  components.DefaultKeyboardInputType(),
		},
	}

	return result
}
