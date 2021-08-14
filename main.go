package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/game"
)

const (
	frameOX     = 0
	frameOY     = 32
	frameWidth  = 32
	frameHeight = 32
	frameNum    = 8
)

func main() {
	game := game.CreateGame()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Walk Good Maybe HD")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
