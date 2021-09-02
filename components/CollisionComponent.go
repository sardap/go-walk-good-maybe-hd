package components

import "github.com/SolarLune/resolv"

type CollisionEvent struct {
	ID   uint64
	Tags []string
}

type CollisionEvents []*CollisionEvent

func (c CollisionEvents) CollidingWith(tags ...string) bool {
	for _, event := range c {
		for _, otherTag := range event.Tags {
			for _, sTag := range tags {
				if otherTag == sTag {
					return true
				}
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
