package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type LevelGenerated struct {
	ecs.BasicEntity
	components.TransformComponent
	components.VelocityComponent
	components.ImageComponent
	components.CollisionComponent
	components.ScrollableComponent
}

type LevelGenerationSystem struct {
	generated map[uint64]*LevelGenerated
}

func CreateLevelGenerationSystem() *LevelGenerationSystem {
	return &LevelGenerationSystem{}
}

func (s *LevelGenerationSystem) New(world *ecs.World) {
	s.generated = make(map[uint64]*LevelGenerated)
}

func (s *LevelGenerationSystem) Update(dt float32) {
}

func (s *LevelGenerationSystem) Add(r LevelGenerationable) {
}

func (s *LevelGenerationSystem) Remove(e ecs.BasicEntity) {
}

type LevelGenerationable interface {
	ecs.BasicFace
}

func (s *LevelGenerationSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(LevelGenerationable))
}
