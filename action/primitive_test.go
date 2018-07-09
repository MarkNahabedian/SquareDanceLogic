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
