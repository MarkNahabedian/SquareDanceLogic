package formationpredicates


import "squaredance/dancer"
// import "squaredance/geometry"


// LeftOf returns true if dancer2 is to the left of Dancer1.
// Dancer1's direction is relevant to this determination but
// Dancer2's direction is not.
func LeftOf(dancer1, dancer2 dancer.Dancer) bool {
	return dancer1.Direction().QuarterLeft().Equal(
    	dancer1.Position().Direction(dancer2.Position()))
}

// RightOf returns true if dancer2 is to the left of Dancer1.
// Dancer1's direction is relevant to this determination but
// Dancer2's direction is not.
func RightOf(dancer1, dancer2 dancer.Dancer) bool {
	return dancer1.Direction().QuarterRight().Equal(
    	dancer1.Position().Direction(dancer2.Position()))
}
