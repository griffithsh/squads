package game

import (
	"fmt"
	"math"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
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
	entities := s.mgr.Get([]string{"Actor", "CombatStats", "MoveIntent", "Position"})

	for _, e := range entities {
		a := s.mgr.Component(e, "Actor").(*Actor)
		stats := s.mgr.Component(e, "CombatStats").(*CombatStats)
		pos := s.mgr.Component(e, "Position").(*Position)
		intent := s.mgr.Component(e, "MoveIntent").(*MoveIntent)

		s.mgr.RemoveComponent(e, intent)

		obstacles := s.obstaclesFor(e)

		var start, goal geom.Key
		var exists func(geom.Key) bool
		var costs func(geom.Key) float64
		var stepToWaypoint func(geom.NavigateStep) Waypoint

		switch a.Size {
		case SMALL:
			startHex := s.field.At(int(pos.Center.X), int(pos.Center.Y))
			goalHex := s.field.At(int(intent.X), int(intent.Y))
			if startHex == nil || goalHex == nil {
				// Don't navigate.
				s.Publish(event.ActorMovementConcluded{Entity: e})
				continue
			}
			start = geom.Key{M: startHex.M, N: startHex.N}
			goal = geom.Key{M: goalHex.M, N: goalHex.N}
			exists = func(k geom.Key) bool {
				return s.field.Get(k.M, k.N) != nil
			}
			costs = func(k geom.Key) float64 {
				hex := s.field.Get(k.M, k.N)

				if hex == nil {
					return math.Inf(0)
				}
				for _, o := range obstacles {
					if o.M == hex.M && o.N == hex.N {
						return o.Cost
					}
				}
				return 1.0
			}
			stepToWaypoint = func(step geom.NavigateStep) Waypoint {
				h := s.field.Get(step.M, step.N)
				return Waypoint{
					X: h.X(),
					Y: h.Y(),
				}
			}
		case MEDIUM:
			startHex := s.field.At4(int(pos.Center.X), int(pos.Center.Y))
			goalHex := s.field.At4(int(intent.X), int(intent.Y))
			if startHex == nil || goalHex == nil {
				// Don't navigate.
				s.Publish(event.ActorMovementConcluded{Entity: e})
				continue
			}
			start = geom.Key{M: startHex.M, N: startHex.N}
			goal = geom.Key{M: goalHex.M, N: goalHex.N}
			exists = func(k geom.Key) bool {
				return s.field.Get4(k.M, k.N) != nil
			}
			costs = func(k geom.Key) float64 {
				hex := s.field.Get4(k.M, k.N)

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
			stepToWaypoint = func(step geom.NavigateStep) Waypoint {
				h := s.field.Get4(step.M, step.N)
				return Waypoint{
					X: h.X(),
					Y: h.Y(),
				}
			}
		case LARGE:
			startHex := s.field.At7(int(pos.Center.X), int(pos.Center.Y))
			goalHex := s.field.At7(int(intent.X), int(intent.Y))
			if startHex == nil || goalHex == nil {
				// Don't navigate.
				s.Publish(event.ActorMovementConcluded{Entity: e})
				continue
			}
			start = geom.Key{M: startHex.M, N: startHex.N}
			goal = geom.Key{M: goalHex.M, N: goalHex.N}
			exists = func(k geom.Key) bool {
				return s.field.Get7(k.M, k.N) != nil
			}
			costs = func(k geom.Key) float64 {
				hex := s.field.Get7(k.M, k.N)

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
			stepToWaypoint = func(step geom.NavigateStep) Waypoint {
				h := s.field.Get7(step.M, step.N)
				return Waypoint{
					X: h.X(),
					Y: h.Y(),
				}
			}
		}

		steps, err := geom.Navigate(start, goal, exists, costs)
		if err != nil {
			fmt.Printf("Navigate: %v\n", err)
			s.Publish(event.ActorMovementConcluded{Entity: e})
			continue
		}

		m := Mover{}
		cost := 0
		for _, step := range steps {
			if int(step.Cost) > stats.ActionPoints {
				break
			}
			cost = int(step.Cost)
			m.Moves = append(m.Moves, stepToWaypoint(step))
		}
		stats.ActionPoints -= cost
		s.mgr.AddComponent(e, &m)
	}
}

// obstaclesFor provides the obstacles that exist in the world that will impede
// (or potentially speed up) the navigation of a character.
func (s *IntentSystem) obstaclesFor(actor ecs.Entity) []ContextualObstacle {
	var obstacles []ContextualObstacle
	for _, e := range s.mgr.Get([]string{"Obstacle"}) {
		// An Actor is not an obstacle to itself.
		if e == actor {
			continue
		}
		obstacle := s.mgr.Component(e, "Obstacle").(*Obstacle)

		switch obstacle.ObstacleType {
		// case SmallActor: // SmallActor handled as default
		case MediumActor:
			hex := s.field.Get4(obstacle.M, obstacle.N)

			if hex == nil {
				continue
			}
			for _, h := range hex.Hexes() {
				// Translate the Obstacles into ContextualObstacles based on
				// how much of an Obstacle this is to the Mover in this context.
				obstacles = append(obstacles, ContextualObstacle{
					M:    h.M,
					N:    h.N,
					Cost: math.Inf(0), // just pretend these all are total obstacles for now
				})
			}

		case LargeActor:
			hex := s.field.Get7(obstacle.M, obstacle.N)

			if hex == nil {
				continue
			}
			for _, h := range hex.Hexes() {
				// Translate the Obstacles into ContextualObstacles based on
				// how much of an Obstacle this is to the Mover in this context.
				obstacles = append(obstacles, ContextualObstacle{
					M:    h.M,
					N:    h.N,
					Cost: math.Inf(0), // just pretend these all are total obstacles for now
				})
			}

		default:
			hex := s.field.Get(obstacle.M, obstacle.N)

			if hex == nil {
				continue
			}
			// Translate the Obstacles into ContextualObstacles based on
			// how much of an Obstacle this is to the Mover in this context.
			obstacles = append(obstacles, ContextualObstacle{
				M:    obstacle.M,
				N:    obstacle.N,
				Cost: math.Inf(0), // just pretend these all are total obstacles for now
			})
		}
	}
	return obstacles
}
