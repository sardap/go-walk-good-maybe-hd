package entity

import (
	"time"

	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type Player struct {
	ecs.BasicEntity
	*components.MainGamePlayerComponent
	*components.TransformComponent
	*components.ImageComponent
	*components.AnimeComponent
	*components.MovementComponent
	*components.VelocityComponent
	*components.CollisionComponent
	*components.ScrollableComponent
	*components.GravityComponent
	*components.IdentityComponent
}

func CreatePlayer() *Player {
	img, _ := assets.LoadEbitenImage(assets.ImageWhaleAirTileSet)

	result := &Player{
		BasicEntity: ecs.NewBasic(),
		MainGamePlayerComponent: &components.MainGamePlayerComponent{
			Speed:     40,
			JumpPower: 30,
			State:     components.MainGamePlayerStateFalling,
		},
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: float64(assets.ImageWhaleAirTileSet.FrameWidth), Y: float64(assets.ImageWhaleAirTileSet.FrameWidth)},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		AnimeComponent: &components.AnimeComponent{
			FrameWidth:    assets.ImageWhaleAirTileSet.FrameWidth,
			FrameHeight:   img.Bounds().Dy(),
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
	}

	return result
}
