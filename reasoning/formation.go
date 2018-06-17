package reasoning

import "squaredance/dancer"

// Formation represents a square dance formation.
type Formation interface {
	// NumberOfDancers returns the number of dancers in the Formation.
	NumberOfDancers() int

	// Dancers returns a slice containing the Dancers in the Formation.
	Dancers() []dancer.Dancer

	// HasDancer returns true of the specified Dancer is in the Formation.
	HasDancer(dancer.Dancer) bool
}

