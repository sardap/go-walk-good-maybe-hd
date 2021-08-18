package components

type CollisionComponent struct {
	Active     bool
	Collisions []*CollisionEvent
}

type CollisionEvent struct {
	ID uint64
}
