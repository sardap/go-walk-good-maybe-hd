package entity

import (
	"image/png"
	"os"
	"path/filepath"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type Player struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
}

func Createplayer() *Player {
	playerFile, _ := os.Open(filepath.Join("assets", "player", "whale_small.png"))
	playerImage, _ := png.Decode(playerFile)

	result := &Player{
		TransformComponent: &components.TransformComponent{
			DrawImageOptions: &ebiten.DrawImageOptions{},
		},
		ImageComponent: &components.ImageComponent{
			Image: ebiten.NewImageFromImage(playerImage),
		},
	}

	return result
}
