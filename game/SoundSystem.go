package game

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/sardap/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
)

type SoundSystem struct {
	ents     map[uint64]Soundable
	audioCtx *audio.Context
}

func CreateSoundSystem() *SoundSystem {
	return &SoundSystem{}
}

func (s *SoundSystem) Priority() int {
	return int(systemPrioritySoundSystem)
}

func (s *SoundSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Soundable)
	s.audioCtx = audio.NewContext(48000)
}

func (s *SoundSystem) Update(dt float32) {
	for _, ent := range s.ents {
		soundCom := ent.GetSoundComponent()
		if soundCom.Active {
			if soundCom.Player == nil {
				switch soundCom.Sound.SoundType {
				case components.SoundTypeMp3:
					buffer := bytes.NewReader(soundCom.Sound.Source)
					var stream io.Reader
					stream, _ = mp3.DecodeWithSampleRate(assets.Mp3SampleRate, buffer)
					if soundCom.Loop {
						mp3Stream := stream.(*mp3.Stream)
						if soundCom.Intro > 0 {
							introLength := int64((soundCom.Intro / time.Second) * 4 * assets.Mp3SampleRate)

							stream = audio.NewInfiniteLoopWithIntro(
								mp3Stream,
								introLength,
								mp3Stream.Length()-introLength)
						} else {
							stream = audio.NewInfiniteLoop(mp3Stream, mp3Stream.Length())
						}
					}

					soundCom.Player, _ = audio.NewPlayer(s.audioCtx, stream)

				default:
					log.Fatalf("Unkown sound type %v", soundCom.Sound.SoundType)
					continue
				}
			}

			if !soundCom.Player.IsPlaying() {
				soundCom.Player.Play()
			}
		}
	}
}

func (s *SoundSystem) Add(r Soundable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *SoundSystem) Remove(e ecs.BasicEntity) {
	if ent, ok := s.ents[e.ID()]; ok {
		ent.GetSoundComponent().Player.Close()
	}

	delete(s.ents, e.ID())
}

type Soundable interface {
	ecs.BasicFace
	components.SoundFace
}

func (s *SoundSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Soundable))
}
