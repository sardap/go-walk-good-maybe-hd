package components

type MainGamePlayerState int

// Idk if this should live here
const (
	MainGamePlayerStateGround MainGamePlayerState = iota
	MainGamePlayerStateFalling
	MainGamePlayerStateJumping
)

type MainGamePlayerComponent struct {
	State     MainGamePlayerState
	Speed     float64
	JumpPower float64
}
