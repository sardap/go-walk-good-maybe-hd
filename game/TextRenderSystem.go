package game

import (
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

type TextRenderSystem struct {
	ents map[uint64]TextRenderable
}

func CreateTextRenderSystem() *TextRenderSystem {
	return &TextRenderSystem{}
}

func (s *TextRenderSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]TextRenderable)
}

func (s *TextRenderSystem) Update(dt float32) {
}

func (s *TextRenderSystem) Render(screen *ebiten.Image) {
	for _, ent := range s.ents {
		trans := ent.GetTransformComponent()
		textCom := ent.GetTextComponent()

		img := ebiten.NewImage(500, 500)
		text.Draw(img, textCom.Text, mplusNormalFont, 0, 50, color.Black)

		op := &ebiten.DrawImageOptions{}
		op.GeoM = *trans.GeoM
		screen.DrawImage(img, op)
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
