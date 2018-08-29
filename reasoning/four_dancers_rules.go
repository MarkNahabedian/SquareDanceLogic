// Definitions and rules about all four dancer square dance formations.
package reasoning

import "fmt"
import "goshua/rete"
import "squaredance/dancer"
import "squaredance/geometry"

type FacingCouples interface {
	Formation
	FacingCouples()
	Couple1() Couple
	Couple2() Couple
	Facing1() FaceToFace  // redundant
	Facing2() FaceToFace  // redundant
	// Roles:
	Beaus() dancer.Dancers    // no-slot
	Belles() dancer.Dancers   // no-slot
	Leaders() dancer.Dancers    // no-slot
	Trailers() dancer.Dancers    // no-slot
}

func (f *implFacingCouples) String() string {
	return fmt.Sprintf("FacingCouples(%s, %s)", f.Couple1(), f.Couple2())
}

func (f *implFacingCouples) Beaus() dancer.Dancers {
	return dancer.Union(f.Couple1().Beaus(), f.Couple2().Beaus())
}

func (f *implFacingCouples) Belles() dancer.Dancers {
	return dancer.Union(f.Couple1().Belles(), f.Couple2().Belles())
}

func (f *implFacingCouples) Leaders() dancer.Dancers {
	return dancer.Union(f.Facing1().Leaders(), f.Facing2().Leaders())
}

func (f *implFacingCouples) Trailers() dancer.Dancers {
	return dancer.Union(f.Facing1().Trailers(), f.Facing2().Trailers())
}

func rule_FacingCouples(node rete.Node, couple1, couple2 Couple, facing1, facing2 FaceToFace) {
	if !OrderedDancers(couple1.Beau(), couple2.Beau()) {
		return
	}
	// The same FaceToFace formation will come in as both facing1
	// and facing2.  These will be de-duped based on their relationship to
	// couple1 and couple2.
	if !HasDancers(facing1, couple1.Beau(), couple2.Belle()) {
		return
	}
	if !HasDancers(facing2, couple2.Beau(), couple1.Belle()) {
		return
	}
	node.Emit(FacingCouples(&implFacingCouples{
		couple1: couple1,
		couple2: couple2,
		facing1: facing1,
		facing2: facing2,
	}))
}


type TandemCouples interface {
	Formation
	TandemCouples()
	Couple1() Couple
	Couple2() Couple
	BeausTandem() Tandem   // redundant
	BellesTandem() Tandem  // redundant
	// Roles
	Beaus() dancer.Dancers       // no-slot
	Belles() dancer.Dancers      // no-slot
	Leaders() dancer.Dancers     // no-slot
	Trailers()  dancer.Dancers   // no-slot
}

func (f *implTandemCouples) String() string {
	return fmt.Sprintf("TandemCouples(%s, %s)", f.Couple1(), f.Couple2())
}

func (f *implTandemCouples) Beaus() dancer.Dancers {
	return f.BeausTandem().Dancers()
}

func (f *implTandemCouples) Belles() dancer.Dancers {
	return f.BellesTandem().Dancers()
}

func (f *implTandemCouples) Leaders() dancer.Dancers {
	return dancer.Union(f.BeausTandem().Leaders(), f.BellesTandem().Leaders())
}

func (f *implTandemCouples) Trailers() dancer.Dancers {
	return dancer.Union(f.BeausTandem().Trailers(), f.BellesTandem().Trailers())
}

func rule_TandemCouples(node rete.Node, couple1, couple2 Couple, tandem1, tandem2 Tandem) {
	if !OrderedDancers(couple1.Beau(), couple2.Beau()) {
		return
	}
	// The same Tandem formation will come in as both tandem1 and tandem2.
	// These will be de-duped based on their relationship to
	// couple1 and couple2.
	if !HasDancers(tandem1, couple1.Beau(), couple2.Beau()) {
		return
	}
	if !HasDancers(tandem2, couple1.Belle(), couple2.Belle()) {
		return
	}
	node.Emit(TandemCouples(&implTandemCouples{
		couple1: couple1,
		couple2: couple2,
		beaustandem: tandem1,
		bellestandem: tandem2,
	}))
}


