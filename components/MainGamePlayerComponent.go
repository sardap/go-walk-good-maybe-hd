package components

type MainGamePlayerState int

// Idk if this should live here
const (
	MainGamePlayerStateGroundIdling MainGamePlayerState = iota
	MainGamePlayerStateGroundMoving
	MainGamePlayerStateFalling
	MainGamePlayerStateJumping
)

type MainGamePlayerComponent struct {
	State     MainGamePlayerState
	Speed     float64
	JumpPower float64
}
