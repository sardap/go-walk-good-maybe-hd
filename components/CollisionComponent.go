package components

import (
	"github.com/SolarLune/resolv"
	"github.com/sardap/walk-good-maybe-hd/utility"
)

type CollisionEvent struct {
	ID   uint64
	Tags []int
}

type CollisionEvents []*CollisionEvent

func (c CollisionEvents) CollidingWith(tags ...int) bool {
	for _, event := range c {
		if utility.ContainsInt(event.Tags, tags...) {
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