type BackToBackCouples interface {
	Formation
	BackToBackCouples()
	Couple1() Couple
	Couple2() Couple
	BackToBack1() BackToBack  // redundant
	BackToBack2() BackToBack  // redundant
	// Roles:
	Beaus() dancer.Dancers     // no-slot
	Belles() dancer.Dancers    // no-slot
	Leaders()dancer.Dancers    // no-slot
	Trailers() dancer.Dancers  // no-slot 
}

func (f *implBackToBackCouples) String() string {
	return fmt.Sprintf("BackToBackCouples(%s, %s)", f.Couple1(), f.Couple2())
}

func (f *implBackToBackCouples) Beaus() dancer.Dancers {
	return dancer.Union(f.Couple1().Beaus(), f.Couple2().Beaus())
}

func (f *implBackToBackCouples) Belles() dancer.Dancers {
	return dancer.Union(f.Couple1().Belles(), f.Couple2().Belles())
}

func (f *implBackToBackCouples) Leaders () dancer.Dancers {
	return f.Dancers()
}

func (f *implBackToBackCouples) Trailers() dancer.Dancers {
	return dancer.Dancers{}
}


func rule_BackToBackCouples(node rete.Node, couple1, couple2 Couple, bb1, bb2 BackToBack) {
	if !OrderedDancers(couple1.Beau(), couple2.Beau()) {
		return
	}
	// The same BackToBack formation will come in as both bb1 and bb2.
	// These will be de-duped based on their relationship to
	// couple1 and couple2.
	if !HasDancers(bb1, couple1.Beau(), couple2.Belle()) {
		return
	}
	if !HasDancers(bb2, couple2.Beau(), couple1.Belle()) {
		return
	}
	node.Emit(BackToBackCouples(&implBackToBackCouples{
		couple1: couple1,
		couple2: couple2,
		backtoback1: bb1,
		backtoback2: bb2,
	}))
}


type BoxOfFour interface {
	Formation
	BoxOfFour()
	MiniWave1() MiniWave
	MiniWave2() MiniWave
	Tandem1() Tandem       // redundant
	Tandem2() Tandem       // redundant
	// Handedness:
	Handedness() Handedness      // no-slot
	// Roles:
	Beaus() dancer.Dancers      // no-slot
	Belles() dancer.Dancers     // no-slot
	Leaders() dancer.Dancers    // no-slot
	Trailers() dancer.Dancers   // no-slot
}

func (f *implBoxOfFour) String() string {
	return fmt.Sprintf("BoxOfFour(%s, %s, %s, %s, %s)",
		f.Handedness(),
		f.Tandem1().Leader(), f.Tandem1().Trailer(),
		f.Tandem2().Leader(), f.Tandem2().Trailer())
}

func (f *implBoxOfFour) Handedness() Handedness {
	return f.MiniWave1().Handedness()
}

func (f *implBoxOfFour) Beaus() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Beaus(), f.MiniWave2().Beaus())
}

func (f *implBoxOfFour) Belles() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Belles(), f.MiniWave2().Belles())
}

func (f *implBoxOfFour) Leaders() dancer.Dancers {
	return dancer.Union(f.Tandem1().Leaders(), f.Tandem2().Leaders())
}

func (f *implBoxOfFour) Trailers() dancer.Dancers {
	return dancer.Union(f.Tandem1().Trailers(), f.Tandem2().Trailers())
}

