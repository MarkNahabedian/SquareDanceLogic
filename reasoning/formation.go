package reasoning

import "squaredance/dancer"
import "squaredance/geometry"

// Formation represents a square dance formation.
type Formation interface {
	// NumberOfDancers returns the number of dancers in the Formation.
	NumberOfDancers() int

	// Dancers returns a slice containing the Dancers in the Formation.
	Dancers() []dancer.Dancer

	// HasDancer returns true of the specified Dancer is in the Formation.
	HasDancer(dancer.Dancer) bool
}

// HasDancers returns true if f contains all of the specified Dancers.
func HasDancers(f Formation, dancers ...dancer.Dancer) bool {
	for _, d := range dancers {
		if !f.HasDancer(d) {
			return false
		}
	}
	return true
}

// OrderedDancers returns true if all of the Dancers are in
// ascending order by Ordinal.
func OrderedDancers(dancers ...dancer.Dancer) bool {
	for i := 0; i < len(dancers) - 1; i++ {
		if dancers[i].Ordinal() >= dancers[i+1].Ordinal() {
			return false
		}
	}
	return true
}

func (t *implTandem) Direction() geometry.Direction {
	return t.Leader().Direction()
}
