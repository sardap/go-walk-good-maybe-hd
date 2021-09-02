package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type Velocityable interface {
	ecs.BasicFace
	components.TransformFace
	components.IdentityFace
	components.CollisionFace
	components.VelocityFace
}

type VelocitySystem struct {
	ents  map[uint64]Velocityable
	space *resolv.Space
}

func CreateVelocitySystem(space *resolv.Space) *VelocitySystem {
	return &VelocitySystem{
		space: space,
	}
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
		colCom := ent.GetCollisionComponent()
		vel := ent.GetVelocityComponent().Vel

		vel = vel.Mul(float64(dt))

		colCom.Collisions = nil

		ground := s.space.FilterByTags(entity.TagGround)

		collision := ground.Resolve(colCom.CollisionShape, vel.X, 0)
		if collision.Colliding() {
			trans.Postion.X += collision.ResolveX
		} else {
			trans.Postion.X += vel.X
		}

		collision = ground.Resolve(colCom.CollisionShape, 0, vel.Y)
		if collision.Colliding() {
			trans.Postion.Y += collision.ResolveY
		} else {
			trans.Postion.Y += vel.Y
		}

		colCom.CollisionShape.SetXY(trans.Postion.X, trans.Postion.Y)

		ent.GetVelocityComponent().Vel = math.Vector2{}
	}
}

func (s *VelocitySystem) Add(r Velocityable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *VelocitySystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *VelocitySystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Velocityable))
}
