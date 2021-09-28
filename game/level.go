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

type spawnProbability struct {
	genFunc     func(*rand.Rand, *ecs.World, LevelBlockable)
	probability float64
}

func createBiscuitEnemy(rand *rand.Rand, w *ecs.World, lb LevelBlockable) {
	lbTrans := lb.GetTransformComponent()
	biscuit := entity.CreateBiscuitEnemy()
	biscuit.Postion.X = utility.RandRangeFloat64(rand, int(lbTrans.Postion.X), int(lbTrans.Postion.X+lbTrans.Size.X-biscuit.Size.X))
	biscuit.Postion.Y = lbTrans.Postion.Y - (biscuit.Size.Y * 1.5)
	biscuit.Layer = ImageLayerObjects
	w.AddEntity(biscuit)
}

func createUfoBiscuitEnemy(rand *rand.Rand, w *ecs.World, lb LevelBlockable) {
	lbTrans := lb.GetTransformComponent()
	ufo := entity.CreateUfoBiscuitEnemy()
	ufo.Postion.X = utility.RandRangeFloat64(rand, int(lbTrans.Postion.X), int(lbTrans.Postion.X+lbTrans.Size.X-ufo.Size.X))
	ufo.Postion.Y = lbTrans.Postion.Y - (ufo.Size.Y * 2.5)
	ufo.Layer = ImageLayerObjects
	w.AddEntity(ufo)
}

func createJumpToken(rand *rand.Rand, w *ecs.World, lb LevelBlockable) {
	lbTrans := lb.GetTransformComponent()
	token := entity.CreateJumpUpToken()
	token.Postion.X = utility.RandRangeFloat64(
		rand,
		int(lbTrans.Postion.X),
		int(lbTrans.Postion.X+lbTrans.Size.X-token.TransformComponent.Size.X),
	)
	token.Postion.Y = lbTrans.Postion.Y - token.TransformComponent.Size.Y
	token.Layer = ImageLayerObjects
	w.AddEntity(token)
}

func createSpeedToken(rand *rand.Rand, w *ecs.World, lb LevelBlockable) {
	lbTrans := lb.GetTransformComponent()
	token := entity.CreateSpeedUpToken()
	token.Postion.X = utility.RandRangeFloat64(
		rand,
		int(lbTrans.Postion.X),
		int(lbTrans.Postion.X+lbTrans.Size.X-token.TransformComponent.Size.X),
	)
	token.Postion.Y = lbTrans.Postion.Y - token.TransformComponent.Size.Y
	token.Layer = ImageLayerObjects
	w.AddEntity(token)
}

type LevelBlock struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.VelocityComponent
	*components.TileImageComponent
	*components.CollisionComponent
	*components.ScrollableComponent
	*components.IdentityComponent
	probabilities []spawnProbability
}

func (l *LevelBlock) GetSpawnProbabilities() []spawnProbability {
	return l.probabilities
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
			Layer:   ImageLayerbuildingForground,
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{entity.TagGround},
		},
		probabilities: []spawnProbability{
			{genFunc: createBiscuitEnemy, probability: 0.5},
			{genFunc: createUfoBiscuitEnemy, probability: 0.5},
			{genFunc: createJumpToken, probability: 0.5},
			{genFunc: createSpeedToken, probability: 0.5},
		},
	}
}

