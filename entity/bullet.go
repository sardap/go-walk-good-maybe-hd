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
	*components.BulletComponent
	*components.CollisionComponent
	*components.DamageComponent
	*components.LifeComponent
	*components.ScrollableComponent
	*components.IdentityComponent
	*components.ImageComponent
	*components.VelocityComponent
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
		BulletComponent: &components.BulletComponent{},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		DamageComponent: &components.DamageComponent{
			BaseDamage: 100,
		},
		LifeComponent: &components.LifeComponent{
			HP: 1,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{TagBullet},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{},
		},
	}

	return result
}
