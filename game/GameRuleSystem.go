package game

import (
	"fmt"

	"github.com/SolarLune/resolv"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type GameRuleSystem struct {
	ents  map[uint64]interface{}
	world *ecs.World
	space *resolv.Space
}

func CreateGameRuleSystem(space *resolv.Space) *GameRuleSystem {
	return &GameRuleSystem{
		space: space,
	}
}

func (s *GameRuleSystem) Priority() int {
	return int(systemPriorityGameRuleSystem)
}

func (s *GameRuleSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]interface{})
	s.world = world
	mainGameInfo.state = gameStateStarting
}

func (s *GameRuleSystem) updatePlayer(dt float32, player *entity.Player) {
	switch mainGameInfo.state {
	case gameStateStarting:
		if player.TransformComponent.Postion.X > 50 {
			mainGameInfo.state = gameStateScrolling
			mainGameInfo.scrollingSpeed.X = xStartScrollSpeed
		}
	case gameStateScrolling:
		break
	}

	playerCom := player.GetMainGamePlayerComponent()

	move := player.GetMovementComponent()
	vel := player.GetVelocityComponent().Vel

	fmt.Printf("Player update\n")

	switch playerCom.State {
	case components.MainGamePlayerStateGroundIdling:
		fmt.Printf("state ground\n")
		if move.MoveUp {
			vel.Y -= player.JumpPower
			player.State = components.MainGamePlayerStateJumping
		} else if move.MoveLeft || move.MoveRight {
			player.State = components.MainGamePlayerStateGroundMoving

			img, _ := assets.LoadEbitenImage(assets.ImageWhaleWalkTileSet)
			player.TileMap.TilesImg = img
			player.TileMap.SetTile(0, 0, 0)
		}
	case components.MainGamePlayerStateGroundMoving:
		if !move.MoveLeft && !move.MoveRight {
			player.State = components.MainGamePlayerStateGroundIdling

			img, _ := assets.LoadEbitenImage(assets.ImageWhaleIdleTileSet)
			player.TileMap.TilesImg = img
			player.TileMap.SetTile(0, 0, 0)
		}

		if move.MoveUp {
			vel.Y -= player.JumpPower
			player.State = components.MainGamePlayerStateJumping
		}

	case components.MainGamePlayerStateJumping:
		fmt.Printf("state Jumping\n")
		player.State = components.MainGamePlayerStateFalling

		img, _ := assets.LoadEbitenImage(assets.ImageWhaleAirTileSet)
		player.TileMap.TilesImg = img
		player.TileMap.SetTile(0, 0, 0)
	case components.MainGamePlayerStateFalling:
		fmt.Printf("state Falling\n")
		if player.Collisions.CollidingWith(entity.TagGround) {
			player.State = components.MainGamePlayerStateGroundIdling

			img, _ := assets.LoadEbitenImage(assets.ImageWhaleIdleTileSet)
			player.TileMap.TilesImg = img
			player.TileMap.SetTile(0, 0, 0)
		}
	default:
		panic("Unimplemented")
	}

	if move.MoveLeft {
		vel.X = -playerCom.Speed
	} else if move.MoveRight {
		vel.X = playerCom.Speed
	}

	// Must reset no matter what
	move.MoveLeft = false
	move.MoveRight = false
	move.MoveUp = false

	player.GetVelocityComponent().Vel = vel
}

type Collideable interface {
	ecs.BasicFace
	components.CollisionFace
	components.IdentityFace
}

type Wrapable interface {
	ecs.BasicFace
	components.TransformFace
	components.WrapFace
}

type Scrollable interface {
	ecs.BasicFace
	components.TransformFace
	components.ScrollableFace
}

type Gravityable interface {
	ecs.BasicFace
	components.VelocityFace
	components.GravityFace
	components.IdentityFace
	components.CollisionFace
}

func (s *GameRuleSystem) Update(dt float32) {
	ground := s.space.FilterByTags(entity.TagGround)

	for _, ent := range s.ents {
		if wrapable, ok := ent.(Wrapable); ok {
			trans := wrapable.GetTransformComponent()
			if trans.Postion.X < -wrapable.GetWrapComponent().Threshold {
				trans.Postion.X = wrapable.GetWrapComponent().Threshold
			}
		}

		if scrollable, ok := ent.(Scrollable); ok {
			trans := scrollable.GetTransformComponent().Postion
			trans = trans.Add(mainGameInfo.scrollingSpeed.Mul(float64(dt)))
			scrollable.GetTransformComponent().Postion = trans
		}

		if building, ok := ent.(*Building); ok {
			trans := building.GetTransformComponent()
			if trans.Postion.X+trans.Size.X < 0 {
				defer s.world.RemoveEntity(building.BasicEntity)
				fmt.Printf("Removing %d\n", building.ID())
			}
		}

		if gravityable, ok := ent.(Gravityable); ok {
			vel := gravityable.GetVelocityComponent()
			colCom := gravityable.GetCollisionComponent()

			collision := ground.Resolve(colCom.CollisionShape, 0, mainGameInfo.gravity*float64(dt))
			if !collision.Colliding() {
				vel.Vel = vel.Vel.Add(math.Vector2{Y: mainGameInfo.gravity})
			}
		}

		if ent, ok := ent.(*entity.Player); ok {
			s.updatePlayer(dt, ent)
		}
	}

	mainGameInfo.level.StartX += mainGameInfo.scrollingSpeed.X * float64(dt)
	generateBuildings(s.world)
}

func (s *GameRuleSystem) Add(r GameRuleable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *GameRuleSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type GameRuleable interface {
	ecs.BasicFace
}

func (s *GameRuleSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(GameRuleable))
}
