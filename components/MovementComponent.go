package components

type MovementComponent struct {
	MoveUp    bool
	MoveDown  bool
	MoveLeft  bool
	MoveRight bool
	Shoot     bool
	// Debug stuff
	ScrollSpeedUp         bool
	ToggleCollsionOverlay bool
}
