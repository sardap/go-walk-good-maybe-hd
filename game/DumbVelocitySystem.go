package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type DumbVelocityable interface {
	ecs.BasicFace
	components.TransformFace
	components.VelocityFace
}

type ExDumbVelocityable interface {
	components.CollisionFace
}

type DumbVelocitySystem struct {
	ents map[uint64]DumbVelocityable
}

func CreateDumbVelocitySystem() *DumbVelocitySystem {
	return &DumbVelocitySystem{}
}

func (s *DumbVelocitySystem) Priority() int {
	return int(systemPriorityDumbVelocitySystem)
}

func (s *DumbVelocitySystem) New(world *ecs.World) {
	s.ents = make(map[uint64]DumbVelocityable)
}

func (s *DumbVelocitySystem) Update(dt float32) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()
		vel := ent.GetVelocityComponent().Vel

		trans.Postion = trans.Postion.Add(vel.Mul(float64(dt)))

		ent.GetVelocityComponent().Vel = math.Vector2{}
	}
}

func (s *DumbVelocitySystem) Add(r DumbVelocityable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *DumbVelocitySystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *DumbVelocitySystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(DumbVelocityable))
}
