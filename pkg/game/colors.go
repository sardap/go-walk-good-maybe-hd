package game

import (
	"image/color"

	"github.com/icza/gox/imagex/colorx"
)

var (
	ClrBlack = color.RGBA{R: 0, G: 0, B: 0, A: 255}
)

func parseHex(hex string) color.RGBA {
	result, _ := colorx.ParseHexColor(hex)
	return result
}
