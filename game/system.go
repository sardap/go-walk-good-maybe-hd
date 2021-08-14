package game

import "github.com/hajimehoshi/ebiten/v2"

type RenderingSystem interface {
	Render(screen *ebiten.Image)
}
