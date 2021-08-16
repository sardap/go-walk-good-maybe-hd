package entity

import (
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
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
		},
		TextComponent: &components.TextComponent{},
	}
}
