package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
)

type ConstantSpeedable interface {
	ecs.BasicFace
	components.VelocityFace
	components.ConstantSpeedFace
}

type ConstantSpeedSystem struct {
	ents map[uint64]ConstantSpeedable
}

func CreateConstantSpeedSystem() *ConstantSpeedSystem {
	return &ConstantSpeedSystem{}
}

func (s *ConstantSpeedSystem) Priority() int {
	return int(systemPriorityConstantSpeedSystem)
}

func (s *ConstantSpeedSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]ConstantSpeedable)
}

func (s *ConstantSpeedSystem) Update(_ float32) {
	for _, ent := range s.ents {
		speedCom := ent.GetConstantSpeedComponent()
		velCom := ent.GetVelocityComponent()

		velCom.Vel = velCom.Vel.Add(speedCom.Speed)
	}
}

func (s *ConstantSpeedSystem) Add(r ConstantSpeedable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *ConstantSpeedSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *ConstantSpeedSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(ConstantSpeedable))
}
