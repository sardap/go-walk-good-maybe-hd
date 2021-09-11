package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type JumpUpToken struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.CollisionComponent
	*components.IdentityComponent
	*components.LifeComponent
	*components.ScrollableComponent
	*components.ImageComponent
	*components.VelocityComponent
}

func CreateJumpUpToken() *JumpUpToken {
	img, _ := assets.LoadEbitenImage(assets.ImageTokenJumpUp)

	result := &JumpUpToken{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(img.Bounds().Max.X),
				Y: float64(img.Bounds().Max.Y),
			},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{TagJumpToken},
		},
		LifeComponent: &components.LifeComponent{
			HP: 1,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		VelocityComponent: &components.VelocityComponent{},
	}

	return result
}
