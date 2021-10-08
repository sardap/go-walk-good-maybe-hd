package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
)

type EnemyBiscuitable interface {
	ecs.BasicFace
	components.TransformFace
	components.BiscuitEnemyFace
	components.CollisionFace
	components.MovementFace
	components.VelocityFace
}

type EnemyBiscuitSystem struct {
	ents  map[uint64]EnemyBiscuitable
	world *ecs.World
	space *resolv.Space
}

func CreateEnemyBiscuitSystem(space *resolv.Space) *EnemyBiscuitSystem {
	return &EnemyBiscuitSystem{
		space: space,
	}
}

func (s *EnemyBiscuitSystem) Priority() int {
	return int(systemPriorityEnemyBiscuitSystem)
}

func (s *EnemyBiscuitSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]EnemyBiscuitable)
	s.world = world
}

func (s *EnemyBiscuitSystem) Update(dt float32) {
	for _, biscuit := range s.ents {
		velCom := biscuit.GetVelocityComponent()
		colCom := biscuit.GetCollisionComponent()
		biscuitCom := biscuit.GetBiscuitEnemyComponent()

		col := s.space.Resolve(colCom.CollisionShape, biscuitCom.Speed.X+10*float64(dt), 50)

		if col.Colliding() {
			velCom.Vel.X += biscuitCom.Speed.X
		} else {
			biscuitCom.Speed.X = -biscuitCom.Speed.X
		}
	}
}

func (s *EnemyBiscuitSystem) Add(r EnemyBiscuitable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *EnemyBiscuitSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *EnemyBiscuitSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(EnemyBiscuitable))
}
