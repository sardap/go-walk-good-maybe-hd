package entity

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sardap/ecs"
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

type TestBox struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.ImageComponent
	*components.CollisionComponent
}

func CreateTestBox() *TestBox {
	rect := ebiten.NewImage(20, 50)
	ebitenutil.DrawRect(rect, 0, 0, 20, 50, color.RGBA{0, 0, 0, 255})

	result := &TestBox{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			GeoM: &ebiten.GeoM{},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ImageComponent: &components.ImageComponent{
			Image: rect,
		},
	}

	result.TransformComponent.Scale(20, 20)

	return result
}
