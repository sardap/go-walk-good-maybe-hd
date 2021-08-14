package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type Player struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
}

func Createplayer() *Player {
	img, _ := assets.LoadImage(assets.ImageWhaleIdle)

	result := &Player{
		TransformComponent: &components.TransformComponent{
			DrawImageOptions: &ebiten.DrawImageOptions{},
		},
		ImageComponent: &components.ImageComponent{
			Image: ebiten.NewImageFromImage(assets.ScaleImage(img)),
		},
	}

	return result
}
