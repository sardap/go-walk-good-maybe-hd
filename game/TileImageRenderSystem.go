package game

import (
	"container/heap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type TileImageRenderSystem struct {
	ents []TileImageRenderable
}

func CreateTileImageRenderSystem() *TileImageRenderSystem {
	return &TileImageRenderSystem{}
}

func (s *TileImageRenderSystem) Priority() int {
	return int(systemPriorityTileImageRenderSystem)
}

func (s *TileImageRenderSystem) New(world *ecs.World) {
	s.ents = make([]TileImageRenderable, 0)
}

func (s *TileImageRenderSystem) Update(dt float32) {
}

func (s *TileImageRenderSystem) Render(cmds *RenderCmds) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(trans.Postion.X, trans.Postion.Y)
		op.GeoM.Scale(scaleMultiplier, scaleMultiplier)

		heap.Push(cmds, &RenderTileMapCmd{
			TileImageComponent: ent.GetTileImageComponent(),
			Options:                 op,
		})
	}
}

func (s *TileImageRenderSystem) Add(r TileImageRenderable) {
	s.ents = append(s.ents, r)
}

func (s *TileImageRenderSystem) remove(i int) {
	s.ents = append(s.ents[:i], s.ents[i+1:]...)
}

func (s *TileImageRenderSystem) Remove(e ecs.BasicEntity) {
	for i, ent := range s.ents {
		if ent.GetBasicEntity().ID() == e.ID() {
			s.remove(i)
			break
		}
	}
}

type TileImageRenderable interface {
	ecs.BasicFace
	components.TransformFace
	components.TileImageFace
}

func (s *TileImageRenderSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(TileImageRenderable))
}
