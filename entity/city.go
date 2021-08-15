package entity

import (
	"image"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
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
			Vel:  &image.Point{},
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
