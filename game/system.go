package game

import (
	"container/heap"

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

func (r RenderCmds) Len() int {
	return len(r)
}

func (r RenderCmds) Less(i, j int) bool {
	return r[i].GetLayer() < r[j].GetLayer()
}

func (r RenderCmds) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
	r[i].SetIndex(i)
	r[j].SetIndex(j)
}

func (r *RenderCmds) Push(x interface{}) {
	n := len(*r)
	item := x.(RenderCmd)
	item.SetIndex(n)
	*r = append(*r, item)
}

func (r *RenderCmds) Pop() interface{} {
	old := *r
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.SetIndex(-1)
	*r = old[0 : n-1]
	return item
}

func (r *RenderCmds) Update(item *RenderImageCmd) {
	heap.Fix(r, item.index)
}

type RenderingSystem interface {
	Render(*RenderCmds)
}
