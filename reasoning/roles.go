package reasoning

import "squaredance/geometry"
import "squaredance/dancer"


type Role interface {
	// Name returns the name of the Role.
	// By convention, Role names are singular, so "Head" rather than "Heads".
	Name() string

	// MeaningfulTo returns true if the Formation applies some meaning to the Role.
	MeaningfulTo(Formation) bool

	// Dancers returns the dancers from the Formation that fit this Role.
	Dancers(Formation) []dancer.Dancer
}

var Roles []Role

func LookupRole(name string) Role {
	for _, r := range Roles {
		if name == r.Name() {
			return r
		}
	}
	return nil
}


// ubiquitousRole is used to implement any Role that is relevant to all formations.
type ubiquitousRole struct {
	name string
	filter func(dancer.Dancer) bool
}

func (r *ubiquitousRole) Name() string { return r.name }

func (r *ubiquitousRole) MeaningfulTo(Formation) bool { return true }

func (r *ubiquitousRole) Dancers(f Formation) []dancer.Dancer {
	result := []dancer.Dancer{}
	for _, d := range f.Dancers() {
		if r.filter(d) {
			result = append(result, d)
		}
	}
	return result
}

func init() {
	Roles = append(Roles,
		&ubiquitousRole {
			name: "OriginalHeads",
			filter: func(d dancer.Dancer) bool {
				cn := d.CoupleNumber()
				return cn & 1 == 1
			},
		},
		&ubiquitousRole {
			name: "OriginalSides",
			filter: func(d dancer.Dancer) bool {
				cn := d.CoupleNumber()
				return cn & 1 == 0
			},
		},
		&ubiquitousRole {
			name: "CurrentHeads",
			filter: func(d dancer.Dancer) bool {
				dd := d.Direction()
				// Only works for 8 dancer sets.
				return (dd.Equal(geometry.Direction(0)) || dd.Equal(geometry.Direction(0).Opposite()))
			},
		},
		&ubiquitousRole {
			name: "CurrentSides",
			filter: func(d dancer.Dancer) bool {
				dd := d.Direction()
				// Only works for 8 dancer sets.
				return (dd.Equal(geometry.Direction(0).QuarterRight()) ||
					dd.Equal(geometry.Direction(0).QuarterLeft()))
			},
		},
	)
}


// Some roles depend on the formation.  For these, the relevant
// formation defines a method of the same name as the role which
// the Role's Dancers method calls for its result.  For these, we
// also define an interface "Has<role>" which has that method as
// its only member.
