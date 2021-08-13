package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type ImageRenderImageSystem struct {
	ents map[uint64]ImageRenderable
}

func CreateRenderSystem() *ImageRenderImageSystem {
	return &ImageRenderImageSystem{}
}

func (s *ImageRenderImageSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]ImageRenderable)
}

func (s *ImageRenderImageSystem) Update(dt float32) {
}

func (s *ImageRenderImageSystem) Render(screen *ebiten.Image) {
	for _, ent := range s.ents {
		screen.DrawImage(ent.GetImageComponent().Image, ent.GetTransformComponent().DrawImageOptions)
	}
}

func (s *ImageRenderImageSystem) Add(r ImageRenderable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *ImageRenderImageSystem) Remove(e ecs.BasicEntity) {

}

type ImageRenderable interface {
	ecs.BasicFace
	components.TransformFace
	components.ImageFace
}

func (s *ImageRenderImageSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(ImageRenderable))
}
