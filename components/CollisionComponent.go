package components

type CollisionEvent struct {
	ID uint64
}

type CollisionEvents []*CollisionEvent

type CollisionComponent struct {
	Active     bool
	Collisions CollisionEvents
}
