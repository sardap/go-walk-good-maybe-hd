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
	systemPriorityLifecycleSystem
	systemPriorityAnimeSystem
	systemPriorityVelocitySystem
	systemPriorityScrollingSystem
	systemPriorityGameRuleSystem
	systemPriorityCollisionSystem
	systemPriorityInputSystem
	systemPrioritySoundSystem
)

type HeapSortable struct {
	index int
}

func (s *HeapSortable) GetIndex() int {
	return s.index
}

func (s *HeapSortable) SetIndex(val int) {
	s.index = val
}

type RenderCmd interface {
	Draw(*ebiten.Image)
	GetLayer() int
	GetIndex() int
	SetIndex(int)
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
