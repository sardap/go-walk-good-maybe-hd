package game

import (
	"github.com/SolarLune/resolv"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
)

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

		colShape := colCom.CollisionShape

		colShape.X -= 5
		colShape.Y -= 5
		colShape.W += 10
		colShape.H += 10

		if collision := s.space.Collision(colShape); collision != nil && collision.Colliding() {
			colCom.Collisions = append(colCom.Collisions, &components.CollisionEvent{
				Tags: collision.ShapeB.GetTags(),
			})
		}

		colShape.X += 5
		colShape.Y += 5
		colShape.W -= 10
		colShape.H -= 10

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

	if r.GetCollisionComponent().CollisionShape != nil {
		return
	}

	trans := r.GetTransformComponent()
	ident := r.GetIdentityComponent()

	rectangle := resolv.NewRectangle(
		trans.Postion.X, trans.Postion.Y,
		trans.Size.X, trans.Size.Y,
	)

	rectangle.AddTags(ident.Tags...)

	s.space.Add(rectangle)

	r.GetCollisionComponent().CollisionShape = rectangle
}

func (s *VelocitySystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type Velocityable interface {
	ecs.BasicFace
	components.TransformFace
	components.IdentityFace
	components.CollisionFace
	components.VelocityFace
}

func (s *VelocitySystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Velocityable))
}
