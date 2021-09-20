package game

import "image/color"

const (
	scaleMultiplier          = 10
	windowWidth              = 240 * scaleMultiplier
	windowHeight             = 160 * scaleMultiplier
	xStartScrollSpeed        = -100.5
	startingGravity          = 500
	minSpaceBetweenBuildings = 30 * scaleMultiplier
)

var (
	swapColor = color.RGBA{
		R: 255, G: 0, B: 247, A: 255,
	}
)
