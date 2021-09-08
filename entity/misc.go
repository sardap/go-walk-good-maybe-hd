package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
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

type SingleScrollableAnime struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.AnimeComponent
	*components.DestoryOnAnimeComponent
	*components.TileImageComponent
	*components.ScrollableComponent
	*components.VelocityComponent
}
