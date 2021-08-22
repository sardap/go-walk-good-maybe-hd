package game

import (
	"fmt"

	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"

	gomath "math"
)

type GameRuleSystem struct {
	ents  map[uint64]interface{}
	world *ecs.World
}

func CreateGameRuleSystem() *GameRuleSystem {
	return &GameRuleSystem{}
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
	if move.MoveLeft {
		vel.X = -playerCom.Speed
	} else if move.MoveRight {
		vel.X = playerCom.Speed
	} else {
		vel.X = 0
	}

	fmt.Printf("Player update\n")

	switch playerCom.State {
	case components.MainGamePlayerStateGround:
		fmt.Printf("state ground\n")
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
		if s.CollidedWith(player, entity.TagGround) {
			player.State = components.MainGamePlayerStateGround

			img, _ := assets.LoadEbitenImage(assets.ImageWhaleIdleTileSet)
			player.TileMap.TilesImg = img
			player.TileMap.SetTile(0, 0, 0)
		}
	default:
		panic("Unimplemented")
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

func (s *GameRuleSystem) CollidedWith(a Collideable, tag string) bool {
	if !a.GetCollisionComponent().Active {
		return false
	}

	for _, event := range a.GetCollisionComponent().Collisions {
		otherInter, ok := s.ents[event.ID]
		if !ok {
			continue
		}

		if other, ok := otherInter.(Collideable); ok && other.GetIdentityComponent().HasTag(tag) {
			return true
		}
	}

	return false
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
			func() {
				vel := gravityable.GetVelocityComponent()

				if s.CollidedWith(gravityable, entity.TagGround) {
					vel.Vel.Y = utility.ClampFloat64(vel.Vel.Y, -gomath.MaxFloat64, 0)
					return
				}

				vel.Vel = vel.Vel.Add(math.Vector2{Y: mainGameInfo.gravity}.Mul(float64(dt)))
				vel.Vel = utility.ClampVec2(
					vel.Vel,
					math.Vector2{
						X: -maxAccelerationX,
						Y: -maxAccelerationY,
					},
					math.Vector2{
						X: maxAccelerationX,
						Y: maxAccelerationY,
					},
				)
			}()
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
