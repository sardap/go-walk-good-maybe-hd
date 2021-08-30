package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type ImageRenderable interface {
	ecs.BasicFace
	components.TransformFace
	components.ImageFace
}

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

		*cmds = append(*cmds, &RenderImageCmd{
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

func (s *ImageRenderSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(ImageRenderable))
}

type RenderImageCmd struct {
	Image   *ebiten.Image
	Options *ebiten.DrawImageOptions
	Layer   components.ImageLayer
}

func (c *RenderImageCmd) Draw(screen *ebiten.Image) {
	screen.DrawImage(c.Image, c.Options)
}

func (c *RenderImageCmd) GetLayer() int {
	return int(c.Layer)
}
