package entity

import (
	"time"

	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type CityMusic struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.SoundComponent
}

func CreateCityMusic() *CityMusic {
	return &CityMusic{
		BasicEntity:        ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{},
		SoundComponent: &components.SoundComponent{
			Sound: components.Sound{
				Source:    assets.MusicPdCity0,
				SoundType: components.SoundTypeMp3,
			},
			Active: true,
			Loop:   true,
			Intro:  time.Duration(8) * time.Second,
		},
	}
}

type CityBackground struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.IdentityComponent
	*components.ImageComponent
	*components.VelocityComponent
	*components.ScrollableComponent
	*components.WrapComponent
}

func CreateCityBackground() *CityBackground {
	img, _ := assets.LoadEbitenImage(assets.ImageBackgroundCity)

	w, h := img.Size()

	return &CityBackground{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: float64(w), Y: float64(h)},
		},
		IdentityComponent: &components.IdentityComponent{},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: math.Vector2{},
		},
		WrapComponent: &components.WrapComponent{
			Threshold: float64(w),
		},
	}
}
