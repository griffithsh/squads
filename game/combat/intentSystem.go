package combat

import (
	"fmt"
	"math"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

type IntentSystem struct {
	mgr *ecs.World
	*event.Bus
	field *geom.Field
}

func NewIntentSystem(mgr *ecs.World, bus *event.Bus, field *geom.Field) *IntentSystem {
	return &IntentSystem{
		mgr:   mgr,
		Bus:   bus,
		field: field,
	}
}

// ContextualObstacle captures how much of an obstacle this is to the navigator.
// A bird can fly right over a tree, a snake is not impeded by a swamp. A horse
// runs fastest when the ground is level and clear. The Cost multiplies the
// normal traversal time. A Cost of 2 implies that taking this path is twice as
// long as it normally would be. A cost of Infinity marks something that is
// completely impassable.
type ContextualObstacle struct {
	M, N int

	Cost float64
}

// Update Actors with Intents.
func (s *IntentSystem) Update() {
	entities := s.mgr.Get([]string{"Actor", "MoveIntent", "Position"})

	for _, e := range entities {
		a := s.mgr.Component(e, "Actor").(*Actor)
		pos := s.mgr.Component(e, "Position").(*game.Position)
		intent := s.mgr.Component(e, "MoveIntent").(*game.MoveIntent)

		s.mgr.RemoveComponent(e, intent)

		var start, goal geom.Key
		var stepToWaypoint func(geom.NavigateStep) Waypoint
		exists := ExistsFuncFactory(s.field, a.Size)
		costs := CostsFuncFactory(s.field, s.mgr, e)

		f := game.AdaptField(s.field, a.Size)
		startHex := f.At(int(pos.Center.X), int(pos.Center.Y))
		goalHex := f.At(int(intent.X), int(intent.Y))
		if startHex == nil || goalHex == nil {
			// Don't navigate.
			s.Publish(&game.CombatActorMovementConcluded{Entity: e})
			continue
		}
		start = startHex.Key()
		goal = goalHex.Key()
		stepToWaypoint = func(step geom.NavigateStep) Waypoint {
			h := f.Get(step.M, step.N)
			return Waypoint{
				X: h.X(),
				Y: h.Y(),
			}
		}

		steps, err := geom.Navigate(start, goal, exists, costs)
		if err != nil {
			fmt.Printf("Navigate: %v\n", err)
			s.Publish(&game.CombatActorMovementConcluded{Entity: e})
			continue
		}

		m := Mover{}
		cost := 0
		for _, step := range steps {
			if int(step.Cost) > a.ActionPoints.Cur {
				break
			}
			cost = int(step.Cost)
			m.Moves = append(m.Moves, stepToWaypoint(step))
		}
		a.ActionPoints.Cur -= cost
		s.Publish(&StatModified{
			Entity: e,
			Stat:   game.ActionStat,
			Amount: -cost,
		})

		s.mgr.AddComponent(e, &m)
	}
}

// ExistsFunc is a function that will return whether a logical hex exists for
// the given M,N coordinates, in a specific context.
type ExistsFunc func(geom.Key) bool

// ExistsFuncFactory constructs ExistsFuncs from a context.
func ExistsFuncFactory(f *geom.Field, sz game.CharacterSize) ExistsFunc {
	return func(k geom.Key) bool {
		return game.AdaptField(f, sz).Get(k.M, k.N) != nil
	}
}

// CostsFunc is a function that will return the cost of moving to M,N in a specific context.
type CostsFunc func(geom.Key) float64

// CostsFuncFactory constructs a CostsFunc that returns the costs of moving to
// an M,N for an Entity from a context.
func CostsFuncFactory(f *geom.Field, mgr *ecs.World, actor ecs.Entity) CostsFunc {
	var obstacles []ContextualObstacle
	for _, e := range mgr.Get([]string{"Obstacle"}) {
		// An Actor is not an obstacle to itself.
		if e == actor {
			continue
		}
		obstacle := mgr.Component(e, "Obstacle").(*game.Obstacle)

		h := game.AdaptFieldObstacle(f, obstacle.ObstacleType).Get(obstacle.M, obstacle.N)
		if h == nil {
			continue
		}
		for _, h := range h.Hexes() {
			// Translate the Obstacles into ContextualObstacles based on
			// how much of an Obstacle this is to the Mover in this context.
			obstacles = append(obstacles, ContextualObstacle{
				M:    h.M,
				N:    h.N,
				Cost: math.Inf(0), // just pretend these all are total obstacles for now
			})
		}
	}
	a := mgr.Component(actor, "Actor").(*Actor)

	return func(k geom.Key) float64 {
		hex := game.AdaptField(f, a.Size).Get(k.M, k.N)

		if hex == nil {
			return math.Inf(0)
		}
		cost := 1.0
		for _, hex := range hex.Hexes() {
			for _, o := range obstacles {
				if o.M == hex.M && o.N == hex.N {
					if math.IsInf(o.Cost, 0) {
						return math.Inf(0)
					}
					cost = cost * o.Cost
				}
			}
		}
		return cost
	}
}
