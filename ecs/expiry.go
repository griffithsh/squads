package ecs

import "time"

// Expiry is a Component that destroys the Entity after the configured duration
// has elapsed.
type Expiry struct {
	Remaining time.Duration
}

// Type of this component.
func (*Expiry) Type() string {
	return "Expiry"
}

// ExpirySystem destroys Entities that have expired.
type ExpirySystem struct {
	mgr *World
}

// NewExpirySystem constructs a new expiry system.
func NewExpirySystem(mgr *World) *ExpirySystem {
	return &ExpirySystem{
		mgr: mgr,
	}
}

// Update the ExpirySystem.
func (es *ExpirySystem) Update(elapsed time.Duration) {
	for _, e := range es.mgr.Get([]string{"Expiry"}) {
		expiry := es.mgr.Component(e, "Expiry").(*Expiry)

		expiry.Remaining -= elapsed
		if expiry.Remaining <= 0 {
			es.mgr.DestroyEntity(e)
		}
	}
}
