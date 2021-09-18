package entity

import (
	gomath "math"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
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

type InputEnt struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.MovementComponent
	*components.InputComponent
}

func CreateDebugInput() *InputEnt {
	result := &InputEnt{
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

func CreateMenuInput() *InputEnt {
	result := &InputEnt{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		MovementComponent:  &components.MovementComponent{},
		InputComponent: &components.InputComponent{
			InputMode: components.InputModeKeyboard,
			Keyboard:  components.DefaultKeyboardInputType(),
			Gamepad:   components.DefaultGamepadInputType(),
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

type BasicTileMap struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.TileImageComponent
}

func CreateLifeDisplay() *BasicTileMap {
	img, _ := assets.LoadEbitenImage(assets.ImageUiLifeAmountTileSet)
	tileMap := components.CreateTileMap(1, 1, img, assets.ImageUiLifeAmountTileSet.FrameWidth)

	return &BasicTileMap{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageUiLifeAmountTileSet.FrameWidth),
				Y: float64(assets.ImageUiLifeAmountTileSet.FrameWidth),
			},
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
	}
}

func CreateJumpDisplay() *BasicTileMap {
	img, _ := assets.LoadEbitenImage(assets.ImageUiJumpAmountTileSet)
	tileMap := components.CreateTileMap(1, 1, img, assets.ImageUiJumpAmountTileSet.FrameWidth)

	return &BasicTileMap{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageUiJumpAmountTileSet.FrameWidth),
				Y: float64(assets.ImageUiJumpAmountTileSet.FrameWidth),
			},
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
	}
}

func CreateSpeedDisplay() *BasicTileMap {
	img, _ := assets.LoadEbitenImage(assets.ImageUiSpeedAmountTileSet)
	tileMap := components.CreateTileMap(1, 1, img, assets.ImageUiSpeedAmountTileSet.FrameWidth)

	return &BasicTileMap{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageUiSpeedAmountTileSet.FrameWidth),
				Y: float64(assets.ImageUiSpeedAmountTileSet.FrameWidth),
			},
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
	}
}
