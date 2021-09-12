package entity

import (
	gomath "math"

	"github.com/EngoEngine/ecs"
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

type KillBox struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.CollisionComponent
	*components.DamageComponent
	*components.IdentityComponent
}

func CreateKillBox() *KillBox {
	return &KillBox{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: 500,
				Y: 10000,
			},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		DamageComponent: &components.DamageComponent{
			BaseDamage: gomath.MaxFloat64,
		},
		IdentityComponent: &components.IdentityComponent{},
	}
}
