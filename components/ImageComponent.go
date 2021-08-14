package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageComponent struct {
	Image   *ebiten.Image
	SubRect *image.Rectangle
}
