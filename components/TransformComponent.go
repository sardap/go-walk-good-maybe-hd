package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type TransformComponent struct {
	*ebiten.GeoM
	Vel *image.Point
}
