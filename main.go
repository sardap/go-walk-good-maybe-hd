package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/game"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	game := game.CreateGame()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Walk Good Maybe HD")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
