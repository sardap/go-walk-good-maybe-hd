package main

import (
	"fmt"
	_ "image/png"
	"log"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/game"

	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if runtime.GOARCH == "js" || runtime.GOOS == "js" {
		fmt.Printf("Remote host %s \n", assets.Remote)
		ebiten.SetFullscreen(true)
	}

	game := game.CreateGame()

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Walk Good Maybe HD")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
