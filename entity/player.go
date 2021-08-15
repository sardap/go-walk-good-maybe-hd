package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type Player struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
	*components.AnimeComponent
	*components.MovementComponent
}

func CreatePlayer() *Player {
	img, _ := assets.LoadImage(assets.ImageWhaleAir)

	img = assets.ScaleImage(img)

	result := &Player{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
		},
		ImageComponent: &components.ImageComponent{
			Image: ebiten.NewImageFromImage(img),
		},
		AnimeComponent: &components.AnimeComponent{
			FrameWidth:    8 * 8 * 2,
			FrameHeight:   img.Bounds().Dy(),
			FrameDuration: 50 * time.Millisecond,
		},
		MovementComponent: &components.MovementComponent{},
	}

	return result
}
