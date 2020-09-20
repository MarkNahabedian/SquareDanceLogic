package action

import "testing"
import "reflect"
import "squaredance/geometry"
import "squaredance/dancer"
import "squaredance/timeline"
import "squaredance/reasoning"


func TestWritePrimitiveCatalog(t *testing.T) {
	WriteCatalog(Primitive)
}


// showHistory writes the position and directiion of each Dancer
// over time to standard output.
func showHistory(tl timeline.Timeline, t *testing.T) {
	for _, d := range tl.Dancers() {
		t.Logf("\nDancer %s\n", d)
		for _, s := range tl.FindSnapshots(d, -1, tl.MostRecent() + 1) {
			t.Logf("    %3d  %s  %v\n", s.Time(), s.Position(), s.Direction())
		}
	}
}


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

func TestQuarterLeft(t *testing.T) {
	dancers := loneTestDancer()
	dancer := dancers[0]
	p := dancer.Position()
	d := dancer.Direction()
	fa := FindAction("QuarterLeft").GetFormationActionFor(dancer)
	if fa == nil {
		t.Fatalf("GetFormationActionFor returned nil")
	}
	fa.DoIt(dancer)
	if !p.Equal(dancer.Position()) {
		t.Errorf("Position changed during QuarterLeft")
	}
	if want, got := d.QuarterLeft(), dancer.Direction(); !want.Equal(got) {
		t.Errorf("Wrong direction, want: %#v, got %#v", want, got)
	}
}

func TestAboutFace(t *testing.T) {
	dancers := loneTestDancer()
	dancer := dancers[0]
	p := dancer.Position()
	d := dancer.Direction()
	fa := FindAction("AboutFace").GetFormationActionFor(dancer)
	if fa == nil {
		t.Fatalf("GetFormationActionFor returned nil")
	}
	fa.DoIt(dancer)
	if !p.Equal(dancer.Position()) {
		t.Errorf("Position changed during AboutFace")
	}
	if want, got := d.Opposite(), dancer.Direction(); !want.Equal(got) {
		t.Errorf("Wrong direction, want: %#v, got %#v", want, got)
	}
}

func TestMeet(t *testing.T) {
	set := dancer.NewSquaredSet(4)
	tl := timeline.NewTimeline(set.Dancers())
	tl.MakeSnapshot(0)
	// Get FaceToFace formations:
	formations := reasoning.FindFormations(set.Dancers(), reflect.TypeOf(func(t reasoning.FaceToFace){}).In(0))
	// Restrict to heads:
	headsff := []reasoning.FaceToFace{}
	for _, ff := range formations {
		heads := reasoning.LookupRole("OriginalHeads").Dancers(ff.Dancers())
		if len(heads) == 2 {
			headsff = append(headsff, ff.(reasoning.FaceToFace))
		}
	}
	if l := len(headsff); l != 2 {
		t.Log(headsff)
		t.Fatalf("Expected 2 Head FaceToFace formations, got %d.", l)
	}
	for _, ff := range headsff {
		fa := FindAction("Meet").GetFormationActionFor(ff)
		if fa == nil {
			t.Fatalf("GetFormationActionFor returned nil")
		}
		fa.DoIt(ff)
	}
	tl.MakeSnapshot(1)
	// Dancer facing directions are unchanged
	for _, d := range set.Dancers() {
		if want, got := tl.FindSnapshot(d, 0).Direction(), tl.FindSnapshot(d, 1).Direction(); got != want {
			t.Errorf("Dancer %s direction changed: want: %v, got: %v",
				d, want, got)
		}
	}
	// Head dancers are close together now
	for _, ff := range headsff {
		distance := tl.FindSnapshot(ff.Dancer1(), 1).Position().Distance(
			tl.FindSnapshot(ff.Dancer2(), 1).Position())
		if distance > geometry.CoupleDistance {
			t.Errorf("Dancers are too far apart: %s, %f should be <= %f.",
				ff, distance, geometry.CoupleDistance)
		}
	}
}

func mini_wave(handedness reasoning.Handedness) reasoning.Formation {
	dancers := dancer.MakeSomeDancers(2)
	var dir geometry.Direction
	switch handedness {
	case reasoning.RightHanded:
		dir = geometry.Direction0
	case reasoning.LeftHanded:
		dir = geometry.Direction0.Opposite()
	}
	dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
		dir.Opposite())
	dancers[1].Move(geometry.Position{ Left: geometry.Left1, Down: geometry.Down0 },
		dir)
	r := reasoning.FindFormations(dancers, reflect.TypeOf(func(t reasoning.MiniWave){}).In(0))
	return r[0]
}

