package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type Token struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.CollisionComponent
	*components.IdentityComponent
	*components.LifeComponent
	*components.ScrollableComponent
	*components.ImageComponent
	*components.VelocityComponent
}

func createToken(img *ebiten.Image, tag string) *Token {
	return &Token{
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
			Tags: []string{tag},
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
}

func CreateJumpUpToken() *Token {
	img, _ := assets.LoadEbitenImage(assets.ImageTokenJumpUp)
	return createToken(img, TagJumpToken)
}

func CreateSpeedUpToken() *Token {
	img, _ := assets.LoadEbitenImage(assets.ImageTokenSpeedUp)
	return createToken(img, TagSpeedToken)
}
