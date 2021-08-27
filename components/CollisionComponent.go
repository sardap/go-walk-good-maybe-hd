package components

import "github.com/SolarLune/resolv"

type CollisionEvent struct {
	ID   uint64
	Tags []string
}

type CollisionEvents []*CollisionEvent

func (c CollisionEvents) CollidingWith(tag string) bool {
	for _, event := range c {
		for _, otherTag := range event.Tags {
			if otherTag == tag {
				return true
			}
		}
	}

	return false
}

type CollisionComponent struct {
	Active         bool
	Collisions     CollisionEvents
	CollisionShape *resolv.Rectangle
}
