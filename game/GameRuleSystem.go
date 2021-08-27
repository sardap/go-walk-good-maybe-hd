package game

import (
	"fmt"
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
	ents  map[uint64]interface{}
	world *ecs.World
	space *resolv.Space
}

func CreateGameRuleSystem(space *resolv.Space) *GameRuleSystem {
	return &GameRuleSystem{
		space: space,
	}
}

func (s *GameRuleSystem) Priority() int {
	return int(systemPriorityGameRuleSystem)
}

func (s *GameRuleSystem) New(world *ecs.World) {
	s.ents = make(map[uint64]interface{})
	s.world = world
	mainGameInfo.state = gameStateStarting
}

func (s *GameRuleSystem) updatePlayer(dt float32, player *entity.Player) {
	switch mainGameInfo.state {
	case gameStateStarting:
		if player.TransformComponent.Postion.X > 50 {
			mainGameInfo.state = gameStateScrolling
			mainGameInfo.scrollingSpeed.X = xStartScrollSpeed
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
	} else if move.MoveRight {
		vel.X = horzSpeed
	}

	// Must reset no matter what
	move.MoveLeft = false
	move.MoveRight = false
	move.MoveUp = false

	player.GetVelocityComponent().Vel = vel
}

type Collideable interface {
	ecs.BasicFace
	components.CollisionFace
	components.IdentityFace
}

type Wrapable interface {
	ecs.BasicFace
	components.TransformFace
	components.WrapFace
}

type Scrollable interface {
	ecs.BasicFace
	components.TransformFace
	components.ScrollableFace
}

type Gravityable interface {
	ecs.BasicFace
	components.VelocityFace
	components.GravityFace
	components.IdentityFace
	components.CollisionFace
}

func (s *GameRuleSystem) Update(dt float32) {
	ground := s.space.FilterByTags(entity.TagGround)

	for _, ent := range s.ents {
		if wrapable, ok := ent.(Wrapable); ok {
			trans := wrapable.GetTransformComponent()
			if trans.Postion.X < -wrapable.GetWrapComponent().Threshold {
				trans.Postion.X = wrapable.GetWrapComponent().Threshold
			}
		}

		if scrollable, ok := ent.(Scrollable); ok {
			trans := scrollable.GetTransformComponent().Postion
			trans = trans.Add(mainGameInfo.scrollingSpeed.Mul(float64(dt)))
			scrollable.GetTransformComponent().Postion = trans
		}

		if building, ok := ent.(*Building); ok {
			trans := building.GetTransformComponent()
			if trans.Postion.X+trans.Size.X < 0 {
				defer s.world.RemoveEntity(building.BasicEntity)
				fmt.Printf("Removing %d\n", building.ID())
			}
		}

		if gravityable, ok := ent.(Gravityable); ok {
			vel := gravityable.GetVelocityComponent()
			colCom := gravityable.GetCollisionComponent()

			collision := ground.Resolve(colCom.CollisionShape, 0, mainGameInfo.gravity*float64(dt))
			if !collision.Colliding() {
				vel.Vel = vel.Vel.Add(math.Vector2{Y: mainGameInfo.gravity})
			}
		}

		if ent, ok := ent.(*entity.Player); ok {
			s.updatePlayer(dt, ent)
		}
	}

	mainGameInfo.level.StartX += mainGameInfo.scrollingSpeed.X * float64(dt)
	generateBuildings(s.world)
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
