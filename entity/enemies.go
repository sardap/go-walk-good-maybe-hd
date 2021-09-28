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
	*components.BiscuitEnemyComponent
	*components.CollisionComponent
	*components.DamageComponent
	*components.LifeComponent
	*components.IdentityComponent
	*components.MovementComponent
	*components.GravityComponent
	*components.TileImageComponent
	*components.ScrollableComponent
	*components.VelocityComponent
}

func CreateBiscuitEnemy() *BiscuitEnemy {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageBiscuitEnemyIdleTileSet)

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
		BiscuitEnemyComponent: &components.BiscuitEnemyComponent{
			Speed: math.Vector2{
				X: 150,
				Y: 0,
			},
		},
		CollisionComponent: &components.CollisionComponent{
			Active: true,
		},
		DamageComponent: &components.DamageComponent{
			BaseDamage: 100,
		},
		LifeComponent: &components.LifeComponent{
			HP: 100,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{TagEnemy},
		},
		MovementComponent: components.CreateMovementComponent(),
		GravityComponent:  &components.GravityComponent{},
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

func CreateBiscuitEnemyDeath() *SingleScrollableAnime {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageBiscuitEnemyDeathTileSet)

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

type UfoBiscuitEnemy struct {
	ecs.BasicEntity
	*components.TransformComponent
	*components.AnimeComponent
	*components.CollisionComponent
	*components.DamageComponent
	*components.LifeComponent
	*components.IdentityComponent
	*components.TileImageComponent
	*components.ScrollableComponent
	*components.UfoBiscuitEnemyComponent
	*components.VelocityComponent
}

func CreateUfoBiscuitEnemy() *UfoBiscuitEnemy {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageBiscutUFOIdleTileSet)

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
		DamageComponent: &components.DamageComponent{
			BaseDamage: 100,
		},
		LifeComponent: &components.LifeComponent{
			HP: 100,
		},
		IdentityComponent: &components.IdentityComponent{
			Tags: []int{TagEnemy, TagUfo},
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

func CreateUfoBiscuitEnemyDeath() *SingleScrollableAnime {
	img, _ := assets.LoadEbitenImageAsset(assets.ImageBiscutUfoDeathTileSet)

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
