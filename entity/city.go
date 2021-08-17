package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
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
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
		},
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
	*components.ImageComponent
	*components.VelocityComponent
}

func CreateCityBackground() *CityBackground {
	img, _ := assets.LoadImage([]byte(assets.ImageBackgroundCity.Data))

	return &CityBackground{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image: ebiten.NewImageFromImage(img),
		},
		VelocityComponent: &components.VelocityComponent{
			Vel: &math.Vector2{},
		},
	}
}
