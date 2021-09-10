package game

import (
	"fmt"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type Lifeable interface {
	ecs.BasicFace
	components.TransformFace
	components.LifeFace
}

type LifeSystem struct {
	ents             map[uint64]Lifeable
	world            *ecs.World
	freePlayerPool   []*entity.SoundPlayer
	activePlayerPool []*entity.SoundPlayer
}

func CreateLifeSystem() *LifeSystem {
	return &LifeSystem{}
}

func (s *LifeSystem) Priority() int {
	return int(systemPriorityDamageSystem)
}

func (s *LifeSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]Lifeable)
	s.world = world
}

func (s *LifeSystem) getPlayer() *entity.SoundPlayer {
	if len(s.freePlayerPool) > 0 {
		result := s.freePlayerPool[len(s.freePlayerPool)-1]
		s.freePlayerPool = s.freePlayerPool[:len(s.freePlayerPool)-1]
		s.activePlayerPool = append(s.activePlayerPool, result)
		return result
	}

	return nil
}

func (s *LifeSystem) freePlayer(toFree *entity.SoundPlayer) {
	toFree.Active = false
	toFree.Restart = false
	for i, player := range s.activePlayerPool {
		if player.ID() == toFree.ID() {
			s.activePlayerPool[i] = s.activePlayerPool[len(s.activePlayerPool)-1]
			s.activePlayerPool = s.activePlayerPool[:len(s.activePlayerPool)-1]
		}
	}
}

func (s *LifeSystem) onRemove(ent Lifeable) {

	if _, ok := ent.(components.BiscuitEnemyFace); ok {
		enemyDeath := s.getPlayer()
		enemyDeath.Sound = components.LoadSound(assets.SoundPdBiscuitDeath)
		enemyDeath.Active = true
		enemyDeath.Restart = true

		biscuitEnemyDeath := entity.CreateBiscuitEnemyDeath()
		biscuitEnemyDeath.Postion = ent.GetTransformComponent().Postion
		biscuitEnemyDeath.Layer = ImagelayerEnemyLayer
		s.world.AddEntity(biscuitEnemyDeath)
	} else if _, ok := ent.(components.UfoBiscuitEnemyFace); ok {
		enemyDeath := s.getPlayer()
		enemyDeath.Sound = components.LoadSound(assets.SoundUfoBiscuitEnemyDeath)
		enemyDeath.Active = true
		enemyDeath.Restart = true

		ufoDeath := entity.CreateUfoBiscuitEnemyDeath()
		ufoDeath.Postion = ent.GetTransformComponent().Postion
		ufoDeath.Layer = ImagelayerEnemyLayer
		s.world.AddEntity(ufoDeath)
	}

	s.world.RemoveEntity(*ent.GetBasicEntity())
}

func (s *LifeSystem) onDamage(ent Lifeable) {
	if player, ok := ent.(*entity.Player); ok {
		fmt.Printf("Player damaged")

		enemyDeath := s.getPlayer()
		enemyDeath.Sound = components.LoadSound(assets.SoundWhaleDamage)
		enemyDeath.Active = true
		enemyDeath.Restart = true
		player.TileMap.Options.InvertColor = true
		// Do I really want to go down this path
		go func() {
			time.Sleep(player.InvincibilityTime)
			player.TileMap.Options.InvertColor = false
		}()
	}
}

func (s *LifeSystem) Update(dt float32) {
	if len(s.freePlayerPool) == 0 {
		for i := 0; i < 10; i++ {
			player := entity.CreateSoundPlayer(assets.SoundPdBiscuitDeath)
			s.freePlayerPool = append(s.freePlayerPool, player)
			s.world.AddEntity(player)
		}
	}

	for _, ent := range s.ents {
		lifeCom := ent.GetLifeComponent()

		defer func(lifeCom *components.LifeComponent) {
			lifeCom.DamageEvents = nil
		}(lifeCom)

		if lifeCom.InvincibilityTimeRemaning > 0 {
			lifeCom.InvincibilityTimeRemaning -= utility.DeltaToDuration(dt)
			continue
		}

		for _, event := range lifeCom.DamageEvents {
			lifeCom.HP -= event.Damage
		}

		if len(lifeCom.DamageEvents) > 0 {
			s.onDamage(ent)
		}

		if lifeCom.HP <= 0 {
			defer s.world.RemoveEntity(*ent.GetBasicEntity())
			defer s.onRemove(ent)
		} else {
			lifeCom.InvincibilityTimeRemaning = lifeCom.InvincibilityTime
		}
	}

	for _, player := range s.activePlayerPool {
		if player.Player != nil && !player.Player.IsPlaying() {
			defer s.freePlayer(player)
		}
	}
}

func (s *LifeSystem) Add(r Lifeable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *LifeSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *LifeSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Lifeable))
}