func rule_BoxOfFour(node rete.Node, mw1, mw2 MiniWave, tandem1, tandem2 Tandem) {
	if !tandem1.Direction().Opposite().Equal(tandem2.Direction()) {
		return
	}
	if mw1.HasDancer(tandem1.Leader()) {
		if !mw1.HasDancer(tandem2.Trailer()) {
			return
		}
		if !mw2.HasDancer(tandem1.Trailer()) {
			return
		}
		if !mw2.HasDancer(tandem2.Leader()) {
			return
		}
	} else if mw1.HasDancer(tandem1.Trailer()) {
		if !mw1.HasDancer(tandem2.Leader()) {
			return
		}
		if !mw2.HasDancer(tandem1.Leader()) {
			return
		}
		if !mw2.HasDancer(tandem2.Trailer()) {
			return
		}
	} else {
  		return
	}
	node.Emit(BoxOfFour(&implBoxOfFour{
		miniwave1: mw1,
		miniwave2: mw2,
		tandem1: tandem1,
		tandem2: tandem2,
	}))
}


type Star interface {
	Formation
	Star()
	MiniWave1() MiniWave
	MiniWave2() MiniWave
	// Handedness:
	Handedness() Handedness    // no-slot
	// Roles:
	Beaus() dancer.Dancers    // no-slot
	Belles() dancer.Dancers   // no-slot
}

func (f *implStar) String() string {
	return fmt.Sprintf("Star(%s, %s, %s, %s, %s)",
		f.Handedness(),
		f.MiniWave1().Dancer1(),
		f.MiniWave2().Dancer1(),
		f.MiniWave1().Dancer2(),
		f.MiniWave2().Dancer2())		
}

func (f *implStar) Handedness() Handedness {
	return f.MiniWave1().Handedness()
}

func (f *implStar) Beaus() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Beaus(), f.MiniWave2().Beaus())
}

func (f *implStar) Belles() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Belles(), f.MiniWave2().Belles())
}

func rule_Star(node rete.Node, mw1, mw2 MiniWave) {
	if !geometry.Center(dancer.Positions(mw1.Dancers()...)...).Equal(
			geometry.Center(dancer.Positions(mw2.Dancers()...)...)) {
		return
	}
	dir := mw1.Dancer1().Direction().QuarterLeft()
	if !(dir.Equal(mw2.Dancer1().Direction()) ||
		dir.Equal(mw2.Dancer2().Direction())) {
		return
	}
	node.Emit(Star(&implStar{
		miniwave1: mw1,
		miniwave2: mw2,
	}))
}


type LineOfFour interface {
	Formation
	LineOfFour()
	LeftCouple() Couple
	CenterCouple() Couple    // redudant
	RightCouple() Couple
	// Roles:
	Beaus() dancer.Dancers    // no-slot
	Belles() dancer.Dancers   // no-slot
	Centers() dancer.Dancers  // no-slot
	Ends() dancer.Dancers     // no-slot
}

func (f *implLineOfFour) String() string {
	return fmt.Sprintf("LineOfFour(%s, %s, %s, %s)",
		f.LeftCouple().Beau(),
		f.LeftCouple().Belle(),
		f.RightCouple().Beau(),
		f.RightCouple().Belle())
}

func (f *implLineOfFour) Beaus() dancer.Dancers {
	return dancer.Union(f.LeftCouple().Beaus(), f.RightCouple().Beaus())
}

func (f *implLineOfFour) Belles() dancer.Dancers {
	return dancer.Union(f.LeftCouple().Belles(), f.RightCouple().Belles())
}

func (f *implLineOfFour) Centers() dancer.Dancers {
	return f.CenterCouple().Dancers()
}

func (f *implLineOfFour) Ends() dancer.Dancers {
	return dancer.SetDifference(f.Dancers(), f.Centers())
}

func rule_LineOfFour(node rete.Node, c1, c2, c3 Couple) {
	if !(c1.Belle() == c2.Beau()) {
		return
	}
	if !(c2.Belle() == c3.Beau()) {
		return
	}
	node.Emit(LineOfFour(&implLineOfFour{
		leftcouple: c1,
		centercouple: c2,
		rightcouple: c3,
	}))
}


