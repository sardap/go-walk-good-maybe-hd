package components

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageLayer int

type Animeable interface {
	AnimeFace
	TileImageFace
}

func ChangeAnimeImage(a Animeable, img *ebiten.Image, frameDuration time.Duration) {
	imgCom := a.GetTileImageComponent()
	animeCom := a.GetAnimeComponent()

	imgCom.TileMap.TilesImg = img
	imgCom.TileMap.SetTile(0, 0, 0)
	animeCom.FrameDuration = frameDuration
	animeCom.FrameRemaining = animeCom.FrameDuration
	animeCom.Cycles = 0
}
