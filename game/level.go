package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type Level struct {
	StartX float64
}

type Building struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.VelocityComponent
	*components.TileImageComponent
	*components.CollisionComponent
	*components.ScrollableComponent
	*components.IdentityComponent
}

func createBuilding0(ent ecs.BasicEntity) *Building {
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

	result := &Building{
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
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{entity.TagGround},
		},
	}

	return result
}

func generateBuildings(w *ecs.World) {
	x := mainGameInfo.level.StartX
	for x < gameWidth/scaleMultiplier {
		ent := ecs.NewBasic()
		block := createBuilding0(ent)
		trans := block.GetTransformComponent()
		block.GetTransformComponent().Postion.Y = gameHeight/scaleMultiplier - trans.Size.Y
		block.GetTransformComponent().Postion.X = x
		x += block.Size.X + float64(utility.RandRange(30, 50))
		w.AddEntity(block)
	}
	mainGameInfo.level.StartX = x
}