type WaveOfFour interface {
	Formation
	WaveOfFour()
	CenterMiniWave() MiniWave    // redundant
	MiniWave1() MiniWave
	MiniWave2() MiniWave
	// Handedness:
	Handedness() Handedness      // no-slot
	// Roles:
	Beaus() dancer.Dancers      // no-slots
	Belles() dancer.Dancers     // no-slots
	Centers() dancer.Dancers    // no-slots
	Ends() dancer.Dancers       // no-slots
}

func (f *implWaveOfFour) String() string {
	return fmt.Sprintf("WaveOfFour(%s, %s, %s, %s, %s)",
		f.Handedness(),
		f.MiniWave1().Dancer1(),
		f.MiniWave1().Dancer2(),
		f.MiniWave2().Dancer1(),
		f.MiniWave2().Dancer1())
}

func (f *implWaveOfFour) Handedness() Handedness {
	return f.MiniWave1().Handedness()
}

func (f *implWaveOfFour) Beaus() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Beaus(), f.MiniWave2().Beaus())
}

func (f *implWaveOfFour) Belles() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Belles(), f.MiniWave2().Belles())
}

func (f *implWaveOfFour) Centers() dancer.Dancers {
	return f.CenterMiniWave().Dancers()
}

func (f *implWaveOfFour) Ends() dancer.Dancers {
	return dancer.SetDifference(f.Dancers(), f.Centers())
}

func rule_WaveOfFour(node rete.Node, mw1, mw2, mw3 MiniWave) {
	if mw1.HasDancer(mw2.Dancer1()) {
		if !mw2.HasDancer(mw2.Dancer2()) {
			return
		}
	} else if mw1.HasDancer(mw2.Dancer2()) {
		if !mw2.HasDancer(mw2.Dancer1()) {
			return
		}
	}
	node.Emit(WaveOfFour(&implWaveOfFour{
		centerminiwave: mw2,
		miniwave1:mw1,
		miniwave2: mw3,
	}))
}


type TwoFacedLine interface {
	Formation
	TwoFacedLine()
	Couple1() Couple
	Couple2() Couple
	CenterMiniWave() MiniWave  // redundant
	// Handedness:
	Handedness() Handedness    // no-slot
	// Roles:
	Beaus() dancer.Dancers    // no-slots
	Belles() dancer.Dancers   // no-slots
	Centers() dancer.Dancers  // no-slots
	Ends() dancer.Dancers     // no-slots
}

func (f *implTwoFacedLine) String() string {
	return fmt.Sprintf("TwoFacedLine(%s, %s, %s, %s, %s)",
    	f.Handedness(),
		f.Couple1().Beau(),
		f.Couple1().Belle(),
		f.Couple2().Beau(),
		f.Couple2().Belle())
}

func (f *implTwoFacedLine) Handedness() Handedness {
	return f.CenterMiniWave().Handedness()
}

func (f *implTwoFacedLine) Beaus() dancer.Dancers {
	return dancer.Union(f.Couple1().Beaus(), f.Couple2().Beaus())
}

func (f *implTwoFacedLine) Belles() dancer.Dancers {
	return dancer.Union(f.Couple1().Belles(), f.Couple2().Belles())
}

func (f *implTwoFacedLine) Centers() dancer.Dancers {
	return f.CenterMiniWave().Dancers()
}

func (f *implTwoFacedLine) Ends() dancer.Dancers {
	return dancer.SetDifference(f.Dancers(), f.Centers())
}

func rule_TwoFacedLine(node rete.Node, c1, c2 Couple, mw MiniWave) {
	if !mw.HasDancer(c1.Beau()) {
		return
	}
	if !mw.HasDancer(c2.Belle()) {
		return
	}
	node.Emit(TwoFacedLine(&implTwoFacedLine{
		couple1: c1,
		couple2: c2,
		centerminiwave: mw,
	}))
}

// Diamonds

// Zs
