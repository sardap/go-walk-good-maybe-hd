package components

import (
	"image/color"

	"golang.org/x/image/font"
)

type TextComponent struct {
	Text  string
	Font  font.Face
	Layer ImageLayer
	Color color.Color
}
