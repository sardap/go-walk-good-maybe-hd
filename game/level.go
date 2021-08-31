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

type LevelBlock struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.VelocityComponent
	*components.TileImageComponent
	*components.CollisionComponent
	*components.ScrollableComponent
	*components.IdentityComponent
}

func createBuilding0(ent ecs.BasicEntity) *LevelBlock {
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

	result := &LevelBlock{
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

type LevelBlockable interface {
	ecs.BasicFace
	components.TransformFace
}

func populateLevelBlock(w *ecs.World, lb LevelBlockable) {
	trans := lb.GetTransformComponent()
	biscuit := entity.CreateBiscuitEnemy()
	biscuit.Postion.X = utility.RandRangeFloat64(int(trans.Postion.X), int(trans.Postion.X+trans.Size.X))
	biscuit.Postion.Y = 30
	biscuit.Layer = enemyLayer
	w.AddEntity(biscuit)
}

func generateCityBuildings(mainGameInfo *MainGameInfo, w *ecs.World) {
	x := mainGameInfo.Level.StartX
	for x < gameWidth/scaleMultiplier {
		ent := ecs.NewBasic()
		levelBlock := createBuilding0(ent)
		trans := levelBlock.GetTransformComponent()
		levelBlock.GetTransformComponent().Postion.Y = gameHeight/scaleMultiplier - trans.Size.Y
		levelBlock.GetTransformComponent().Postion.X = x
		x += levelBlock.Size.X + float64(utility.RandRange(30, 50))
		w.AddEntity(levelBlock)
		populateLevelBlock(w, levelBlock)
	}
	mainGameInfo.Level.StartX = x
}
