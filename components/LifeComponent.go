package components

type DamageEvent struct {
	Damage float64
}

type DamageEvents []*DamageEvent

type LifeComponent struct {
	HP           float64
	DamageEvents DamageEvents
}
