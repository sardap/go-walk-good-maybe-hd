package game

import (
	"container/heap"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type CollisionSystem struct {
	ents    map[uint64]Collisionable
	overlay *ebiten.Image
}

func CreateCollisionSystem() *CollisionSystem {
	return &CollisionSystem{}
}

func (s *CollisionSystem) Priority() int {
	return int(systemPriorityCollisionSystem)
}

func (s *CollisionSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Collisionable)
	s.overlay = ebiten.NewImage(gameWidth, gameHeight)
}

func getRect(ent Collisionable) (x1, x2, y1, y2 float64) {
	trans := ent.GetTransformComponent()

	var imgW, imgH int
	if ent.GetImageComponent().SubRect != nil {
		imgW, imgH = ent.GetImageComponent().SubRect.Dx(), ent.GetImageComponent().SubRect.Dy()
	} else {
		imgW, imgH = ent.GetImageComponent().Image.Size()
	}

	x1, w := trans.Element(0, 2), trans.Element(1, 1)*float64(imgW)
	y1, h := trans.Element(1, 2), trans.Element(0, 0)*float64(imgH)

	return x1, x1 + w, y1, y1 + h

}

func (s *CollisionSystem) Update(dt float32) {
	for _, entA := range s.ents {

		if !entA.GetCollisionComponent().Active {
			continue
		}

		entCol := entA.GetCollisionComponent()
		entCol.Collisions = nil

		aX1, aX2, aY1, aY2 := getRect(entA)

		for _, entB := range s.ents {
			if !entB.GetCollisionComponent().Active ||
				entA.GetBasicEntity().ID() == entB.GetBasicEntity().ID() {
				continue
			}

			bX1, bX2, bY1, bY2 := getRect(entA)

			if aX1 < bX2 && aX2 > bX1 && aY1 < bY2 && aY2 > bY1 {
				entCol.Collisions = append(entCol.Collisions, &components.CollisionEvent{
					ID: entB.GetBasicEntity().ID(),
				})
			}
		}
	}
}

func (s *CollisionSystem) Render(cmds *RenderCmds) {
	s.overlay.Fill(color.RGBA{0, 0, 0, 0})

	for _, ent := range s.ents {
		if !ent.GetCollisionComponent().Active {
			continue
		}

		x1, x2, y1, y2 := getRect(ent)
		w := x2 - x1
		h := y2 - y1

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
	heap.Push(cmds, &RenderCmd{
		Image:   s.overlay,
		Options: op,
		Layer:   debugImageLayer,
	})
}

func (s *CollisionSystem) Add(r Collisionable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *CollisionSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type Collisionable interface {
	ecs.BasicFace
	components.TransformFace
	components.CollisionFace
	components.ImageFace
}

func (s *CollisionSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Collisionable))
}
