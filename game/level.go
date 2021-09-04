package game

import (
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type Level struct {
	StartX float64
	Width  float64
	Height float64
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

func createBuilding1(ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImage(assets.ImageBuilding1TileSet)

	width := utility.RandRange(4, 6)
	height := utility.RandRange(5, 10)

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding0TileSet.FrameWidth)
	// Roof
	tileMap.SetTile(0, 0, assets.IndexBuilding1RoofLeft)
	tileMap.SetRow(1, 0, assets.IndexBuilding1RoofMiddle)
	tileMap.SetTile(width-1, 0, assets.IndexBuilding1RoofRight)
	// Body
	tileMap.SetCol(0, 1, assets.IndexBuilding1MiddleLeft)
	for x := 1; x < width-1; x++ {
		tileMap.SetCol(x, 1, assets.IndexBuilding1BotMiddle)
	}
	tileMap.SetCol(width-1, 1, assets.IndexBuilding1MiddleRight)

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
	if lb.GetTransformComponent().Postion.X < 50 {
		return
	}

	trans := lb.GetTransformComponent()
	biscuit := entity.CreateBiscuitEnemy()
	biscuit.Postion.X = utility.RandRangeFloat64(int(trans.Postion.X), int(trans.Postion.X+trans.Size.X))
	biscuit.Postion.Y = 30
	biscuit.Layer = enemyLayer
	w.AddEntity(biscuit)
}

func createLevelBlock(basic ecs.BasicEntity) *LevelBlock {
	val := rand.Float64()
	switch {
	case val <= 0.5:
		return createBuilding0(basic)
	case val > 0.5:
		return createBuilding1(basic)
	}

	panic("random number bug")
}

func generateCityBuildings(mainGameInfo *MainGameInfo, w *ecs.World) {
	x := mainGameInfo.Level.StartX
	for x < mainGameInfo.Level.Width {
		ent := ecs.NewBasic()
		levelBlock := createLevelBlock(ent)
		trans := levelBlock.GetTransformComponent()
		levelBlock.GetTransformComponent().Postion.Y = mainGameInfo.Level.Height - trans.Size.Y
		levelBlock.GetTransformComponent().Postion.X = x
		x += levelBlock.Size.X + float64(utility.RandRange(minSpaceBetweenBuildings, minSpaceBetweenBuildings+20*scaleMultiplier))
		w.AddEntity(levelBlock)
		populateLevelBlock(w, levelBlock)
	}
	mainGameInfo.Level.StartX = x
}
