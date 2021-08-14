package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type GameRuleSystem struct {
	ents map[uint64]GameRuleable
}

func CreateGameRuleSystem() *GameRuleSystem {
	return &GameRuleSystem{}
}

func (s *GameRuleSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]GameRuleable)
}

func (s *GameRuleSystem) Update(dt float32) {
	for _, ent := range s.ents {
		if moveInter, ok := ent.(components.MovementFace); ok {
			move := moveInter.GetMovementComponent()
			trans := ent.(components.TransformFace).GetTransformComponent()

			var velX float64
			if move.MoveLeft {
				velX -= 10
				move.MoveLeft = false
			}
			if move.MoveRight {
				velX += 10
				move.MoveRight = false
			}

			var velY float64
			if move.MoveDown {
				velY += 10
				move.MoveDown = false
			}
			if move.MoveUp {
				velY -= 10
				move.MoveUp = false
			}

			trans.GeoM.Translate(velX, velY)
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
