package game

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/math"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type GameRuleSystem struct {
	ents         map[uint64]interface{}
	world        *ecs.World
	space        *resolv.Space
	mainGameInfo *MainGameInfo
}

func CreateGameRuleSystem(mainGameInfo *MainGameInfo, space *resolv.Space) *GameRuleSystem {
	return &GameRuleSystem{
		space:        space,
		mainGameInfo: mainGameInfo,
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

func (s *GameRuleSystem) updatePlayer(dt float32, player *entity.Player) {
	switch s.mainGameInfo.State {
	case gameStateStarting:
		if player.TransformComponent.Postion.X > 50 {
			s.mainGameInfo.State = gameStateScrolling
			s.mainGameInfo.ScrollingSpeed.X = xStartScrollSpeed
		}
	case gameStateScrolling:
		break
	}

	playerCom := player.GetMainGamePlayerComponent()

	move := player.GetMovementComponent()
	vel := player.GetVelocityComponent().Vel

	changeToPrepareJump := func() {
		player.State = components.MainGamePlayerStatePrepareJumping

		img, _ := assets.LoadEbitenImage(assets.ImageWhaleJumpTileSet)
		components.ChangeAnimeImage(player, img, 125*time.Millisecond)
	}

	changeToJumping := func() {
		player.State = components.MainGamePlayerStateJumping

		player.SoundComponent.Active = true
		player.SoundComponent.Restart = true

		img, _ := assets.LoadEbitenImage(assets.ImageWhaleAirTileSet)
		components.ChangeAnimeImage(player, img, 50*time.Millisecond)
		player.JumpTime = 0
	}

	changeToFlying := func() {
		player.State = components.MainGamePlayerStateFlying
	}

	changeToIdle := func() {
		player.State = components.MainGamePlayerStateGroundIdling

		img, _ := assets.LoadEbitenImage(assets.ImageWhaleIdleTileSet)
		components.ChangeAnimeImage(player, img, 50*time.Millisecond)
	}

	changeToWalk := func() {
		player.State = components.MainGamePlayerStateGroundMoving

		img, _ := assets.LoadEbitenImage(assets.ImageWhaleWalkTileSet)
		components.ChangeAnimeImage(player, img, 50*time.Millisecond)
	}

	horzSpeed := playerCom.Speed

	switch playerCom.State {
	case components.MainGamePlayerStateGroundIdling:
		if move.MoveUp {
			changeToPrepareJump()
		} else if move.MoveLeft || move.MoveRight {
			changeToWalk()
		}

	case components.MainGamePlayerStateGroundMoving:
		if move.MoveUp {
			changeToPrepareJump()
		} else if !move.MoveLeft && !move.MoveRight {
			changeToIdle()
		}

	case components.MainGamePlayerStatePrepareJumping:
		horzSpeed = 0

		if player.Cycles >= 1 {
			changeToJumping()
		}

	case components.MainGamePlayerStateJumping:
		horzSpeed /= 2

		vel.Y -= player.JumpPower
		player.JumpTime += utility.DeltaToDuration(dt)

		if player.JumpTime > time.Duration(1000)*time.Millisecond {
			changeToFlying()
		}

	case components.MainGamePlayerStateFlying:
		horzSpeed /= 2

		if player.Collisions.CollidingWith(entity.TagGround) {
			changeToIdle()
		}

	default:
		panic("Unimplemented")
	}

	if move.MoveLeft {
		vel.X = -horzSpeed
		player.TileMap.Options.InvertX = true
	} else if move.MoveRight {
		vel.X = horzSpeed
		player.TileMap.Options.InvertX = false
	}

	player.ShootCooldownRemaning -= utility.DeltaToDuration(dt)
	if move.Shoot && player.ShootCooldownRemaning < 0 {
		player.ShootCooldownRemaning = player.ShootCooldown
		bullet := entity.CreateBullet()
		bullet.Postion.X = player.Postion.X + player.Size.X + 0.5
		bullet.Postion.Y = player.Postion.Y + player.Size.Y/2
		bullet.Layer = bulletImageLayer
		bullet.Speed.X = 70
		s.world.AddEntity(bullet)
	}

	// Must reset no matter what
	move.MoveLeft = false
	move.MoveRight = false
	move.MoveUp = false
	move.Shoot = false

	player.GetVelocityComponent().Vel = vel
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
	components.BulletFace
	components.TransformFace
	components.VelocityFace
	components.CollisionFace
}

func (s *GameRuleSystem) Update(dt float32) {
	ground := s.space.FilterByTags(entity.TagGround)

	for _, ent := range s.ents {
		// Here broken
		if wrapable, ok := ent.(Wrapable); ok {
			trans := wrapable.GetTransformComponent()
			wrap := wrapable.GetWrapComponent()
			trans.Postion = utility.WrapVec2(trans.Postion, wrap.Min, wrap.Max)
		}

		if scrollable, ok := ent.(Scrollable); ok {
			velCom := scrollable.GetVelocityComponent()
			velCom.Vel = velCom.Vel.Add(s.mainGameInfo.ScrollingSpeed)
		}

		if building, ok := ent.(*LevelBlock); ok {
			trans := building.GetTransformComponent()
			if trans.Postion.X+trans.Size.X < 0 {
				defer s.world.RemoveEntity(building.BasicEntity)
			}
		}

		if gravityable, ok := ent.(Gravityable); ok {
			vel := gravityable.GetVelocityComponent()
			colCom := gravityable.GetCollisionComponent()

			collision := ground.Resolve(colCom.CollisionShape, 0, s.mainGameInfo.Gravity*float64(dt))
			if !collision.Colliding() {
				vel.Vel = vel.Vel.Add(math.Vector2{Y: s.mainGameInfo.Gravity})
			}
		}

		if bullet, ok := ent.(Bulletable); ok {
			velCom := bullet.GetVelocityComponent()
			velCom.Vel = velCom.Vel.Add(bullet.GetBulletComponent().Speed)
			postion := bullet.GetTransformComponent().Postion
			colCom := bullet.GetCollisionComponent()
			if postion.X > gameWidth/scaleMultiplier || postion.X < 0 ||
				colCom.Collisions.CollidingWith(entity.TagGround) {
				defer s.world.RemoveEntity(*bullet.GetBasicEntity())
			}
		}

		if ent, ok := ent.(*entity.Player); ok {
			s.updatePlayer(dt, ent)
		}
	}

	s.mainGameInfo.Level.StartX += s.mainGameInfo.ScrollingSpeed.X * float64(dt)
	generateCityBuildings(s.mainGameInfo, s.world)
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
