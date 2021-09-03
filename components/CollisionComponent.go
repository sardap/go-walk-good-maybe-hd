package components

import (
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type CollisionEvent struct {
	ID   uint64
	Tags []string
}

type CollisionEvents []*CollisionEvent

func (c CollisionEvents) CollidingWith(tags ...string) bool {
	for _, event := range c {
		if utility.ContainsString(event.Tags, tags...) {
			return true
		}
	}

	return false
}

type CollisionComponent struct {
	Active         bool
	Collisions     CollisionEvents
	CollisionShape *resolv.Rectangle
}
