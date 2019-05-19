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

func (s *IntentSystem) Update() {
	entities := s.mgr.Get([]string{"Actor", "MoveIntent", "Position"})

	for _, e := range entities {
		a := s.mgr.Component(e, "Actor").(*Actor)
		pos := s.mgr.Component(e, "Position").(*Position)
		intent := s.mgr.Component(e, "MoveIntent").(*MoveIntent)

		s.mgr.RemoveComponent(e, intent)

		obstacles := s.obstaclesFor(e)

		var steps []geom.Positioned
		var err error
		switch a.Size {
		case SMALL:
			steps, err = geom.Navigate(s.field.At(int(pos.Center.X), int(pos.Center.Y)), s.field.At(int(intent.X), int(intent.Y)), obstacles)
		case MEDIUM:
			steps, err = geom.Navigate4(s.field.At4(int(pos.Center.X), int(pos.Center.Y)), s.field.At4(int(intent.X), int(intent.Y)), obstacles)
		case LARGE:
			steps, err = geom.Navigate7(s.field.At7(int(pos.Center.X), int(pos.Center.Y)), s.field.At7(int(intent.X), int(intent.Y)), obstacles)
		}
		if err != nil {
			fmt.Printf("no path there: %v\n", err)
			s.Publish(event.ActorMovementConcluded{Entity: e})
		} else {
			m := Mover{}

			for _, step := range steps {
				m.Moves = append(m.Moves, Waypoint{X: step.X(), Y: step.Y()})
			}

			s.mgr.AddComponent(e, &m)
		}
	}
}

// obstaclesFor provides the obstacles that exist in the world that will impede
// (or potentially speed up) the navigation of a character.
func (s *IntentSystem) obstaclesFor(actor ecs.Entity) []geom.ContextualObstacle {
	var obstacles []geom.ContextualObstacle
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
				obstacles = append(obstacles, geom.ContextualObstacle{
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
				obstacles = append(obstacles, geom.ContextualObstacle{
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
			obstacles = append(obstacles, geom.ContextualObstacle{
				M:    obstacle.M,
				N:    obstacle.N,
				Cost: math.Inf(0), // just pretend these all are total obstacles for now
			})
		}
	}
	return obstacles
}
