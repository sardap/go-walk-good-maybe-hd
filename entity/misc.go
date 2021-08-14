package entity

import (
	"image"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type InputInfo struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.TextComponent
}

func CreateInputInfo() *InputInfo {
	return &InputInfo{
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
			Vel:  &image.Point{},
		},
		TextComponent: &components.TextComponent{},
	}
}
