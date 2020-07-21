package combat

import (
	"fmt"
	"math"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/graph"
)

// MoveIntent is a Component that indicates that this Entity should move.
type MoveIntent struct {
	X, Y float64
}

// Type of this Component.
func (MoveIntent) Type() string {
	return "MoveIntent"
}

// IntentSystem processes intents of Participants in a combat.
type IntentSystem struct {
	mgr *ecs.World
	*event.Bus
	field *geom.Field
}

// NewIntentSystem constructs a new IntentSystem.
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

// Update Characters with Intents.
func (s *IntentSystem) Update() {
	entities := s.mgr.Get([]string{"Participant", "MoveIntent", "Position"})

	for _, e := range entities {
		participant := s.mgr.Component(e, "Participant").(*Participant)
		pos := s.mgr.Component(e, "Position").(*game.Position)
		intent := s.mgr.Component(e, "MoveIntent").(*MoveIntent)

		s.mgr.RemoveComponent(e, intent)

		var start, goal geom.Key
		var stepToWaypoint func(graph.Step) Waypoint

		startHex := s.field.At(pos.Center.X, pos.Center.Y)
		goalHex := s.field.At(intent.X, intent.Y)
		if startHex == nil || goalHex == nil || startHex.Key() == goalHex.Key() {
			// Don't navigate.
			s.Publish(&ParticipantMovementConcluded{Entity: e})
			continue
		}
		start = startHex.Key()
		goal = goalHex.Key()
		stepToWaypoint = func(step graph.Step) Waypoint {
			k := step.V.(geom.Key)
			x, y := s.field.Ktow(k)
			return Waypoint{
				X: x,
				Y: y,
			}
		}

		costs := CostsFuncFactory(s.field, s.mgr, e)
		edges := EdgeFuncFactory(s.field)
		guess := HeuristicFactory(s.field)
		steps := graph.NewSearcher(costs, edges, guess).Search(start, goal)
		if steps == nil {
			fmt.Printf("Search: no path to %v\n", goal)
			s.Publish(&ParticipantMovementConcluded{Entity: e})
			continue
		}

		m := Mover{}
		cost := 0
		for _, step := range steps {
			if step.Cost > float64(participant.ActionPoints.Cur) {
				break
			}
			cost = int(step.Cost)
			m.Moves = append(m.Moves, stepToWaypoint(step))
		}
		participant.ActionPoints.Cur -= cost
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
func ExistsFuncFactory(f *geom.Field) ExistsFunc {
	return func(k geom.Key) bool {
		return f.Get(k) != nil
	}
}

// CostsFuncFactory constructs a CostsFunc that returns the costs of moving from
// one location to another for an Entity from a context.
func CostsFuncFactory(f *geom.Field, mgr *ecs.World, participantEntity ecs.Entity) graph.CostFunc {
	var obstacles []ContextualObstacle
	for _, e := range mgr.Get([]string{"Obstacle"}) {
		// A Participant is not an obstacle to itself.
		if e == participantEntity {
			continue
		}
		obstacle := mgr.Component(e, "Obstacle").(*game.Obstacle)

		if h := f.Get(geom.Key{obstacle.M, obstacle.N}); h == nil {

			continue
		}
		// Translate the Obstacles into ContextualObstacles based on
		// how much of an Obstacle this is to the Mover in this context.
		obstacles = append(obstacles, ContextualObstacle{
			M:    obstacle.M,
			N:    obstacle.N,
			Cost: math.Inf(0), // just pretend these all are total obstacles for now...
		})
	}

	return func(vFrom, vTo graph.Vertex) float64 {
		// from := vFrom.(geom.Key)
		to := vTo.(geom.Key)

		hex := f.Get(to)

		if hex == nil {
			return math.Inf(0)
		}
		cost := 10.0
		for _, o := range obstacles {
			if to == (geom.Key{o.M, o.N}) {
				if math.IsInf(o.Cost, 0) {
					return math.Inf(0)
				}
				cost = cost * o.Cost
			}
		}
		return cost
	}
}

// EdgeFuncFactory generates a function that returns the connected Keys of a Key
// given the context of a Field.
func EdgeFuncFactory(f *geom.Field) graph.EdgeFunc {
	return func(v graph.Vertex) []graph.Vertex {
		key, ok := v.(geom.Key)
		if !ok {
			return []graph.Vertex{}
		}
		candidates := []geom.Key{
			key.ToN(),
			key.ToS(),
			key.ToNW(),
			key.ToNE(),
			key.ToSW(),
			key.ToSE(),
		}
		result := make([]graph.Vertex, 0, 6)
		for _, adj := range candidates {
			if f.Get(adj) != nil {
				result = append(result, adj)
			}
		}

		return result
	}
}

// HeuristicFactory returns a function that calculates the as-the-crow-flies
// distance between two geom.Keys.
// TODO: there's an optimisation here, because this should not require a
// geom.Field to calculate this, but could use the number of steps between any
// two M,N coordinate pairs.
func HeuristicFactory(f *geom.Field) graph.Heuristic {
	return func(v1, v2 graph.Vertex) float64 {
		a, b := v1.(geom.Key), v2.(geom.Key)
		return geom.DistanceSquared(f.Get(a), f.Get(b))
	}
}
