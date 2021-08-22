package game

import (
	"fmt"

	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
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

type Moveable interface {
	ecs.BasicFace
	components.MovementFace
	components.VelocityFace
}

type Wrapable interface {
	ecs.BasicFace
	components.TransformFace
	components.WrapFace
}

type Scrollable interface {
	ecs.BasicFace
	components.VelocityFace
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
		if ent, ok := ent.(*entity.Player); ok {
			s.updatePlayer(dt, ent)
		}

		if moveable, ok := ent.(Moveable); ok {
			move := moveable.GetMovementComponent()
			vel := moveable.GetVelocityComponent().Vel
			if move.MoveLeft {
				vel.X -= move.Speed
				move.MoveLeft = false
			}
			if move.MoveRight {
				vel.X += move.Speed
				move.MoveRight = false
			}

			if move.MoveDown {
				vel.Y += move.Speed
				move.MoveDown = false
			}
			if move.MoveUp {
				vel.Y -= move.Speed
				move.MoveUp = false
			}

			moveable.GetVelocityComponent().Vel = vel
		}

		if wrapable, ok := ent.(Wrapable); ok {
			trans := wrapable.GetTransformComponent()
			if trans.Postion.X < -wrapable.GetWrapComponent().Threshold {
				trans.Postion.X = wrapable.GetWrapComponent().Threshold
			}
		}

		if scrollable, ok := ent.(Scrollable); ok {
			vel := scrollable.GetVelocityComponent().Vel
			vel = vel.Add(mainGameInfo.scrollingSpeed)
			scrollable.GetVelocityComponent().Vel = vel
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

			if s.CollidedWith(gravityable, entity.TagGround) {
				vel.Acc = math.Vector2{}
				continue
			}

			vel.Acc = vel.Acc.Add(math.Vector2{Y: mainGameInfo.gravity}.Mul(float64(dt)))
			vel.Acc = math.ClampVec2(
				vel.Acc,
				math.Vector2{
					X: -maxAccelerationX,
					Y: -maxAccelerationY,
				},
				math.Vector2{
					X: maxAccelerationX,
					Y: maxAccelerationY,
				},
			)
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
