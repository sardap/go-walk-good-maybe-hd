package game

import (
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type ScrollingSystem struct {
	ents        map[uint64]Scrollable
	scrollSpeed math.Vector2
}

func CreateScrollingSystem(scrollSpeed math.Vector2) *ScrollingSystem {
	return &ScrollingSystem{
		scrollSpeed: scrollSpeed,
	}
}

func (s *ScrollingSystem) Priority() int {
	return int(systemPriorityScrollingSystem)
}

func (s *ScrollingSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Scrollable)
}

func (s *ScrollingSystem) Update(dt float32) {
	for _, ent := range s.ents {
		vel := ent.GetVelocityComponent().Vel
		vel = vel.Add(s.scrollSpeed)
		ent.GetVelocityComponent().Vel = vel
	}
}

func (s *ScrollingSystem) Add(r Scrollable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *ScrollingSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type Scrollable interface {
	ecs.BasicFace
	components.VelocityFace
	components.ScrollableFace
}

func (s *ScrollingSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Scrollable))
}
