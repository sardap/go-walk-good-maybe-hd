package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type Player struct {
	ecs.BasicEntity
	*components.MainGamePlayerComponent
	*components.TransformComponent
	*components.TileImageComponent
	*components.AnimeComponent
	*components.MovementComponent
	*components.VelocityComponent
	*components.CollisionComponent
	*components.ScrollableComponent
	*components.GravityComponent
	*components.IdentityComponent
	*components.SoundComponent
	*components.InputComponent
}

func CreatePlayer() *Player {
	img, _ := assets.LoadEbitenImage(assets.ImageWhaleAirTileSet)

	tileMap := components.CreateTileMap(1, 1, img, assets.ImageWhaleAirTileSet.FrameWidth)

	result := &Player{
		BasicEntity: ecs.NewBasic(),
		MainGamePlayerComponent: &components.MainGamePlayerComponent{
			Speed:         700,
			JumpPower:     900,
			State:         components.MainGamePlayerStateFlying,
			ShootCooldown: 250 * time.Millisecond,
		},
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageWhaleAirTileSet.FrameWidth),
				Y: float64(assets.ImageWhaleAirTileSet.FrameWidth),
			},
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration: 50 * time.Millisecond,
		},
		MovementComponent: &components.MovementComponent{},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ScrollableComponent: &components.ScrollableComponent{},
		GravityComponent:    &components.GravityComponent{},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{TagPlayer},
		},
		SoundComponent: &components.SoundComponent{
			Active: false,
		},
		InputComponent: &components.InputComponent{
			InputMode: components.InputModeKeyboard,
			Gamepad:   components.DefaultGamepadInputType(),
			Keyboard:  components.DefaultKeyboardInputType(),
		},
	}

	return result
}
