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

func createLevelBlock(ent ecs.BasicEntity, tileMap *components.TileMap, width, height int) *LevelBlock {
	return &LevelBlock{
		BasicEntity: ent,
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(width * tileMap.TileWidth),
				Y: float64(height * tileMap.TileWidth),
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
}

func createBuilding0(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImage(assets.ImageBuilding0TileSet)

	width := utility.RandRange(rand, 5, 9)
	height := utility.RandRange(rand, 5, 10)

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding0TileSet.FrameWidth)
	// Edges
	tileMap.SetRow(0, 0, assets.IndexBuilding0Wall)
	tileMap.SetCol(width-1, 1, assets.IndexBuilding0Wall)
	// Body
	tileMap.SetCol(0, 1, assets.IndexBuilding0Wall)
	for x := 1; x < width-1; x++ {
		tileMap.SetCol(x, 1, assets.IndexBuilding0Window)
	}

	return createLevelBlock(ent, tileMap, width, height)
}

func createBuilding1(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImage(assets.ImageBuilding1TileSet)

	width := utility.RandRange(rand, 4, 6)
	height := utility.RandRange(rand, 5, 10)

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding1TileSet.FrameWidth)
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

	return createLevelBlock(ent, tileMap, width, height)
}

func createBuilding2(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImage(assets.ImageBuilding2TileSet)

	width := utility.RandRange(rand, 4, 6)
	if width%2 == 0 {
		width++
	}

	height := utility.RandRange(rand, 5, 10)
	for height%3 != 0 {
		height++
	}
	height++

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding2TileSet.FrameWidth)
	// Left section
	tileMap.SetTile(0, 0, assets.IndexBuilding2LeftRoof)
	tileMap.SetCol(0, 1, assets.IndexBuilding2LeftMiddle)

	// Middle section
	for x := 1; x < width-1; x++ {
		if x%2 == 0 {
			tileMap.SetTile(x, 0, assets.IndexBuilding2RoofClean)
		} else {
			tileMap.SetTile(x, 0, assets.IndexBuilding2RoofHanging)
		}

		for y := 1; y < height; y += 3 {
			if x%2 == 0 {
				tileMap.SetTile(x, y, assets.IndexBuilding2Window)
				tileMap.SetTile(x, y+1, assets.IndexBuilding2Window)
			} else {
				tileMap.SetTile(x, y, assets.IndexBuilding2MiddleWhite)
				tileMap.SetTile(x, y+1, assets.IndexBuilding2MiddleWhiteClean)
			}

			tileMap.SetTile(x, y+2, assets.IndexBuilding2MiddleBlue)
		}
	}

	// Right section
	tileMap.SetTile(width-1, 0, assets.IndexBuilding2RightRoof)
	tileMap.SetCol(width-1, 1, assets.IndexBuilding2RightMiddle)

	return createLevelBlock(ent, tileMap, width, height)
}

func createBuilding3(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImage(assets.ImageBuilding3TileSet)

	width := utility.RandRange(rand, 3, 5)

	height := utility.RandRange(rand, 5, 8)
	for height%2 != 0 {
		height++
	}

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding3TileSet.FrameWidth)
	// Left section
	tileMap.SetTile(0, 0, assets.IndexBuilding3RoofTop)
	tileMap.SetCol(0, 1, assets.IndexBuilding3RoofLeft)
	tileMap.SetCol(0, 2, assets.IndexBuilding3LeftMiddle)

	// Middle section
	for x := 1; x < width-1; x++ {
		tileMap.SetTile(x, 0, assets.IndexBuilding3RoofTop)
		tileMap.SetTile(x, 1, assets.IndexBuilding3RoofMiddle)

		for y := 2; y < height; y += 2 {
			tileMap.SetTile(x, y, assets.IndexBuilding3WindowMiddle)
			tileMap.SetTile(x, y+1, assets.IndexBuilding3WhiteMiddle)
		}
	}

	// Right section
	tileMap.SetTile(width-1, 0, assets.IndexBuilding3RoofTop)
	tileMap.SetTile(width-1, 1, assets.IndexBuilding3RoofRight)
	tileMap.SetCol(width-1, 2, assets.IndexBuilding3RightMiddle)

	return createLevelBlock(ent, tileMap, width, height)
}

type LevelBlockable interface {
	ecs.BasicFace
	components.TransformFace
}

func populateLevelBlock(rand *rand.Rand, w *ecs.World, lb LevelBlockable) {
	if lb.GetTransformComponent().Postion.X < 50 {
		return
	}

	trans := lb.GetTransformComponent()
	biscuit := entity.CreateBiscuitEnemy()
	biscuit.Postion.X = utility.RandRangeFloat64(rand, int(trans.Postion.X), int(trans.Postion.X+trans.Size.X))
	biscuit.Postion.Y = 30
	biscuit.Layer = enemyLayer
	w.AddEntity(biscuit)
}

func createRandomLevelBlock(rand *rand.Rand, basic ecs.BasicEntity) *LevelBlock {
	val := rand.Float64()
	switch {
	case val <= 0.25:
		return createBuilding0(rand, basic)
	case val <= 0.50:
		return createBuilding1(rand, basic)
	case val <= 0.75:
		return createBuilding2(rand, basic)
	case val <= 1:
		return createBuilding3(rand, basic)
	}

	panic("random number bug")
}

func generateCityBuildings(info *Info, w *ecs.World) {
	mainGameInfo := info.MainGameInfo
	x := mainGameInfo.Level.StartX
	for x < mainGameInfo.Level.Width {
		ent := ecs.NewBasic()
		levelBlock := createRandomLevelBlock(info.Rand, ent)
		trans := levelBlock.GetTransformComponent()
		levelBlock.GetTransformComponent().Postion.Y = mainGameInfo.Level.Height - trans.Size.Y
		levelBlock.GetTransformComponent().Postion.X = x
		x += levelBlock.Size.X + float64(utility.RandRange(info.Rand, minSpaceBetweenBuildings, minSpaceBetweenBuildings+20*scaleMultiplier))
		w.AddEntity(levelBlock)
		populateLevelBlock(info.Rand, w, levelBlock)
	}
	mainGameInfo.Level.StartX = x
}
