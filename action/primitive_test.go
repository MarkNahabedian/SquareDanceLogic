package action

import "testing"
import "squaredance/geometry"
import "squaredance/dancer"
// import "squaredance/reasoning"

func loneTestDancer() dancer.Dancers {
	dancers := dancer.MakeSomeDancers(1)
	dancer := dancers[0]
	here := geometry.NewPositionDownLeft(0, 0)
	facing := geometry.Direction0.QuarterRight()
	dancer.Move(here, facing)
	return dancers
}

func TestQuarterRight(t *testing.T) {
	dancers := loneTestDancer()
	dancer := dancers[0]
	p := dancer.Position()
	d := dancer.Direction()
	fa := FindAction("QuarterRight").GetFormationActionFor(dancer)
	if fa == nil {
		t.Fatalf("GetFormationActionFor returned nil")
	}
	fa.DoIt(dancer)
	if !p.Equal(dancer.Position()) {
		t.Errorf("Position changed during QuarterRight")
	}
	if want, got := d.QuarterRight(), dancer.Direction(); !want.Equal(got) {
		t.Errorf("Wrong direction, want: %#v, got %#v", want, got)
	}
}

