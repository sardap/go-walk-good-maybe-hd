package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type DestoryBoundable interface {
	ecs.BasicFace
	components.TransformFace
	components.DestoryBoundFace
}

type DestoryBoundSystem struct {
	ents  map[uint64]DestoryBoundable
	world *ecs.World
}

func CreateDestoryBoundSystem() *DestoryBoundSystem {
	return &DestoryBoundSystem{}
}

func (s *DestoryBoundSystem) Priority() int {
	return int(systemPriorityDestoryBoundSystem)
}

func (s *DestoryBoundSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]DestoryBoundable)
	s.world = world
}

func (s *DestoryBoundSystem) Update(_ float32) {
	for _, ent := range s.ents {
		transCom := ent.GetTransformComponent()
		destroyCom := ent.GetDestoryBoundComponent()

		if transCom.Postion.X > destroyCom.Max.X ||
			transCom.Postion.X+transCom.Size.X < destroyCom.Min.X ||
			transCom.Postion.Y > destroyCom.Max.Y ||
			transCom.Postion.Y+transCom.Size.Y < destroyCom.Min.Y {
			defer s.world.RemoveEntity(*ent.GetBasicEntity())
		}
	}
}

func (s *DestoryBoundSystem) Add(r DestoryBoundable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *DestoryBoundSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *DestoryBoundSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(DestoryBoundable))
}

func (s *DestoryBoundSystem) ContainsEnt(id uint64) bool {
	_, ok := s.ents[id]
	return ok
}
