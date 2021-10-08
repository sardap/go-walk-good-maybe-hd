package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/pkg/assets"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
	"github.com/sardap/walk-good-maybe-hd/pkg/math"
)

type SpeedLine struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.DestoryOnAnimeComponent
	*components.ImageComponent
	*components.ScrollableComponent
	*components.VelocityComponent
}

func CreateSpeedLine() *SpeedLine {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageSpeedLine)

	return &SpeedLine{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageBiscuitEnemyIdleTileSet.FrameWidth),
				Y: float64(img.Bounds().Dy()),
			},
		},
		DestoryOnAnimeComponent: &components.DestoryOnAnimeComponent{
			CyclesTilDeath: 1,
		},
		ImageComponent: &components.ImageComponent{
			Active: true,
			Image:  img,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 4,
		},
		VelocityComponent: &components.VelocityComponent{},
	}
}
