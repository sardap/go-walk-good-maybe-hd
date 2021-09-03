package game

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type Playerable interface {
	ecs.BasicFace
	components.MainGamePlayerFace
}

type PlayerSystem struct {
	ents         map[uint64]*entity.Player
	mainGameInfo *MainGameInfo
	world        *ecs.World
}

func CreatePlayerSystem(mainGameInfo *MainGameInfo) *PlayerSystem {
	return &PlayerSystem{
		mainGameInfo: mainGameInfo,
	}
}

func (s *PlayerSystem) Priority() int {
	return int(systemPriorityPlayerSystem)
}

func (s *PlayerSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]*entity.Player)
	s.world = world
}

func (s *PlayerSystem) changeToPrepareJump(player *entity.Player) {
	player.State = components.MainGamePlayerStatePrepareJumping

	img, _ := assets.LoadEbitenImage(assets.ImageWhaleJumpTileSet)
	components.ChangeAnimeImage(player, img, 125*time.Millisecond)
}

func (s *PlayerSystem) changeToJumping(player *entity.Player) {
	player.State = components.MainGamePlayerStateJumping

	player.SoundComponent.Active = true
	player.SoundComponent.Restart = true

	img, _ := assets.LoadEbitenImage(assets.ImageWhaleAirTileSet)
	components.ChangeAnimeImage(player, img, 50*time.Millisecond)
	player.JumpTime = 0
}

func (s *PlayerSystem) changeToFlying(player *entity.Player) {
	player.State = components.MainGamePlayerStateFlying
}

func (s *PlayerSystem) changeToIdle(player *entity.Player) {
	player.State = components.MainGamePlayerStateGroundIdling

	img, _ := assets.LoadEbitenImage(assets.ImageWhaleIdleTileSet)
	components.ChangeAnimeImage(player, img, 50*time.Millisecond)
}

func (s *PlayerSystem) changeToWalk(player *entity.Player) {
	player.State = components.MainGamePlayerStateGroundMoving

	img, _ := assets.LoadEbitenImage(assets.ImageWhaleWalkTileSet)
	components.ChangeAnimeImage(player, img, 50*time.Millisecond)
}

func (s *PlayerSystem) Update(dt float32) {
	for _, player := range s.ents {
		switch s.mainGameInfo.State {
		case gameStateStarting:
			if player.TransformComponent.Postion.X > 50 {
				s.mainGameInfo.State = gameStateScrolling
				s.mainGameInfo.ScrollingSpeed.X = xStartScrollSpeed
			}
		case gameStateScrolling:
		}

		playerCom := player.GetMainGamePlayerComponent()

		move := player.GetMovementComponent()
		vel := player.GetVelocityComponent().Vel

		horzSpeed := playerCom.Speed

		switch playerCom.State {
		case components.MainGamePlayerStateGroundIdling:
			if move.MoveUp {
				s.changeToPrepareJump(player)
			} else if move.MoveLeft || move.MoveRight {
				s.changeToWalk(player)
			}

		case components.MainGamePlayerStateGroundMoving:
			if move.MoveUp {
				s.changeToPrepareJump(player)
			} else if !move.MoveLeft && !move.MoveRight {
				s.changeToIdle(player)
			}

		case components.MainGamePlayerStatePrepareJumping:
			horzSpeed = 0

			if player.Cycles >= 1 {
				s.changeToJumping(player)
			}

		case components.MainGamePlayerStateJumping:
			horzSpeed /= 2

			vel.Y -= player.JumpPower
			player.JumpTime += utility.DeltaToDuration(dt)

			if player.JumpTime > time.Duration(1000)*time.Millisecond {
				s.changeToFlying(player)
			}

		case components.MainGamePlayerStateFlying:
			horzSpeed /= 2

			if player.Collisions.CollidingWith(entity.TagGround) {
				s.changeToIdle(player)
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
			bullet.Speed.X = 750
			s.world.AddEntity(bullet)
		}

		// Must reset no matter what
		move.MoveLeft = false
		move.MoveRight = false
		move.MoveUp = false
		move.Shoot = false

		player.GetVelocityComponent().Vel = vel
	}
}

func (s *PlayerSystem) Add(r Playerable) {
	player, ok := r.(*entity.Player)
	if !ok {
		panic("Invlaid player given")
	}

	s.ents[r.GetBasicEntity().ID()] = player
}

func (s *PlayerSystem) Remove(e ecs.BasicEntity) {
	delete(s.ents, e.ID())
}

func (s *PlayerSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(Playerable))
}
