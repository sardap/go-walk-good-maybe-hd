package game

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type systemPriority int

const (
	systemPriorityImageRenderSystem systemPriority = iota
	systemPriorityKaraokeRenderSystem
	systemPriorityTileImageRenderSystem
	systemPriorityTextRenderSystem
	systemPriorityMainGameUiSystem
	systemPrioritySoundSystem
	systemPriorityDamageSystem
	systemPriorityDestoryBoundSystem
	systemPriorityAnimeSystem
	systemPriorityResolvSystem
	systemPriorityVelocitySystem
	systemPriorityDumbVelocitySystem
	systemPriorityConstantSpeedSystem
	systemPriorityScrollingSystem
	systemPriorityGameRuleSystem
	systemPriorityPlayerSystem
	systemPriorityEnemyBiscuitSystem
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
