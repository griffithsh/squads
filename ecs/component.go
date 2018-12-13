package ecs

// Component is a thing that stores data and can be interacted with by a System.
type Component interface {
	Type() string
}
