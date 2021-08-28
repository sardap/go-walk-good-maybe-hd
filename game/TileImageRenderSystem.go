package game

import (
	"image"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
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

		*cmds = append(*cmds, &RenderTileMapCmd{
			TransformComponent: ent.GetTransformComponent(),
			TileImageComponent: ent.GetTileImageComponent(),
			Options:            op,
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

type RenderTileMapCmd struct {
	HeapSortable
	*components.TransformComponent
	*components.TileImageComponent
	Options *ebiten.DrawImageOptions
}

func (c *RenderTileMapCmd) Draw(screen *ebiten.Image) {
	tileSize := c.TileMap.TileWidth
	tileXNum := c.TileMap.TileXNum

	for i, t := range c.TileMap.Map {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.Postion.X, c.Postion.Y)
		op.GeoM.Translate(float64((i%tileXNum)*tileSize), float64((i/tileXNum)*tileSize))
		op.GeoM.Scale(scaleMultiplier, scaleMultiplier)

		sx := int(t) * tileSize
		sy := 0
		screen.DrawImage(c.TileMap.TilesImg.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	}
}

func (c *RenderTileMapCmd) GetLayer() int {
	return int(c.Layer)
}
