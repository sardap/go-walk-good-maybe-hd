package game

import (
	"math/rand"

	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type block interface {
	ecs.BasicFace
	ecs.Identifier
	components.TransformFace
}

type windowBlock struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.VelocityComponent
	*components.TileImageComponent
	*components.CollisionComponent
	*components.ScrollableComponent
}

func createBuilding0(ent ecs.BasicEntity) block {
	tileSet, _ := assets.LoadEbitenImage(assets.ImageBuilding0TileSet)

	tileMap := make([]int, 10)
	for i := range tileMap {
		if rand.Int()%2 == 0 {
			tileMap[i] = assets.IndexBuilding0Wall
		} else {
			tileMap[i] = assets.IndexBuilding0Window
		}
	}

	w := float64(5 * assets.ImageBuilding0TileSet.FrameWidth)
	h := float64(2 * assets.ImageBuilding0TileSet.FrameWidth)

	window := &windowBlock{
		BasicEntity: ent,
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{X: w, Y: h},
		},
		VelocityComponent: &components.VelocityComponent{},
		TileImageComponent: &components.TileImageComponent{
			TilesImg:  tileSet,
			TilesMap:  tileMap,
			TileWidth: assets.ImageBuilding0TileSet.FrameWidth,
			TileXNum:  5,
			Layer:     buildingForground,
		},
		CollisionComponent:  &components.CollisionComponent{},
		ScrollableComponent: &components.ScrollableComponent{},
	}

	return window
}
