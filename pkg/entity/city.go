package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/pkg/assets"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
	"github.com/sardap/walk-good-maybe-hd/pkg/math"
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
			Sound:  components.LoadSound(assets.MusicPdCity0),
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
	img, _ := assets.LoadEbitenImageAsset(assets.ImageBackgroundCity)

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
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 0.25,
		},
		WrapComponent: &components.WrapComponent{
			Max: math.Vector2{X: float64(w), Y: 0},
			Min: math.Vector2{X: -float64(w), Y: 0},
		},
	}
}

type CitySkyBackground struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
}

func CreateCitySkyBackground() *CitySkyBackground {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageSkyCity)

	w, h := img.Size()

	return &CitySkyBackground{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: float64(w), Y: float64(h)},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
	}
}

type CityFogBackground struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.IdentityComponent
	*components.ImageComponent
	*components.VelocityComponent
	*components.ScrollableComponent
	*components.WrapComponent
}

func CreateCityFogBackground() *CityFogBackground {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageCityFog)

	w, h := img.Size()

	return &CityFogBackground{
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
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 0.175,
		},
		WrapComponent: &components.WrapComponent{
			Max: math.Vector2{X: float64(w), Y: 0},
			Min: math.Vector2{X: -float64(w), Y: 0},
		},
	}
}
