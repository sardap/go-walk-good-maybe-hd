package game

import "image/color"

const (
	scaleMultiplier          = 10
	windowWidth              = 1920
	windowHeight             = 1080
	xStartScrollSpeed        = -100.5
	startingGravity          = 500
	minSpaceBetweenBuildings = 30 * scaleMultiplier
)

var (
	swapColor = color.RGBA{
		R: 255, G: 0, B: 247, A: 255,
	}
)
