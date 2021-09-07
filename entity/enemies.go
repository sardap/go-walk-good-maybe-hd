package entity

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type BiscuitEnemy struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.AnimeComponent
	*components.CollisionComponent
	*components.BiscuitEnemyComponent
	*components.IdentityComponent
	*components.GravityComponent
	*components.TileImageComponent
	*components.ScrollableComponent
	*components.VelocityComponent
}

func CreateBiscuitEnemy() *BiscuitEnemy {
	img, _ := assets.LoadEbitenImage(assets.ImageBiscuitEnemyIdleTileSet)

	tileMap := components.CreateTileMap(1, 1, img, assets.ImageBiscuitEnemyIdleTileSet.FrameWidth)
	tileMap.SetTile(0, 0, 0)

	return &BiscuitEnemy{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageBiscuitEnemyIdleTileSet.FrameWidth),
				Y: float64(img.Bounds().Dy()),
			},
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration:  200 * time.Millisecond,
			FrameRemaining: 200 * time.Millisecond,
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		BiscuitEnemyComponent: &components.BiscuitEnemyComponent{},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{TagEnemy},
		},
		GravityComponent: &components.GravityComponent{},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		VelocityComponent: &components.VelocityComponent{},
	}
}

type SingleScrollableAnime struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.AnimeComponent
	*components.DestoryOnAnimeComponent
	*components.TileImageComponent
	*components.ScrollableComponent
	*components.VelocityComponent
}

func CreateBiscuitEnemyDeath() *SingleScrollableAnime {
	img, _ := assets.LoadEbitenImage(assets.ImageBiscuitEnemyDeathTileSet)

	tileMap := components.CreateTileMap(1, 1, img, assets.ImageBiscuitEnemyDeathTileSet.FrameWidth)
	tileMap.SetTile(0, 0, 0)

	return &SingleScrollableAnime{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageBiscuitEnemyIdleTileSet.FrameWidth),
				Y: float64(img.Bounds().Dy()),
			},
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration:  200 * time.Millisecond,
			FrameRemaining: 200 * time.Millisecond,
		},
		DestoryOnAnimeComponent: &components.DestoryOnAnimeComponent{
			CyclesTilDeath: 1,
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		VelocityComponent: &components.VelocityComponent{},
	}
}

func CreateUfoBiscuitEnemyDeath() *SingleScrollableAnime {
	img, _ := assets.LoadEbitenImage(assets.ImageBiscutUfoDeathTileSet)

	tileMap := components.CreateTileMap(1, 1, img, assets.ImageBiscutUfoDeathTileSet.FrameWidth)
	tileMap.SetTile(0, 0, 0)

	return &SingleScrollableAnime{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageBiscuitEnemyIdleTileSet.FrameWidth),
				Y: float64(img.Bounds().Dy()),
			},
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration:  100 * time.Millisecond,
			FrameRemaining: 100 * time.Millisecond,
		},
		DestoryOnAnimeComponent: &components.DestoryOnAnimeComponent{
			CyclesTilDeath: 1,
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		VelocityComponent: &components.VelocityComponent{},
	}
}

type UfoBiscuitEnemy struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.AnimeComponent
	*components.CollisionComponent
	*components.IdentityComponent
	*components.TileImageComponent
	*components.ScrollableComponent
	*components.UfoBiscuitEnemyComponent
	*components.VelocityComponent
}

func CreateUfoBiscuitEnemy() *UfoBiscuitEnemy {
	img, _ := assets.LoadEbitenImage(assets.ImageBiscutUFOIdleTileSet)

	tileMap := components.CreateTileMap(1, 1, img, assets.ImageBiscutUFOIdleTileSet.FrameWidth)
	tileMap.SetTile(0, 0, 0)

	return &UfoBiscuitEnemy{
		BasicEntity: ecs.NewBasic(),
		TransformComponent: &components.TransformComponent{
			Size: math.Vector2{
				X: float64(assets.ImageBiscutUFOIdleTileSet.FrameWidth),
				Y: float64(img.Bounds().Dy()),
			},
		},
		AnimeComponent: &components.AnimeComponent{
			FrameDuration:  200 * time.Millisecond,
			FrameRemaining: 200 * time.Millisecond,
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []string{TagEnemy, TagUfo},
		},
		TileImageComponent: &components.TileImageComponent{
			Active:  true,
			TileMap: tileMap,
		},
		ScrollableComponent: &components.ScrollableComponent{
			Modifier: 1,
		},
		UfoBiscuitEnemyComponent: &components.UfoBiscuitEnemyComponent{
			ShootTime:         1 * time.Second,
			ShootTimeRemaning: 1 * time.Second,
		},
		VelocityComponent: &components.VelocityComponent{},
	}
}
