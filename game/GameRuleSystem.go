package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type GameRuleSystem struct {
	ents            map[uint64]interface{}
	world           *ecs.World
	enemyDeathSound *entity.SoundPlayer
	info            *Info
	mainGameInfo    *MainGameInfo
}

func CreateGameRuleSystem(info *Info) *GameRuleSystem {
	return &GameRuleSystem{
		info:         info,
		mainGameInfo: info.MainGameInfo,
	}
}

func (s *GameRuleSystem) Priority() int {
	return int(systemPriorityGameRuleSystem)
}

func (s *GameRuleSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]interface{})
	s.world = world
	s.mainGameInfo.State = gameStateStarting
}

type Wrapable interface {
	ecs.BasicFace
	components.TransformFace
	components.WrapFace
}

type Scrollable interface {
	ecs.BasicFace
	components.ScrollableFace
	components.VelocityFace
}

type Gravityable interface {
	ecs.BasicFace
	components.CollisionFace
	components.GravityFace
	components.IdentityFace
	components.VelocityFace
}

type Bulletable interface {
	ecs.BasicFace
	components.TransformFace
	components.BulletFace
	components.CollisionFace
	components.VelocityFace
}

type UfoBiscuitEnemyable interface {
	ecs.BasicFace
	components.TransformFace
	components.UfoBiscuitEnemyFace
	components.CollisionFace
}

type DestoryOnAnimeable interface {
	ecs.BasicFace
	components.AnimeFace
	components.DestoryOnAnimeFace
}

func (s *GameRuleSystem) Update(dt float32) {
	if s.enemyDeathSound == nil {
		s.enemyDeathSound = entity.CreateSoundPlayer(assets.SoundPdBiscuitDeath)
		s.world.AddEntity(s.enemyDeathSound)
	}

	switch s.mainGameInfo.State {
	case gameStateScrolling:
		s.mainGameInfo.ScrollingSpeed.X -= 1 * float64(dt)
	}

	for _, ent := range s.ents {
		if wrapable, ok := ent.(Wrapable); ok {
			trans := wrapable.GetTransformComponent()
			wrap := wrapable.GetWrapComponent()
			trans.Postion = utility.WrapVec2(trans.Postion, wrap.Min, wrap.Max)
		}

		if scrollable, ok := ent.(Scrollable); ok {
			velCom := scrollable.GetVelocityComponent()
			velCom.Vel = velCom.Vel.Add(s.mainGameInfo.ScrollingSpeed.Mul(scrollable.GetScrollableComponent().Modifier))
		}

		if building, ok := ent.(*LevelBlock); ok {
			trans := building.GetTransformComponent()
			if trans.Postion.X+trans.Size.X < 0 {
				defer s.world.RemoveEntity(building.BasicEntity)
			}
		}

		if gravityable, ok := ent.(Gravityable); ok {
			vel := gravityable.GetVelocityComponent()
			vel.Vel = vel.Vel.Add(math.Vector2{Y: s.mainGameInfo.Gravity})
		}

		if bullet, ok := ent.(Bulletable); ok {
			velCom := bullet.GetVelocityComponent()
			velCom.Vel = velCom.Vel.Add(bullet.GetBulletComponent().Speed)
			colCom := bullet.GetCollisionComponent()
			if colCom.Collisions.CollidingWith(entity.TagGround) {
				defer s.world.RemoveEntity(*bullet.GetBasicEntity())
			}
		}

		if ufo, ok := ent.(UfoBiscuitEnemyable); ok {
			ufoCom := ufo.GetUfoBiscuitEnemyComponent()
			transCom := ufo.GetTransformComponent()
			ufoCom.ShootTimeRemaning -= utility.DeltaToDuration(dt)

			if ufoCom.ShootTimeRemaning < 0 {
				bullet := entity.CreateEnemyBullet()
				bullet.Postion.X = transCom.Postion.X + transCom.Size.X/2 - bullet.TransformComponent.Size.X/2
				bullet.Postion.Y = transCom.Postion.Y + transCom.Size.Y + 5
				bullet.Speed.Y = 300
				bullet.Layer = ImageLayerbullet
				bullet.Options.InvertY = true
				bullet.Options.InvertX = false
				defer s.world.AddEntity(bullet)

				ufoCom.ShootTimeRemaning = ufoCom.ShootTime
			}
		}

		if ent, ok := ent.(DestoryOnAnimeable); ok {
			anime := ent.GetAnimeComponent()
			if anime.Cycles >= ent.GetDestoryOnAnimeComponent().CyclesTilDeath {
				defer s.world.RemoveEntity(*ent.GetBasicEntity())
			}
		}
	}

	s.mainGameInfo.Level.StartX += s.mainGameInfo.ScrollingSpeed.X * float64(dt)
	generateCityBuildings(s.info, s.world)
}

func (s *GameRuleSystem) Add(r GameRuleable) {
	s.ents[r.GetBasicEntity().ID()] = r
}

func (s *GameRuleSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

type GameRuleable interface {
	ecs.BasicFace
}

func (s *GameRuleSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(GameRuleable))
}
