package game

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type systemPriority int

const (
	systemPriorityImageRenderSystem systemPriority = iota
	systemPriorityTileImageRenderSystem
	systemPriorityTextRenderSystem
	systemPrioritySoundSystem
	systemPriorityLifecycleSystem
	systemPriorityAnimeSystem
	systemPriorityResolvSystem
	systemPriorityVelocitySystem
	systemPriorityScrollingSystem
	systemPriorityGameRuleSystem
	systemPriorityCollisionSystem
	systemPriorityInputSystem
)

type RenderCmd interface {
	Draw(*ebiten.Image)
	GetLayer() int
}

type RenderCmds []RenderCmd

func (r RenderCmds) Sort() {
	sort.Slice(r, func(i, j int) bool {
		return r[i].GetLayer() < r[j].GetLayer()
	})
}

type RenderingSystem interface {
	Render(*RenderCmds)
}
