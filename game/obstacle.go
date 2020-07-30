package game

import "fmt"

// ObstacleType is an enum.
type ObstacleType int

//go:generate stringer -type=ObstacleType

// ParseObstacleType parses an ObstacleType from a string.
func ParseObstacleType(s string) (ObstacleType, error) {
	for i := 0; i < len(_ObstacleType_index)-1; i++ {
		if ObstacleType(i).String() == s+"Obstacle" {
			return ObstacleType(i), nil
		}
	}
	return NonObstacle, fmt.Errorf("unrecognised obstacle %s", s)
}

// ObstacleTypes represent the thing that is the obstacle. These might be static
// obstacle types, like a Tree or Rock, or another Character could be an Obstacle
// too.
const (
	NonObstacle ObstacleType = iota
	CharacterObstacle
	CrevasseObstacle
	TreeObstacle
	DeepWaterObstacle
	MudObstacle
)

// Obstacle is a Component that blocks a Hex.
type Obstacle struct {
	M, N int

	ObstacleType ObstacleType
}

// Type of the Component.
func (o *Obstacle) Type() string {
	return "Obstacle"
}
