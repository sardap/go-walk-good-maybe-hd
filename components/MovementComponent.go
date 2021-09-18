package components

type MovementComponent struct {
	MoveUp           bool
	MoveDown         bool
	MoveLeft         bool
	MoveRight        bool
	Shoot            bool
	ChangeToGamepad  bool
	ChangeToKeyboard bool
	// Debug stuff
	FastGameSpeed         bool
	ToggleCollsionOverlay bool
	// Menu Stuff
	Select bool
}
