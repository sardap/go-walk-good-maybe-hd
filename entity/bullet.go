package entity

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
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

func CreateBullet(img *ebiten.Image) *Bullet {
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
			Tags: []int{TagBullet},
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

func CreatePlayerBullet() *Bullet {
	img, _ := assets.LoadEbitenImage(assets.ImageBulletSmallGreen)

	return CreateBullet(img)
}

func CreateEnemyBullet() *Bullet {
	img, _ := assets.LoadEbitenImageColorSwap(
		assets.ImageBulletSmallGreen,
		map[color.RGBA]color.RGBA{
			{R: 75, G: 205, B: 75, A: 255}: {R: 205, G: 75, B: 75, A: 255},
			{R: 72, G: 150, B: 72, A: 255}: {R: 150, G: 72, B: 72, A: 255},
			{R: 85, G: 185, B: 85, A: 255}: {R: 185, G: 85, B: 85, A: 255},
		},
	)

	return CreateBullet(img)
}
