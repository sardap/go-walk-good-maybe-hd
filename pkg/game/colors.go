package game

import (
	"image/color"

	"github.com/icza/gox/imagex/colorx"
)

func parseHex(hex string) color.RGBA {
	result, _ := colorx.ParseHexColor(hex)
	return result
}
