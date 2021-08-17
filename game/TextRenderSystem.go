package game

import (
	"container/heap"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/sardap/walk-good-maybe-hd/components"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type textImageCache struct {
	text string
	img  *ebiten.Image
}

type TextRenderSystem struct {
	ents      map[uint64]TextRenderable
	textCache map[uint64]*textImageCache
}

func CreateTextRenderSystem() *TextRenderSystem {
	return &TextRenderSystem{}
}

func (s *TextRenderSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]TextRenderable)
	s.textCache = make(map[uint64]*textImageCache)
}

func (s *TextRenderSystem) Update(dt float32) {
}

func (s *TextRenderSystem) Render(cmds *RenderCmds) {
	for id, ent := range s.ents {
		trans := ent.GetTransformComponent()
		textCom := ent.GetTextComponent()

		value, ok := s.textCache[id]
		if !ok || value.text != textCom.Text {
			if ok {
				value.img.Dispose()
			} else {
				value = &textImageCache{}
				s.textCache[id] = value
			}

			img := ebiten.NewImage(500, 500)
			text.Draw(img, textCom.Text, mplusNormalFont, 0, 50, color.Black)
			value.text = textCom.Text
			value.img = img
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM = *trans.GeoM

		heap.Push(cmds, &RenderCmd{
			Image:   value.img,
			Options: op,
			Layer:   textCom.Layer,
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
