package action

import "testing"
import "goshua/rete"
import "goshua/rete/rule_compiler/runtime"
import "squaredance/geometry"
import "squaredance/dancer"
import "squaredance/timeline"
import "squaredance/reasoning"


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


func get_formation(dancers dancer.Dancers, formation_name string) []reasoning.Formation {
	root_node := rete.MakeRootNode()
	for _, rule := range runtime.AllRules {
		rule.Inserter()(root_node)
	}
	result := []reasoning.Formation{}
	rete.Walk(root_node, func(n rete.Node) {
		ttn, ok := n.(*rete.TypeTestNode)
		if !ok {
			return
		}
		if ttn.TypeName() != formation_name {
			return
		}
		rete.Connect(n, rete.MakeActionNode(func(item interface{}) {
			result = append(result, item.(reasoning.Formation))
		}))
	})
	for _, d := range dancers {
		root_node.Receive(d)
	}
	return result
}


func TestMeet(t *testing.T) {
	set := dancer.NewSquaredSet(4)
	tl := timeline.NewTimeline(set.Dancers())
	tl.MakeSnapshot(0)
	root_node := rete.MakeRootNode()
	for _, rule := range runtime.AllRules {
		rule.Inserter()(root_node)
	}
	headsff := []reasoning.FaceToFace{}
	rete.Walk(root_node, func(n rete.Node) {
		ttn, ok := n.(*rete.TypeTestNode)
		if !ok {
			return
		}
		if ttn.TypeName() != "FaceToFace" {
			return
		}
		rete.Connect(n, rete.MakeActionNode(func(item interface{}) {
			ff := item.(reasoning.FaceToFace)
			heads := reasoning.LookupRole("OriginalHeads").Dancers(ff.Dancers())
			if len(heads) == 2 {
				headsff = append(headsff, ff)
			}
		}))
	})
	root_node.Receive(set)
	if l := len(headsff); l != 2 {
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
	for _, ff1 := range headsff {
		ff := ff1.(reasoning.FaceToFace)
		distance := tl.FindSnapshot(ff.Dancer1(), 1).Position().Distance(
			tl.FindSnapshot(ff.Dancer2(), 1).Position())
		if distance > geometry.CoupleDistance {
			t.Errorf("Dancers are too far apart: %s, %f should be <= %f.",
				ff, distance, geometry.CoupleDistance)
		}
	}
}

func face_to_face() reasoning.Formation {
	dancers := dancer.MakeSomeDancers(2)
	dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0},
		geometry.Direction0)
	dancers[1].Move(geometry.Position{Left: geometry.Left0, Down: geometry.Down1},
		geometry.FullCircle / 2)
	r := get_formation(dancers, "FaceToFace")
	return r[0]
}

func back_to_back() reasoning.Formation {
	dancers := dancer.MakeSomeDancers(2)
	dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0},
		geometry.FullCircle / 2)
	dancers[1].Move(geometry.Position{Left: geometry.Left0, Down: geometry.Down1},
		geometry.Direction0)
	r := get_formation(dancers, "BackToBack")
	return r[0]
}

func TestForwardLeft(t *testing.T) {
	// Start with FaceToFace dancers.  End in RightHand MiniWave
	dancers := face_to_face().(reasoning.FaceToFace)
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("ForwardLeft").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := get_formation(dancers.Dancers(), "MiniWave")
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wri=on number of MiniWaves: got %d, want %d.", got, want)
		return
	}
	mw := dancers2[0]
	if want, got := reasoning.RightHanded, mw.(reasoning.MiniWave).Handedness(); got != want {
		t.Errorf("Wrong handedness: want %v, got %v.", want, got)
	}
	timeline.ShowHistory(tl)
}

func TestForwardRight(t *testing.T) {
	// Start with FaceToFace dancers.  End in RightHand MiniWave
	dancers := face_to_face().(reasoning.FaceToFace)
	tl := timeline.NewTimeline(dancers.Dancers())
	tl.MakeSnapshot(0)
	fa := FindAction("ForwardRight").GetFormationActionFor(dancers)
	if fa == nil {
		t.Fatalf("GetFormationActionFor did not find a FormationAction for %#v", dancers)
	}
	fa.DoIt(dancers)
	tl.MakeSnapshot(1)
	dancers2 := get_formation(dancers.Dancers(), "MiniWave")
	if want, got := 1, len(dancers2); got != want {
		t.Errorf("Wri=on number of MiniWaves: got %d, want %d.", got, want)
		return
	}
	mw := dancers2[0]
	if want, got := reasoning.LeftHanded, mw.(reasoning.MiniWave).Handedness(); got != want {
		t.Errorf("Wrong handedness: want %v, got %v.", want, got)
	}
	timeline.ShowHistory(tl)
}


// "ForwardRight"

// "BackwardLeft"

// "BackwardRight"

// "BackToFace"

