package game

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
)

var (
	audioCtx *audio.Context
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
	if audioCtx == nil {
		audioCtx = audio.NewContext(48000)
	}
	s.audioCtx = audioCtx
}

func (s *SoundSystem) Update(dt float32) {
	for _, ent := range s.ents {
		soundCom := ent.GetSoundComponent()

		// Reloading the assets is fucking stupid should be using rewind with
		// checking if the sound has been changed
		if soundCom.Restart {
			soundCom.Player = nil
			soundCom.Restart = false
		}

		if !soundCom.Active {
			if soundCom.Player != nil && soundCom.Player.IsPlaying() {
				soundCom.Player.Rewind()
				soundCom.Player.Pause()
			}

			continue
		}

		if soundCom.Player == nil {
			var stream io.Reader

			switch soundCom.Sound.SoundType {
			case assets.SoundTypeMp3:
				buffer := bytes.NewReader(soundCom.Sound.Source)
				stream, _ = mp3.DecodeWithSampleRate(soundCom.Sound.SampleRate, buffer)
				if soundCom.Loop {
					mp3Stream := stream.(*mp3.Stream)
					if soundCom.Intro > 0 {
						introLength := int64(int(soundCom.Intro/time.Second) * 4 * soundCom.Sound.SampleRate)

						stream = audio.NewInfiniteLoopWithIntro(
							mp3Stream,
							introLength,
							mp3Stream.Length()-introLength,
						)
					} else {
						stream = audio.NewInfiniteLoop(mp3Stream, mp3Stream.Length())
					}
				}

			case assets.SoundTypeWav:
				buffer := bytes.NewReader(soundCom.Sound.Source)
				stream, _ = wav.DecodeWithSampleRate(soundCom.Sound.SampleRate, buffer)

			default:
				log.Fatalf("Unknown sound type %v", soundCom.Sound.SoundType)
				continue
			}

			soundCom.Player, _ = audio.NewPlayer(s.audioCtx, stream)

			if soundCom.Sound.Volume > 0 {
				soundCom.Player.SetVolume(soundCom.Sound.Volume)
			}

			soundCom.Player.Play()
		}

		if !soundCom.Player.IsPlaying() {
			soundCom.Active = false
		}
	}
}

func (s *SoundSystem) Add(r Soundable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *SoundSystem) Remove(e ecs.BasicEntity) {
	if ent, ok := s.ents[e.ID()]; ok {
		if ent.GetSoundComponent().Player != nil {
			ent.GetSoundComponent().Player.Close()
		}
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
