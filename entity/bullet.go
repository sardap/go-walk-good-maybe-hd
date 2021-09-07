package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type Bullet struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
	*components.VelocityComponent
	*components.CollisionComponent
	*components.ScrollableComponent
	*components.IdentityComponent
	*components.BulletComponent
}

func CreateBullet() *Bullet {
	img, _ := assets.LoadEbitenImage(assets.ImageBulletSmallGreen)

	result := &Bullet{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(img.Bounds().Dx()),
				Y: float64(img.Bounds().Dy()),
			},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{TagBullet},
		},
		BulletComponent: &components.BulletComponent{},
	}

	return result
}
