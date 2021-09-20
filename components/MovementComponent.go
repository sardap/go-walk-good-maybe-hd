package components

// DO NOT DIRECTLY CREATE USE CreateMovementComponent()
type MovementComponent struct {
	PressedDuration []int
	JustPressed     []bool
	JustReleased    []bool
}

func CreateMovementComponent() *MovementComponent {
	return &MovementComponent{
		PressedDuration: make([]int, InputKindLength),
		JustPressed:     make([]bool, InputKindLength),
		JustReleased:    make([]bool, InputKindLength),
	}
}

func (m *MovementComponent) InputPressedDuration(input InputKind) int {
	return m.PressedDuration[input]
}

func (m *MovementComponent) InputJustPressed(input InputKind) bool {
	return m.JustPressed[input]
}

func (m *MovementComponent) InputPressed(input InputKind) bool {
	return m.PressedDuration[input] > 0
}

func (m *MovementComponent) InputJustReleased(input InputKind) bool {
	return m.JustReleased[input]
}
