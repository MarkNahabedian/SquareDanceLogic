package action

import "testing"
import "squaredance/dancer"
import "squaredance/geometry"


func TestBreating(t *testing.T) {
	dancers := dancer.MakeSomeDancers(4)
	left0 := geometry.Left0
	left1 := geometry.Left(0.7)
	down0 := geometry.Down0
	down1 := geometry.Down(0.8)
	dancers[0].Move(geometry.Position{ Left: left0, Down: down0 },
		geometry.Direction0)
	dancers[1].Move(geometry.Position{ Left: left1, Down: down0 },
		geometry.Direction0)
	dancers[2].Move(geometry.Position{ Left: left0, Down: down1},
		geometry.Direction0)
	dancers[3].Move(geometry.Position{ Left: left1, Down: down1},
		geometry.Direction0)
	center_before := dancers.Center()
	Breathe(dancers)
	center_after := dancers.Center()
	if !center_before.Equal(center_after) {
		t.Errorf("Center moved from %v to %v",
			center_before, center_after)
	}
	check_distance := func(dancer1, dancer2 dancer.Dancer) {
		if d := dancer.Distance(dancer1, dancer2); d < geometry.CoupleDistance {
			t.Errorf("Too close: %f, %v, %v", d, dancer1, dancer2)
		}
	}
	check_distance(dancers[0], dancers[1])
	check_distance(dancers[2], dancers[3])
	check_distance(dancers[0], dancers[2])
	check_distance(dancers[1], dancers[3])
}

