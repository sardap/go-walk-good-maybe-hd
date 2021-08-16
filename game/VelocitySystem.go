package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type VelocitySystem struct {
	ents map[uint64]Velocityable
}

func CreateVelocitySystem() *VelocitySystem {
	return &VelocitySystem{}
}

func (s *VelocitySystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Velocityable)
}

func (s *VelocitySystem) Update(dt float32) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()
		vel := ent.GetVelocityComponent()

		trans.GeoM.Translate(float64(vel.Vel.X), float64(vel.Vel.Y))
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
