package formationpredicates


import "squaredance/dancer"
import "squaredance/geometry"

// Near returns true if the two dancers are near each other.
func Near(dancer1, dancer2 dancer.Dancer) bool {
	// *** Should we add a bit of fudge?
	return dancer1.Position().Distance(dancer2.Position()) <= geometry.CoupleDistance
}
