package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageComponent struct {
	Active  bool
	Options DrawOptions
	Image   *ebiten.Image
	SubRect *image.Rectangle
	Layer   ImageLayer
}

func (i *ImageComponent) Size() (width, height int) {
	if i.SubRect == nil {
		return i.Image.Size()
	}

	return i.SubRect.Dx(), i.SubRect.Dy()
}