func createBuilding0(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImageAsset(assets.ImageBuilding0TileSet)

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
	tileSet, _ := assets.LoadEbitenImageAsset(assets.ImageBuilding1TileSet)

	width := utility.RandRange(rand, 4, 6)
	height := utility.RandRange(rand, 3, 5)

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
	tileSet, _ := assets.LoadEbitenImageAsset(assets.ImageBuilding2TileSet)

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
	tileSet, _ := assets.LoadEbitenImageAsset(assets.ImageBuilding3TileSet)

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

func createBuilding4(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImageAsset(assets.ImageBuilding4TileSet)

	width := utility.RandRange(rand, 5, 6)

	height := utility.RandRange(rand, 6, 8)
	for height%2 != 0 {
		height++
	}
	height++

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding4TileSet.FrameWidth)
	// Left section
	tileMap.SetTile(0, 0, assets.IndexBuilding4RoofLeft)
	tileMap.SetCol(0, 1, assets.IndexBuilding4LeftMiddle)

	// Middle section
	for x := 1; x < width-1; x++ {
		tileMap.SetTile(x, 0, assets.IndexBuilding4RoofMiddle)

		for y := 1; y < height; y += 2 {
			val := rand.Float64()
			switch {
			case val <= 0.166:
				tileMap.SetTile(x, y, assets.IndexBuilding4SignYellowTop)
				tileMap.SetTile(x, y+1, assets.IndexBuilding4SignYellowBot)
			case val <= 0.333:
				tileMap.SetTile(x, y, assets.IndexBuilding4SignOrangeTop)
				tileMap.SetTile(x, y+1, assets.IndexBuilding4SignOrangeBot)
			case val <= 0.500:
				tileMap.SetTile(x, y, assets.IndexBuilding4SignGreenTop)
				tileMap.SetTile(x, y+1, assets.IndexBuilding4SignGreenBot)
			case val <= 0.666:
				tileMap.SetTile(x, y, assets.IndexBuilding4SignRed)
				tileMap.SetTile(x, y+1, assets.IndexBuilding4MiddlePlain)
			case val <= 0.833:
				tileMap.SetTile(x, y, assets.IndexBuilding4SignBlue)
				tileMap.SetTile(x, y+1, assets.IndexBuilding4MiddlePlain)
			case val <= 1:
				tileMap.SetTile(x, y, assets.IndexBuilding4MiddlePlain)
				tileMap.SetTile(x, y+1, assets.IndexBuilding4MiddlePlain)
			}
		}
	}

	// Right section
	tileMap.SetTile(width-1, 0, assets.IndexBuilding4RoofRight)
	tileMap.SetCol(width-1, 1, assets.IndexBuilding4RightMiddle)

	return createLevelBlock(ent, tileMap, width, height)
}

func createBuilding5(rand *rand.Rand, ent ecs.BasicEntity) *LevelBlock {
	tileSet, _ := assets.LoadEbitenImageAsset(assets.ImageBuilding5TileSet)

	width := utility.RandRange(rand, 4, 6)

	height := utility.RandRange(rand, 4, 6)
	for height%2 != 0 {
		height++
	}
	height++

	tileMap := components.CreateTileMap(width, height, tileSet, assets.ImageBuilding5TileSet.FrameWidth)
	// Left section
	tileMap.SetTile(0, 0, assets.IndexBuilding5RoofLeft)
	tileMap.SetCol(0, 1, assets.IndexBuilding5MiddleLeft)

	// Middle section
	for x := 1; x < width; x++ {
		tileMap.SetTile(x, 0, assets.IndexBuilding5RoofMiddle)
		tileMap.SetCol(x, 1, assets.IndexBuilding5MiddleWindow)
	}

	// Right section
	tileMap.SetTile(width-1, 0, assets.IndexBuilding5RoofRight)
	tileMap.SetCol(width-1, 1, assets.IndexBuilding5MiddleRight)

	return createLevelBlock(ent, tileMap, width, height)
}

type LevelBlockable interface {
	ecs.BasicFace
	components.TransformFace
	GetSpawnProbabilities() []spawnProbability
}

func populateLevelBlock(rand *rand.Rand, w *ecs.World, lb LevelBlockable) {
	// Don't spawn on player spawn
	if lb.GetTransformComponent().Postion.X < 50 {
		return
	}

	for _, prob := range lb.GetSpawnProbabilities() {
		if prob.probability < rand.Float64() {
			prob.genFunc(rand, w, lb)
			return
		}
	}
}

type genFunc func(*rand.Rand, ecs.BasicEntity) *LevelBlock

func createRandomLevelBlock(rand *rand.Rand, basic ecs.BasicEntity) *LevelBlock {

	funcs := []genFunc{createBuilding0, createBuilding1, createBuilding2, createBuilding3, createBuilding4, createBuilding5}

	val := rand.Float64()
	// Probably don't need to iterrate over the whole thing
	for i, genFunc := range funcs {
		if val < (1.0/float64(len(funcs)))*float64(i+1.0) || i == len(funcs)-1 {
			return genFunc(rand, basic)
		}
	}

	panic("random number bug")
}
