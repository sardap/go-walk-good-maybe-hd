package game

import (
	"container/heap"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type systemPriority int

const (
	systemPriorityImageRenderSystem systemPriority = iota
	systemPriorityAnimeSystem
	systemPriorityTextRenderSystem
	systemPriorityVelocitySystem
	systemPriorityScrollingSystem
	systemPriorityGameRuleSystem
	systemPriorityCollisionSystem
	systemPriorityInputSystem
	systemPrioritySoundSystem
)

type RenderCmd struct {
	Image   *ebiten.Image
	Options *ebiten.DrawImageOptions
	Layer   components.ImageLayer
	index   int
}

type RenderCmds []*RenderCmd

func (r RenderCmds) Len() int {
	return len(r)
}

func (r RenderCmds) Less(i, j int) bool {
	return r[i].Layer < r[j].Layer
}

func (r RenderCmds) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
	r[i].index = i
	r[j].index = j
}

func (r *RenderCmds) Push(x interface{}) {
	n := len(*r)
	item := x.(*RenderCmd)
	item.index = n
	*r = append(*r, item)
}

func (r *RenderCmds) Pop() interface{} {
	old := *r
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*r = old[0 : n-1]
	return item
}

func (r *RenderCmds) Update(item *RenderCmd) {
	heap.Fix(r, item.index)
}

type RenderingSystem interface {
	Render(*RenderCmds)
}

type boundingBox interface {
	components.TransformFace
	components.ImageFace
}

func bounds(b boundingBox) (float64, float64, float64, float64) {
	trans := b.GetTransformComponent()
	imgW, imgH := b.GetImageComponent().Size()

	x1, w := trans.Element(0, 2), trans.Element(1, 1)*float64(imgW)
	y1, h := trans.Element(1, 2), trans.Element(0, 0)*float64(imgH)

	return x1, w, y1, h
}
