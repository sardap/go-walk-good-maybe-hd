package game

import (
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/assets"
	"github.com/sardap/walk-good-maybe-hd/components"
	"github.com/sardap/walk-good-maybe-hd/entity"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

const (
	startingPlayerJumpPower  = 1250
	maxPlayerJump            = 2000
	startingPlayerAirHorzMod = 0.5
	maxPlayerAirHorzMod      = 1
)

type Playerable interface {
	ecs.BasicFace
	components.MainGamePlayerFace
}

type PlayerSystem struct {
	ents          map[uint64]*entity.Player
	mainGameScene *MainGameScene
	world         *ecs.World
}

func CreatePlayerSystem(mainGameScene *MainGameScene) *PlayerSystem {
	return &PlayerSystem{
		mainGameScene: mainGameScene,
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

	img, _ := assets.LoadEbitenImageAsset(assets.ImageWhaleJumpTileSet)
	components.ChangeAnimeImage(player, img, 125*time.Millisecond)
}

func (s *PlayerSystem) changeToJumping(player *entity.Player) {
	player.State = components.MainGamePlayerStateJumping

	player.Sound = components.LoadSound(assets.SoundByJumpTwo)
	player.SoundComponent.Active = true
	player.SoundComponent.Restart = true

	img, _ := assets.LoadEbitenImageAsset(assets.ImageWhaleAirTileSet)
	components.ChangeAnimeImage(player, img, 50*time.Millisecond)
	player.JumpPowerRemaning = player.JumpPower
}

func (s *PlayerSystem) changeToFlying(player *entity.Player) {
	player.State = components.MainGamePlayerStateFlying
}

func (s *PlayerSystem) changeToIdle(player *entity.Player) {
	player.State = components.MainGamePlayerStateGroundIdling

	img, _ := assets.LoadEbitenImageAsset(assets.ImageWhaleIdleTileSet)
	components.ChangeAnimeImage(player, img, 50*time.Millisecond)
}

func (s *PlayerSystem) changeToWalk(player *entity.Player) {
	player.State = components.MainGamePlayerStateGroundMoving

	img, _ := assets.LoadEbitenImageAsset(assets.ImageWhaleWalkTileSet)
	components.ChangeAnimeImage(player, img, 50*time.Millisecond)
}

func (s *PlayerSystem) Update(dt float32) {
	for _, player := range s.ents {
		switch s.mainGameScene.State {
		case gameStateStarting:
			if player.TransformComponent.Postion.X > 50 {
				s.mainGameScene.State = gameStateScrolling
				s.mainGameScene.ScrollingSpeed.X = xStartScrollSpeed
			}
		case gameStateScrolling:
		}

		playerCom := player.GetMainGamePlayerComponent()

		move := player.GetMovementComponent()
		vel := player.GetVelocityComponent().Vel

		horzSpeed := playerCom.Speed

		// Token stuff
		if player.Collisions.CollidingWith(entity.TagJumpToken) {
			player.Sound = components.LoadSound(assets.SoundByCollect5)
			player.SoundComponent.Active = true
			player.SoundComponent.Restart = true

			player.JumpPower = utility.ClampFloat64(player.JumpPower+startingPlayerJumpPower*0.1, 0, maxPlayerJump)
		}

		if player.Collisions.CollidingWith(entity.TagSpeedToken) {
			player.Sound = components.LoadSound(assets.SoundJdwBlowOne)
			player.SoundComponent.Active = true
			player.SoundComponent.Restart = true

			rand := s.mainGameScene.Rand
			for i := 0; i < rand.Intn(10)+7; i++ {
				speedLine := entity.CreateSpeedLine()
				speedLine.Postion.X = windowWidth + float64(rand.Intn(2000))
				speedLine.Postion.Y = float64(rand.Intn(windowHeight)) + rand.Float64()
				speedLine.ImageComponent.Layer = ImageLayerObjects
				defer s.world.AddEntity(speedLine)
			}

			player.AirHorzSpeedModifier = utility.ClampFloat64(player.AirHorzSpeedModifier+0.1, 0.5, 1)
			extraSpeed := xStartScrollSpeed * 4
			s.mainGameScene.ScrollingSpeed.X += extraSpeed
			go func(extraSpeed float64) {
				time.Sleep(2 * time.Second)
				s.mainGameScene.ScrollingSpeed.X -= extraSpeed
			}(extraSpeed)
		}

		// Player State
		switch playerCom.State {
		case components.MainGamePlayerStateGroundIdling:
			if move.InputPressed(components.InputKindJump) {
				s.changeToPrepareJump(player)
			} else if move.InputPressed(components.InputKindMoveLeft) || move.InputPressed(components.InputKindMoveRight) {
				s.changeToWalk(player)
			}

		case components.MainGamePlayerStateGroundMoving:
			if move.InputPressed(components.InputKindJump) {
				s.changeToPrepareJump(player)
			} else if !move.InputPressed(components.InputKindMoveLeft) && !move.InputPressed(components.InputKindMoveRight) {
				s.changeToIdle(player)
			}

		case components.MainGamePlayerStatePrepareJumping:
			horzSpeed = 0

			if player.Cycles >= 1 {
				s.changeToJumping(player)
			}

		case components.MainGamePlayerStateJumping:
			horzSpeed *= player.AirHorzSpeedModifier

			vel.Y -= player.JumpPowerRemaning
			player.JumpPowerRemaning -= float64(dt) * player.JumpPower / 2

			if player.JumpPowerRemaning < 0 {
				s.changeToFlying(player)
			}

		case components.MainGamePlayerStateFlying:
			horzSpeed *= player.AirHorzSpeedModifier

			if player.Collisions.CollidingWith(entity.TagGround) {
				s.changeToIdle(player)
			}

		default:
			panic("Unimplemented")
		}

		if move.InputPressed(components.InputKindMoveLeft) {
			vel.X = -horzSpeed
			player.TileMap.Options.InvertX = true
		} else if move.InputPressed(components.InputKindMoveRight) {
			vel.X = horzSpeed
			player.TileMap.Options.InvertX = false
		}

		player.ShootCooldownRemaning -= utility.DeltaToDuration(dt)
		if move.InputPressed(components.InputKindShoot) && player.ShootCooldownRemaning < 0 {
			player.ShootCooldownRemaning = player.ShootCooldown
			bullet := entity.CreatePlayerBullet()
			bullet.Postion.X = player.Postion.X
			bullet.Postion.Y = player.Postion.Y + player.Size.Y/2
			bullet.Layer = ImageLayerbullet
			bullet.Speed.X = 750

			xOffset := player.Size.X + 0.5
			// Flip bullet if facing the other way
			if player.TileMap.Options.InvertX {
				xOffset = -xOffset + 90
				bullet.Speed.X = -bullet.Speed.X
				bullet.ImageComponent.Options.InvertX = true
			}
			bullet.Postion.X += xOffset

			s.world.AddEntity(bullet)

			player.Sound = components.LoadSound(assets.SoundByLaserFour)
			player.SoundComponent.Active = true
			player.SoundComponent.Restart = true
		}

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