func TestForwardLeft(t *testing.T) {
	// Start with FaceToFace dancers.  End in RightHand MiniWave
	dancers := reasoning.MakeSampleFormation(reasoning.AllFormationTypes["FaceToFace"])
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("ForwardLeft").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := reasoning.FindFormations(dancers.Dancers(), reflect.TypeOf(func(t reasoning.MiniWave){}).In(0))
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wrong number of MiniWaves: got %d, want %d.", got, want)
		return
	}
	mw := dancers2[0]
	if want, got := reasoning.RightHanded, mw.(reasoning.MiniWave).Handedness(); got != want {
		t.Errorf("Wrong handedness: want %v, got %v.", want, got)
	}
	showHistory(tl, t)
}

func TestForwardRight(t *testing.T) {
	// Start with FaceToFace dancers.  End in RightHand MiniWave
	dancers := reasoning.MakeSampleFormation(reasoning.AllFormationTypes["FaceToFace"])
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("ForwardRight").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := reasoning.FindFormations(dancers.Dancers(),
		reflect.TypeOf(func(t reasoning.MiniWave){}).In(0))
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wrong number of MiniWaves: got %d, want %d.", got, want)
		return
	}
	mw := dancers2[0]
	if want, got := reasoning.LeftHanded, mw.(reasoning.MiniWave).Handedness(); got != want {
		t.Errorf("Wrong handedness: want %v, got %v.", want, got)
	}
	showHistory(tl, t)
}

func TestBackwardLeft(t *testing.T) {
	t.Logf("TestBackwardLeft\n")
	// Start with BackToBack dancers.  End in RightHand MiniWave
	dancers := reasoning.MakeSampleFormation(reasoning.AllFormationTypes["BackToBack"])
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("BackwardLeft").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := reasoning.FindFormations(dancers.Dancers(),
		reflect.TypeOf(func(t reasoning.MiniWave){}).In(0))
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wrong number of MiniWaves: got %d, want %d.", got, want)
		return
	}
	mw := dancers2[0]
	if want, got := reasoning.RightHanded, mw.(reasoning.MiniWave).Handedness(); got != want {
		t.Errorf("Wrong handedness: want %v, got %v.", want, got)
	}
	showHistory(tl, t)
}

func TestBackwardRight(t *testing.T) {
	t.Logf("TestBackwardRight\n")
	// Start with BackToBack dancers.  End in RightHand MiniWave
	dancers := reasoning.MakeSampleFormation(reasoning.AllFormationTypes["BackToBack"])
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("BackwardRight").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := reasoning.FindFormations(dancers.Dancers(),
		reflect.TypeOf(func(t reasoning.MiniWave){}).In(0))
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wrong number of MiniWaves: got %d, want %d.", got, want)
		return
	}
	mw := dancers2[0]
	if want, got := reasoning.LeftHanded, mw.(reasoning.MiniWave).Handedness(); got != want {
		t.Errorf("Wrong handedness: want %v, got %v.", want, got)
	}
	showHistory(tl, t)
}

func TestBackToFaceRight(t *testing.T) {
	t.Logf("TestBackToFaceRight\n")
	// Start with a RightHanded MiniWave.  End in FaceToFace dancers.
	dancers := mini_wave(reasoning.RightHanded).(reasoning.MiniWave)
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("BackToFace").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := reasoning.FindFormations(dancers.Dancers(),
		reflect.TypeOf(func(t reasoning.FaceToFace){}).In(0))
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wrong number of FaceToFace formations: got %d, want %d.", got, want)
		return
	}
	showHistory(tl, t)
}

func TestBackToFaceLeft(t *testing.T) {
	t.Logf("TestBackToFaceLeft\n")
	// Start with a LeftHanded MiniWave.  End in FaceToFace dancers.
	dancers := mini_wave(reasoning.LeftHanded).(reasoning.MiniWave)
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("BackToFace").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := reasoning.FindFormations(dancers.Dancers(),
		reflect.TypeOf(func(t reasoning.FaceToFace){}).In(0))
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wrong number of FaceToFace formations: got %d, want %d.", got, want)
		return
	}
	showHistory(tl, t)
}

