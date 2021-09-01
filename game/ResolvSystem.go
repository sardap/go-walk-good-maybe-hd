package game

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type Resolvable interface {
	ecs.BasicFace
	components.TransformFace
	components.IdentityFace
	components.CollisionFace
}

type ResolvSystem struct {
	ents           map[uint64]Resolvable
	space          *resolv.Space
	overlay        *ebiten.Image
	OverlayEnabled bool
}

func CreateResolvSystem(space *resolv.Space) *ResolvSystem {
	return &ResolvSystem{
		space: space,
	}
}

func (s *ResolvSystem) Priority() int {
	return int(systemPriorityResolvSystem)
}

func (s *ResolvSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Resolvable)
	s.overlay = ebiten.NewImage(windowWidth, gameHeight)
}

func (s *ResolvSystem) Update(dt float32) {
	if inpututil.IsKeyJustReleased(ebiten.KeyO) {
		s.OverlayEnabled = !s.OverlayEnabled
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

func (s *ResolvSystem) Render(cmds *RenderCmds) {
	s.overlay.Fill(color.RGBA{0, 0, 0, 0})

	if !s.OverlayEnabled {
		return
	}

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

func (s *ResolvSystem) Add(r Resolvable) {
	s.ents[r.GetBasicEntity().ID()] = r

	if r.GetCollisionComponent().CollisionShape != nil {
		if !s.space.Contains(r.GetCollisionComponent().CollisionShape) {
			s.space.Add(r.GetCollisionComponent().CollisionShape)
		}
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

func (s *ResolvSystem) Remove(e ecs.BasicEntity) {
	if ent, ok := s.ents[e.ID()]; ok {
		s.space.Remove(ent.GetCollisionComponent().CollisionShape)
	}

	delete(s.ents, e.ID())
}

func (s *ResolvSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Resolvable))
}
