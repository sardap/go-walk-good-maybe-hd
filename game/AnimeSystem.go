package game

import (
	"image"

	"github.com/sardap/ecs"
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

func frameCount(anime *components.AnimeComponent, img *components.ImageComponent) int {
	return img.Image.Bounds().Max.X / anime.FrameWidth
}

func (s *AnimeSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Animeable)
}

func (s *AnimeSystem) Update(dt float32) {
	for _, ent := range s.ents {
		anime := ent.GetAnimeComponent()
		img := ent.GetImageComponent()
		anime.FrameRemaining -= utility.DeltaToDuration(dt)
		if anime.FrameRemaining < 0 {
			anime.CurrentFrame = utility.WrapInt(anime.CurrentFrame+1, 0, frameCount(anime, img))
			anime.FrameRemaining = anime.FrameDuration
			sx, sy := anime.CurrentFrame*anime.FrameWidth, 0
			img.SubRect = &image.Rectangle{
				image.Pt(sx, sy), image.Pt(sx+anime.FrameWidth, sy+anime.FrameHeight),
			}
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
	components.ImageFace
}

func (s *AnimeSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Animeable))
}
