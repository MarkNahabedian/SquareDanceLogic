package reasoning

import "squaredance/dancer"


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

// InFrontOf returns trur if dancer2 is in front of dancer1, that is,
// dancer1 is facing dancer2.
func InFrontOf(dancer1, dancer2 dancer.Dancer) bool {
	return dancer1.Direction().Equal(
		dancer1.Position().Direction(dancer2.Position()))
}

// Behind returns true if dancer2 is behind dancer1.
func Behind(dancer1, dancer2 dancer.Dancer) bool {
	return dancer1.Direction().QuarterRight().QuarterRight().Equal(
		dancer1.Position().Direction(dancer2.Position()))
}
