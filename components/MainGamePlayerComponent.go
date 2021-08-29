package components

import "time"

type MainGamePlayerState int

// Idk if this should live here
const (
	MainGamePlayerStateGroundIdling MainGamePlayerState = iota
	MainGamePlayerStateGroundMoving
	MainGamePlayerStateFlying
	MainGamePlayerStatePrepareJumping
	MainGamePlayerStateJumping
)

type MainGamePlayerComponent struct {
	State                 MainGamePlayerState
	Speed                 float64
	JumpPower             float64
	JumpTime              time.Duration
	ShootCooldown         time.Duration
	ShootCooldownRemaning time.Duration
}
