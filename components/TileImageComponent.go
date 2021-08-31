package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/math"
)

type TileMap struct {
	Options   DrawOptions
	TilesImg  *ebiten.Image
	Map       []int16
	TileWidth int
	TileXNum  int
}

func CreateTileMap(width int, height int, tiles *ebiten.Image, tileWIdth int) *TileMap {
	result := &TileMap{
		TilesImg:  tiles,
		TileWidth: tileWIdth,
		TileXNum:  width,
		Map:       make([]int16, width*height),
		Options: DrawOptions{
			Scale: math.Vector2{X: 1, Y: 1},
		},
	}

	for i := range result.Map {
		result.Map[i] = -1
	}

	return result
}

func (t *TileMap) Get(x, y int) int16 {
	return t.Map[t.TileXNum*y+x]
}

func (t *TileMap) SetTile(x, y int, tileIdx int16) {
	t.Map[t.TileXNum*y+x] = tileIdx
}

func (t *TileMap) SetRow(startX, y int, tileIdx int16) {
	for x := startX; x < t.TileXNum; x++ {
		t.SetTile(x, y, tileIdx)
	}
}

func (t *TileMap) SetCol(x, startY int, tileIdx int16) {
	yNum := len(t.Map) / t.TileXNum
	for y := startY; y < yNum; y++ {
		t.SetTile(x, y, tileIdx)
	}
}

type TileImageComponent struct {
	Active  bool
	TileMap *TileMap
	Layer   ImageLayer
}
