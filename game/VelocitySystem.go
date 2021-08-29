package game

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	ents    map[uint64]Velocityable
	space   *resolv.Space
	overlay *ebiten.Image
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
	s.overlay = ebiten.NewImage(gameWidth, gameHeight)
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

	for _, ent := range s.ents {
		colCom := ent.GetCollisionComponent()

		colShape := colCom.CollisionShape

		colShape.X -= 2.5
		colShape.Y -= 2.5
		colShape.W += 5
		colShape.H += 5

		if collision := s.space.Collision(colShape); collision != nil && collision.Colliding() {
			colCom.Collisions = append(colCom.Collisions, &components.CollisionEvent{
				Tags: collision.ShapeB.GetTags(),
			})
		}

		colShape.X += 2.5
		colShape.Y += 2.5
		colShape.W -= 5
		colShape.H -= 5
	}
}

func (s *VelocitySystem) Render(cmds *RenderCmds) {
	s.overlay.Fill(color.RGBA{0, 0, 0, 0})

	for _, ent := range s.ents {
		if !ent.GetCollisionComponent().Active {
			continue
		}

		x1 := ent.GetTransformComponent().Postion.X
		y1 := ent.GetTransformComponent().Postion.Y
		w := ent.GetTransformComponent().Size.X
		h := ent.GetTransformComponent().Size.Y

		clr := color.RGBA{255, 0, 0, 255}
		// Left Top to Right Top
		ebitenutil.DrawLine(s.overlay, x1, y1, x1+w, y1, clr)
		// Right Top to Right Bottom
		ebitenutil.DrawLine(s.overlay, x1+w, y1, x1+w, y1+h, clr)
		// Right Bottom to Left Bottom
		ebitenutil.DrawLine(s.overlay, x1+w, y1+h, x1, y1+h, clr)
		// Left Bottom to Left top
		ebitenutil.DrawLine(s.overlay, x1, y1+h, x1, y1, clr)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleMultiplier, scaleMultiplier)
	*cmds = append(*cmds, &RenderImageCmd{
		Image:   s.overlay,
		Options: op,
		Layer:   debugImageLayer,
	})
}

func (s *VelocitySystem) Add(r Velocityable) {
	s.ents[r.GetBasicEntity().ID()] = r

	if r.GetCollisionComponent().CollisionShape != nil && s.space.Contains(r.GetCollisionComponent().CollisionShape) {
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
	if ent, ok := s.ents[e.ID()]; ok {
		s.space.Remove(ent.GetCollisionComponent().CollisionShape)
	}

	delete(s.ents, e.ID())
}

func (s *VelocitySystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Velocityable))
}
