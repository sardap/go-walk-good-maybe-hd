package game

import (
	"container/heap"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sardap/walk-good-maybe-hd/components"
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

type RenderImageCmd struct {
	HeapSortable
	Image   *ebiten.Image
	Options *ebiten.DrawImageOptions
	Layer   components.ImageLayer
}

func (c *RenderImageCmd) Draw(screen *ebiten.Image) {
	screen.DrawImage(c.Image, c.Options)
}

func (c *RenderImageCmd) GetLayer() int {
	return int(c.Layer)
}

type RenderTileMapCmd struct {
	HeapSortable
	*components.TransformComponent
	*components.TileImageComponent
	Options *ebiten.DrawImageOptions
}

func (c *RenderTileMapCmd) Draw(screen *ebiten.Image) {
	tileSize := c.TileMap.TileWidth
	tileXNum := c.TileMap.TileXNum

	for i, t := range c.TileMap.Map {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(c.Postion.X, c.Postion.Y)
		op.GeoM.Translate(float64((i%tileXNum)*tileSize), float64((i/tileXNum)*tileSize))
		op.GeoM.Scale(scaleMultiplier, scaleMultiplier)

		sx := (int(t) % tileXNum) * tileSize
		sy := (int(t) / tileXNum) * tileSize
		screen.DrawImage(c.TileMap.TilesImg.SubImage(image.Rect(sx, sy, sx+tileSize, sy+tileSize)).(*ebiten.Image), op)
	}
}

func (c *RenderTileMapCmd) GetLayer() int {
	return int(c.Layer)
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
