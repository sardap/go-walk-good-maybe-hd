package game

import (
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type gameState int

const (
	gameStateStarting gameState = iota
	gameStateScrolling
)

type GameRuleSystem struct {
	ents  map[uint64]interface{}
	world *ecs.World
	state gameState
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
	s.state = gameStateStarting
}

func (s *GameRuleSystem) updatePlayer(dt float32, player *entity.Player) {
	switch s.state {
	case gameStateStarting:
		if player.MoveRight {
			var scrollable *Scrollable
			s.world.AddSystemInterface(CreateScrollingSystem(math.Vector2{X: -2.5, Y: 0}), scrollable, nil)
			s.state = gameStateScrolling
		}
	}
}

type Moveable interface {
	ecs.BasicFace
	components.MovementFace
	components.VelocityFace
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
				vel.X -= 10
				move.MoveLeft = false
			}
			if move.MoveRight {
				vel.X += 10
				move.MoveRight = false
			}

			if move.MoveDown {
				vel.Y += 10
				move.MoveDown = false
			}
			if move.MoveUp {
				vel.Y -= 10
				move.MoveUp = false
			}

			moveable.GetVelocityComponent().Vel = vel
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
