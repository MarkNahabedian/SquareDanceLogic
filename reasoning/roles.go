// For each role we define a function named by that role which takes a
// formation as argument and retuirns a sloice of Dancer.
package reasoning

import "squaredance/geometry"
import "squaredance/dancer"

type Role func(Formation) []dancer.Dancer

var Roles []Role

// Some roles depend on the formation.  For these, the relevant
// formation defines a method of the same name as the role function which
// the role function calls for its result.  For these, we also define an
// interface "Has<role>" which has that method as its only member.

// Other roles depend only on the dancers themselves.

func OriginalHead(f Formation) []dancer.Dancer {
	result := []dancer.Dancer{}
	for _, d := range f.Dancers() {
		cn := d.CoupleNumber()
		if cn == 1 || cn == 3 {
			result = append(result, d)
		}
	}
	return result
}

func OriginalSide(f Formation) []dancer.Dancer {
	result := []dancer.Dancer{}
	for _, d := range f.Dancers() {
		cn := d.CoupleNumber()
		if cn == 2 || cn == 4 {
			result = append(result, d)
		}
	}
	return result
}

func CurrentHead(f Formation) []dancer.Dancer {
	result := []dancer.Dancer{}
	for _, d := range f.Dancers() {
		dd := d.Direction()
		if dd.Equal(geometry.Direction(0)) || dd.Equal(geometry.Direction(0).Opposite()) {
			result = append(result, d)
		}
	}
	return result
}

func CurrentSide(f Formation) []dancer.Dancer {
	result := []dancer.Dancer{}
	for _, d := range f.Dancers() {
		dd := d.Direction()
		if dd.Equal(geometry.Direction(0).QuarterRight()) || dd.Equal(geometry.Direction(0).QuarterLeft()) {
			result = append(result, d)
		}
	}
	return result
}

func init() {
	Roles = append(Roles, OriginalHead, OriginalSide, CurrentHead, CurrentSide)
}