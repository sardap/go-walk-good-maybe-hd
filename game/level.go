package game

import (
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
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

	width := utility.RandRange(5, 9)
	height := utility.RandRange(5, 10)

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding0TileSet.FrameWidth)
	// Edges
	tileMap.SetRow(0, 0, assets.IndexBuilding0Wall)
	tileMap.SetCol(width-1, 1, assets.IndexBuilding0Wall)
	// Body
	tileMap.SetCol(0, 1, assets.IndexBuilding0Wall)
	for x := 1; x < width-1; x++ {
		tileMap.SetCol(x, 1, assets.IndexBuilding0Window)
	}

	window := &windowBlock{
		BasicEntity: ent,
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(width * assets.ImageBuilding0TileSet.FrameWidth),
				Y: float64(height * assets.ImageBuilding0TileSet.FrameWidth),
			},
		},
		VelocityComponent: &components.VelocityComponent{},
		TileImageComponent: &components.TileImageComponent{
			TileMap: tileMap,
			Layer:   buildingForground,
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ScrollableComponent: &components.ScrollableComponent{},
	}

	return window
}
