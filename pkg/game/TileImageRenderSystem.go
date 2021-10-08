package game

import (
	"image"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
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
		*cmds = append(*cmds, &RenderTileMapCmd{
			TransformComponent: ent.GetTransformComponent(),
			TileImageComponent: ent.GetTileImageComponent(),
		})
	}
}

func (s *TileImageRenderSystem) Add(r TileImageRenderable) {
	s.ents = append(s.ents, r)
}

func (s *TileImageRenderSystem) Remove(e ecs.BasicEntity) {
	for i, ent := range s.ents {
		if ent.GetBasicEntity().ID() == e.ID() {
			s.ents = append(s.ents[:i], s.ents[i+1:]...)
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
	*components.TransformComponent
	*components.TileImageComponent
}

func (c *RenderTileMapCmd) Draw(screen *ebiten.Image) {
	tileSize := c.TileMap.TileWidth
	tileXNum := c.TileMap.TileXNum

	options := c.TileMap.Options

	for i, t := range c.TileMap.Map {
		op := &ebiten.DrawImageOptions{}

		if options.InvertX {
			op.GeoM.Scale(-1, 1)
			op.GeoM.Translate(float64(tileSize), 0)
		}

		if options.InvertY {
			op.GeoM.Scale(1, -1)
			op.GeoM.Translate(0, float64(tileSize))
		}

		op.GeoM.Translate(c.Postion.X, c.Postion.Y)
		op.GeoM.Translate(float64((i%tileXNum)*tileSize), float64((i/tileXNum)*tileSize))
		op.GeoM.Scale(options.Scale.X, options.Scale.Y)

		sx := int(t) * tileSize
		rect := image.Rect(sx, 0, sx+tileSize, c.TileMap.TilesImg.Bounds().Dy())
		subImg := c.TileMap.TilesImg.SubImage(rect).(*ebiten.Image)

		if options.InvertColor {
			op.ColorM.Scale(-1, -1, -1, 1)
			op.ColorM.Translate(1, 1, 1, 0)
		}

		screen.DrawImage(subImg, op)
	}
}

func (c *RenderTileMapCmd) GetLayer() int {
	return int(c.Layer)
}
