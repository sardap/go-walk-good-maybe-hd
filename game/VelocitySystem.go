package game

import (
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type VelocitySystem struct {
	ents map[uint64]Velocityable
}

func CreateVelocitySystem() *VelocitySystem {
	return &VelocitySystem{}
}

func (s *VelocitySystem) Priority() int {
	return int(systemPriorityVelocitySystem)
}

func (s *VelocitySystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Velocityable)
}

func (s *VelocitySystem) Update(dt float32) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()
		vel := ent.GetVelocityComponent().Vel

		trans.Postion = trans.Postion.Add(vel.Mul(float64(dt)))

		ent.GetVelocityComponent().Vel = math.Vector2{}
	}
}

func (s *VelocitySystem) Add(r Velocityable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *VelocitySystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type Velocityable interface {
	ecs.BasicFace
	components.TransformFace
	components.VelocityFace
}

func (s *VelocitySystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Velocityable))
}
