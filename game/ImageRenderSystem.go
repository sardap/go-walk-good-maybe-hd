package game

import (
	"container/heap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type ImageRenderSystem struct {
	ents map[uint64]ImageRenderable
}

func CreateImageRenderSystem() *ImageRenderSystem {
	return &ImageRenderSystem{}
}

func (s *ImageRenderSystem) Priority() int {
	return int(systemPriorityImageRenderSystem)
}

func (s *ImageRenderSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]ImageRenderable)
}

func (s *ImageRenderSystem) Update(dt float32) {
}

func (s *ImageRenderSystem) Render(cmds *RenderCmds) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()
		imgCom := ent.GetImageComponent()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(trans.Postion.X, trans.Postion.Y)
		op.GeoM.Scale(scaleMultiplier, scaleMultiplier)

		var img *ebiten.Image
		if imgCom.SubRect != nil {
			img = imgCom.Image.SubImage(*imgCom.SubRect).(*ebiten.Image)
		} else {
			img = imgCom.Image
		}

		heap.Push(cmds, &RenderImageCmd{
			Image:   img,
			Options: op,
			Layer:   imgCom.Layer,
		})
	}
}

func (s *ImageRenderSystem) Add(r ImageRenderable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *ImageRenderSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type ImageRenderable interface {
	ecs.BasicFace
	components.TransformFace
	components.ImageFace
}

func (s *ImageRenderSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(ImageRenderable))
}
