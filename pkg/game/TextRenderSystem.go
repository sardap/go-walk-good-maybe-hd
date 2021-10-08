package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sardap/walk-good-maybe-hd/pkg/components"
)

type TextRenderSystem struct {
	ents map[uint64]TextRenderable
}

func CreateTextRenderSystem() *TextRenderSystem {
	return &TextRenderSystem{}
}

func (s *TextRenderSystem) Priority() int {
	return int(systemPriorityTextRenderSystem)
}

func (s *TextRenderSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]TextRenderable)
}

func (s *TextRenderSystem) Update(dt float32) {
}

func (s *TextRenderSystem) Render(cmds *RenderCmds) {
	for _, ent := range s.ents {
		*cmds = append(*cmds, &RenderTextCmd{
			TransCom: ent.GetTransformComponent(),
			TextCom:  ent.GetTextComponent(),
		})
	}
}

func (s *TextRenderSystem) Add(r TextRenderable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *TextRenderSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type TextRenderable interface {
	ecs.BasicFace
	components.TransformFace
	components.TextFace
}

func (s *TextRenderSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(TextRenderable))
}

type RenderTextCmd struct {
	TransCom *components.TransformComponent
	TextCom  *components.TextComponent
}

func (c *RenderTextCmd) Draw(screen *ebiten.Image) {
	transCom := c.TransCom
	b := text.BoundString(c.TextCom.Font, c.TextCom.Text)
	x := int(c.TransCom.Postion.X)
	y := b.Dy() + int(transCom.Postion.Y)
	text.Draw(screen, c.TextCom.Text, c.TextCom.Font, x, y, c.TextCom.Color)
}

func (c *RenderTextCmd) GetLayer() int {
	return int(c.TextCom.Layer)
}
