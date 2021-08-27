package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type AnimeSystem struct {
	ents map[uint64]Animeable
}

func CreateAnimeSystem() *AnimeSystem {
	return &AnimeSystem{}
}

func (s *AnimeSystem) Priority() int {
	return int(systemPriorityAnimeSystem)
}

func frameCount(img *components.TileImageComponent) int {
	return img.TileMap.TilesImg.Bounds().Max.X / img.TileMap.TileWidth
}

func (s *AnimeSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Animeable)
}

func (s *AnimeSystem) Update(dt float32) {
	for _, ent := range s.ents {
		anime := ent.GetAnimeComponent()
		img := ent.GetTileImageComponent()

		anime.FrameRemaining -= utility.DeltaToDuration(dt)
		if anime.FrameRemaining < 0 {
			nextFrame := img.TileMap.Map[0] + 1
			if nextFrame >= int16(frameCount(img)) {
				anime.Cycles++
				nextFrame = 0
			}
			img.TileMap.Map[0] = nextFrame
			anime.FrameRemaining = anime.FrameDuration
		}
	}
}

func (s *AnimeSystem) Add(r Animeable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *AnimeSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type Animeable interface {
	ecs.BasicFace
	components.AnimeFace
	components.TileImageFace
}

func (s *AnimeSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Animeable))
}
