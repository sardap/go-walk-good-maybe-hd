package game

import (
	"github.com/EngoEngine/ecs"
	"github.com/sardap/walk-good-maybe-hd/pkg/assets"
	"github.com/sardap/walk-good-maybe-hd/pkg/entity"
)

type MainGameUiSystem struct {
	world    *ecs.World
	player   *entity.Player
	lifeEnt  *entity.BasicTileMap
	jumpEnt  *entity.BasicTileMap
	speedEnt *entity.BasicTileMap
}

func CreateMainGameUiSystem() *MainGameUiSystem {
	return &MainGameUiSystem{}
}

func (s *MainGameUiSystem) Priority() int {
	return int(systemPriorityMainGameUiSystem)
}

func (s *MainGameUiSystem) New(world *ecs.World) {
	s.world = world
}

func (s *MainGameUiSystem) Update(dt float32) {
	if s.player != nil {
		// Life ent
		switch {
		case s.player.HP <= s.player.MaxHp/3:
			s.lifeEnt.TileMap.SetCol(0, 0, assets.IndexUiLifeAmountAmount1)
		case s.player.HP <= s.player.MaxHp/2:
			s.lifeEnt.TileMap.SetCol(0, 0, assets.IndexUiLifeAmountAmount2)
		default:
			s.lifeEnt.TileMap.SetCol(0, 0, assets.IndexUiLifeAmountAmount3)
		}

		const jumpInc = (maxPlayerJump - startingPlayerJumpPower) / 4
		jp := s.player.JumpPower - startingPlayerJumpPower
		// Jump ent
		switch {
		case jp >= jumpInc*4:
			s.jumpEnt.TileMap.SetCol(0, 0, assets.IndexUiJumpAmountAmount4)
		case jp >= jumpInc*3:
			s.jumpEnt.TileMap.SetCol(0, 0, assets.IndexUiJumpAmountAmount3)
		case jp >= jumpInc*2:
			s.jumpEnt.TileMap.SetCol(0, 0, assets.IndexUiJumpAmountAmount2)
		default:
			s.jumpEnt.TileMap.SetCol(0, 0, assets.IndexUiJumpAmountAmount1)
		}

		const speedInc = (maxPlayerAirHorzMod - startingPlayerAirHorzMod) / 4
		sp := s.player.AirHorzSpeedModifier - startingPlayerAirHorzMod
		// speed ent
		switch {
		case sp >= speedInc*5:
			s.speedEnt.TileMap.SetCol(0, 0, assets.IndexUiSpeedAmountAmount5)
		case sp >= speedInc*4:
			s.speedEnt.TileMap.SetCol(0, 0, assets.IndexUiSpeedAmountAmount4)
		case sp >= speedInc*3:
			s.speedEnt.TileMap.SetCol(0, 0, assets.IndexUiSpeedAmountAmount3)
		case sp >= speedInc*2:
			s.speedEnt.TileMap.SetCol(0, 0, assets.IndexUiSpeedAmountAmount2)
		default:
			s.speedEnt.TileMap.SetCol(0, 0, assets.IndexUiSpeedAmountAmount1)
		}

	}
}

func (s *MainGameUiSystem) Render(cmds *RenderCmds) {
}

func (s *MainGameUiSystem) Add(r GameRuleable) {
	if player, ok := r.(*entity.Player); ok {
		// life ent
		s.player = player
		lifeEnt := entity.CreateLifeDisplay()
		lifeEnt.Postion.X = windowWidth - lifeEnt.Size.X - 20
		lifeEnt.Postion.Y = 20
		lifeEnt.Layer = ImageLayerUi
		s.lifeEnt = lifeEnt
		s.world.AddEntity(s.lifeEnt)

		jumpEnt := entity.CreateJumpDisplay()
		jumpEnt.Postion.X = lifeEnt.Postion.X - jumpEnt.Size.X - 20
		jumpEnt.Postion.Y = 20
		jumpEnt.Layer = ImageLayerUi
		s.jumpEnt = jumpEnt
		s.world.AddEntity(jumpEnt)

		speedEnt := entity.CreateSpeedDisplay()
		speedEnt.Postion.X = jumpEnt.Postion.X - speedEnt.Size.X - 20
		speedEnt.Postion.Y = 20
		speedEnt.Layer = ImageLayerUi
		s.speedEnt = speedEnt
		s.world.AddEntity(speedEnt)
	}
}

func (s *MainGameUiSystem) Remove(e ecs.BasicEntity) {
	if s.player != nil && e.ID() == s.player.ID() {
		s.player = nil
		defer s.world.RemoveEntity(*s.lifeEnt.GetBasicEntity())
		s.lifeEnt = nil
		defer s.world.RemoveEntity(*s.jumpEnt.GetBasicEntity())
		s.jumpEnt = nil
		defer s.world.RemoveEntity(*s.speedEnt.GetBasicEntity())
		s.speedEnt = nil
	}
}

func (s *MainGameUiSystem) AddByInterface(o ecs.Identifier) {
	s.Add(o.(GameRuleable))
}
