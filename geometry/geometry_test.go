package geometry

import "testing"

func TestDirection(t *testing.T) {
	d1 := FullCircle.DivideBy(4)

	if !d1.Opposite().Add(d1).Equal(FullCircle) {
		t.Errorf("Direction inverse/add failed.")
	}

	if !FullCircle.Add(FullCircle).Equal(FullCircle.Subtract(FullCircle)) {
		t.Errorf("Direction add/subtract failed.")
	}

	if !FullCircle.DivideBy(4.0).MultiplyBy(4.0).Equal(FullCircle) {
		t.Errorf("Direction multiply/divide failed.")
	}
}

func TestDownLeft(t *testing.T) {
	if float64(Down0.Add(Down1).Subtract(Down1)) != 0 {
		t.Errorf("Down arithmetic failed")
	}
	if float64(Left0.Add(Left1).Subtract(Left1)) != 0 {
		t.Errorf("Left arithmetic failed")
	}
}

func TestPosition(t *testing.T) {
	d1 := FullCircle.DivideBy(4)
	origin := NewPositionDownLeft(0.0, 0.0)
	p := NewPosition(d1, 2.0)
	if p.Distance(origin) != 2.0 {
		t.Errorf("NewPosition or Distaqnce failed.")
	}
	if !origin.Direction(p).Equal(d1) {
		t.Errorf("NewPosition or Direction failed.")
	}
	if !(p.Direction(origin).Equal(FullCircle.DivideBy(4).Opposite())) {
		t.Errorf("Right isoscelese triangle failed Distance.")
	}
}

func TestCenter(t *testing.T) {
	expectedCenter := NewPosition(0.0, 0.0)
	if !expectedCenter.Equal(expectedCenter) {
		t.Errorf("Position.Equal failed")
	}
	c := Center([]Position{
		expectedCenter.Add(NewPosition(0.25, 1.0)),
		expectedCenter.Add(NewPosition(0.75, 1.0)),
	})
	if !expectedCenter.Equal(c) {
		t.Errorf("Position Center failed: %v", c)
	}
}
