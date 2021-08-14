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
}

func Createplayer() *Player {
	img, _ := assets.LoadImage(assets.ImageWhaleAir)

	img = assets.ScaleImage(img)

	result := &Player{
		TransformComponent: &components.TransformComponent{
			DrawImageOptions: &ebiten.DrawImageOptions{},
		},
		ImageComponent: &components.ImageComponent{
			Image: ebiten.NewImageFromImage(img),
		},
		AnimeComponent: &components.AnimeComponent{
			FrameWidth:    8 * 8 * 2,
			FrameHeight:   img.Bounds().Dy(),
			FrameDuration: 50 * time.Millisecond,
		},
	}

	return result
}
