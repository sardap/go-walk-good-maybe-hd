package game

import (
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
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
	}
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
