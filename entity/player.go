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
}

func CreatePlayer() *Player {
	img, _ := assets.LoadImage([]byte(assets.ImageWhaleAir.Data))

	result := &Player{
		BasicEntity:             ecs.NewBasic(),
		MainGamePlayerComponent: &components.MainGamePlayerComponent{},
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: float64(assets.ImageWhaleAir.FrameWidth), Y: float64(assets.ImageWhaleAir.FrameWidth)},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		AnimeComponent: &components.AnimeComponent{
			FrameWidth:    assets.ImageWhaleAir.FrameWidth,
			FrameHeight:   img.Bounds().Dy(),
			FrameDuration: 50 * time.Millisecond,
		},
		MovementComponent: &components.MovementComponent{
			Speed: 40,
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ScrollableComponent: &components.ScrollableComponent{},
	}

	return result
}
