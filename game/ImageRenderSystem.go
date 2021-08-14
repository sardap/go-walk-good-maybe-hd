package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type ImageRenderSystem struct {
	ents map[uint64]ImageRenderable
}

func CreateImageRenderSystem() *ImageRenderSystem {
	return &ImageRenderSystem{}
}

func (s *ImageRenderSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]ImageRenderable)
}

func (s *ImageRenderSystem) Update(dt float32) {
}

func (s *ImageRenderSystem) Render(screen *ebiten.Image) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()
		img := ent.GetImageComponent()

		op := &ebiten.DrawImageOptions{}
		op.GeoM = *trans.GeoM

		if img.SubRect == nil {
			screen.DrawImage(img.Image, op)
		} else {
			screen.DrawImage(img.Image.SubImage(*img.SubRect).(*ebiten.Image), op)
		}

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
