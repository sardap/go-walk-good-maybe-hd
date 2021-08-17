package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
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
}

func CreatePlayer() *Player {
	img, _ := assets.LoadImage([]byte(assets.ImageWhaleAir.Data))

	result := &Player{
		BasicEntity:             ecs.NewBasic(),
		MainGamePlayerComponent: &components.MainGamePlayerComponent{},
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
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
		MovementComponent: &components.MovementComponent{},
		VelocityComponent: &components.VelocityComponent{
			Vel: &math.Vector2{},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
	}

	return result
}
