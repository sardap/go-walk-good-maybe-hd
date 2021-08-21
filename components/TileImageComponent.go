package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type TileImageComponent struct {
	Active    bool
	TilesImg  *ebiten.Image
	TilesMap  []int
	TileWidth int
	TileXNum  int
	Layer     ImageLayer
}
