package components

import "time"

type DamageEvent struct {
	Damage float64
}

type DamageEvents []*DamageEvent

type LifeComponent struct {
	HP                        float64
	InvincibilityTime         time.Duration
	InvincibilityTimeRemaning time.Duration
	DamageEvents              DamageEvents
}
