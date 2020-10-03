// Definitions and rules about all four dancer square dance formations.
package reasoning

import "fmt"
import "goshua/rete"
import "squaredance/dancer"
import "squaredance/geometry"


func swap_positions(dancer1, dancer2 dancer.Dancer) {
	d1pos := dancer1.Position()
	dancer1.Move(dancer2.Position(), dancer1.Direction())
	dancer2.Move(d1pos, dancer2.Direction())
}

type FacingCouples interface {
	Formation
	FacingCouples()
	Couple1() Couple      // defimpl:"read couple1"  fe:"dancers"
	Couple2() Couple      // defimpl:"read couple2"  fe:"dancers"
	Facing1() FaceToFace  // defimpl:"read facing1"
	Facing2() FaceToFace  // defimpl:"read facing2"
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	Leaders() dancer.Dancers
	Trailers() dancer.Dancers
}

func (f *FacingCouplesImpl) String() string {
	return fmt.Sprintf("FacingCouples(%s, %s)", f.Couple1(), f.Couple2())
}

func (f *FacingCouplesImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.Couple1().Beaus(), f.Couple2().Beaus())
}

func (f *FacingCouplesImpl) Belles() dancer.Dancers {
	return dancer.Union(f.Couple1().Belles(), f.Couple2().Belles())
}

func (f *FacingCouplesImpl) Leaders() dancer.Dancers {
	return dancer.Union(f.Facing1().Leaders(), f.Facing2().Leaders())
}

func (f *FacingCouplesImpl) Trailers() dancer.Dancers {
	return dancer.Union(f.Facing1().Trailers(), f.Facing2().Trailers())
}

func make_FacingCouples_sample() Formation {
	couple1 := make_Couple_sample().(*CoupleImpl)
	couple2 := make_Couple_sample().(*CoupleImpl)
	down := geometry.NewPositionDownLeft(geometry.Down1, geometry.Left0)
	couple2.Beau().MoveBy(down).Rotate(geometry.Direction2)
	couple2.Belle().MoveBy(down).Rotate(geometry.Direction2)
	swap_positions(couple2.Beau(),couple2.Belle())
	dancer.Reorder(couple1.Beau(), couple1.Belle(), couple2.Beau(), couple2.Belle())
	sample := FacingCouples(&FacingCouplesImpl {
		couple1: couple1,
		couple2: couple2,
		facing1: &FaceToFaceImpl {
			dancer1: couple1.Beau(),
			dancer2: couple2.Belle(),
		},
		facing2: &FaceToFaceImpl {
			dancer1: couple1.Belle(),
			dancer2: couple2.Beau(),
		},
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_FacingCouples_sample)
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
	node.Emit(FacingCouples(&FacingCouplesImpl{
		couple1: couple1,
		couple2: couple2,
		facing1: facing1,
		facing2: facing2,
	}))
}


type TandemCouples interface {
	Formation
	TandemCouples()
	// Should we call these LeadingCouple and TrailingCouple?
	Couple1() Couple        // defimpl:"read couple1" fe:"dancers"
	Couple2() Couple        // defimpl:"read couple2" fe:"dancers"
	BeausTandem() Tandem   // defimpl:"read beaustandem"
	BellesTandem() Tandem  // defimpl:"read bellestandem"
	// Roles
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	Leaders() dancer.Dancers
	Trailers()  dancer.Dancers
}

func (f *TandemCouplesImpl) String() string {
	return fmt.Sprintf("TandemCouples(%s, %s)", f.Couple1(), f.Couple2())
}

func (f *TandemCouplesImpl) Beaus() dancer.Dancers {
	return f.BeausTandem().Dancers()
}

func (f *TandemCouplesImpl) Belles() dancer.Dancers {
	return f.BellesTandem().Dancers()
}

func (f *TandemCouplesImpl) Leaders() dancer.Dancers {
	return dancer.Union(f.BeausTandem().Leaders(), f.BellesTandem().Leaders())
}

func (f *TandemCouplesImpl) Trailers() dancer.Dancers {
	return dancer.Union(f.BeausTandem().Trailers(), f.BellesTandem().Trailers())
}

func make_TandemCouples_sample() Formation {
	couple1 := make_Couple_sample().(*CoupleImpl)
	couple2 := make_Couple_sample().(*CoupleImpl)
	down := geometry.NewPositionDownLeft(geometry.Down1, geometry.Left0)
	couple2.Beau().MoveBy(down)
	couple2.Belle().MoveBy(down)
	tandem1 := &TandemImpl {
		leader: couple2.Beau(),
		trailer: couple1.Beau(),
	}
	tandem2 := &TandemImpl {
		leader: couple2.Belle(),
		trailer: couple1.Belle(),
	}
	dancer.Reorder(tandem1.Leader(), tandem1.Trailer(), tandem2.Leader(), tandem2.Trailer())
	sample := TandemCouples(&TandemCouplesImpl {
		couple1: couple1,
		couple2: couple2,
		beaustandem: tandem1,
		bellestandem: tandem2,
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_TandemCouples_sample)
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
	node.Emit(TandemCouples(&TandemCouplesImpl{
		couple1: couple1,
		couple2: couple2,
		beaustandem: tandem1,
		bellestandem: tandem2,
	}))
}


type BackToBackCouples interface {
	Formation
	BackToBackCouples()
	Couple1() Couple           // defimpl:"read couple1" fe:"dancers"
	Couple2() Couple           // defimpl:"read couple2" fe:"dancers"
	BackToBack1() BackToBack  // defimpl:"read backtoback1"
	BackToBack2() BackToBack  // defimpl:"read backtoback2"
	// Roles:
	Beaus() dancer.Dancers     // fe:"no-slot"
	Belles() dancer.Dancers    // fe:"no-slot"
	Leaders()dancer.Dancers    // fe:"no-slot"
	Trailers() dancer.Dancers  // fe:"no-slot"
}

func (f *BackToBackCouplesImpl) String() string {
	return fmt.Sprintf("BackToBackCouples(%s, %s)", f.Couple1(), f.Couple2())
}

func (f *BackToBackCouplesImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.Couple1().Beaus(), f.Couple2().Beaus())
}

func (f *BackToBackCouplesImpl) Belles() dancer.Dancers {
	return dancer.Union(f.Couple1().Belles(), f.Couple2().Belles())
}

func (f *BackToBackCouplesImpl) Leaders () dancer.Dancers {
	return f.Dancers()
}

func (f *BackToBackCouplesImpl) Trailers() dancer.Dancers {
	return dancer.Dancers{}
}

func make_BackToBackCouples_sample() Formation {
	couple1 := make_Couple_sample().(*CoupleImpl)
	couple2 := make_Couple_sample().(*CoupleImpl)
	// Move couple1 down by one position:
	down := geometry.NewPositionDownLeft(geometry.Down1, geometry.Left0)
	couple2.Beau().MoveBy(down)
	couple2.Belle().MoveBy(down)
	// couple2 trade:
	swap_positions(
		couple1.Beau().Rotate(geometry.Direction2),
		couple1.Belle().Rotate(geometry.Direction2))
	dancer.Reorder(couple1.Beau(), couple1.Belle(), couple2.Beau(), couple2.Belle())
	sample := BackToBackCouples(&BackToBackCouplesImpl {
		couple1: couple1,
		couple2: couple2,
		backtoback1: &BackToBackImpl {
			dancer1: couple1.Beau(),
			dancer2: couple2.Belle(),
		},
		backtoback2: &BackToBackImpl {
			dancer1: couple2.Beau(),
			dancer2: couple1.Belle(),
		},
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_BackToBackCouples_sample)
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
	node.Emit(BackToBackCouples(&BackToBackCouplesImpl{
		couple1: couple1,
		couple2: couple2,
		backtoback1: bb1,
		backtoback2: bb2,
	}))
}


type BoxOfFour interface {
	Formation
	BoxOfFour()
	MiniWave1() MiniWave   // defimpl:"read miniwave1" fe:"dancers"
	MiniWave2() MiniWave   // defimpl:"read miniwave2" fe:"dancers"
	Tandem1() Tandem       // defimpl:"read tandem1"
	Tandem2() Tandem       // defimpl:"read tandem2"
	// Handedness:
	Handedness() Handedness
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	Leaders() dancer.Dancers
	Trailers() dancer.Dancers
}

func (f *BoxOfFourImpl) String() string {
	return fmt.Sprintf("BoxOfFour(%s, %s, %s, %s, %s)",
		f.Handedness(),
		f.Tandem1().Leader(), f.Tandem1().Trailer(),
		f.Tandem2().Leader(), f.Tandem2().Trailer())
}

func (f *BoxOfFourImpl) Handedness() Handedness {
	return f.MiniWave1().Handedness()
}

func (f *BoxOfFourImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Beaus(), f.MiniWave2().Beaus())
}

func (f *BoxOfFourImpl) Belles() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Belles(), f.MiniWave2().Belles())
}

func (f *BoxOfFourImpl) Leaders() dancer.Dancers {
	return dancer.Union(f.Tandem1().Leaders(), f.Tandem2().Leaders())
}

func (f *BoxOfFourImpl) Trailers() dancer.Dancers {
	return dancer.Union(f.Tandem1().Trailers(), f.Tandem2().Trailers())
}

func make_BoxOfFour_sample() Formation {
	tandem1 := make_Tandem_sample().(Tandem)
	tandem2 := make_Tandem_sample().(Tandem)
	right := geometry.NewPositionDownLeft(geometry.Down0, -geometry.Left1)
	tandem2.Leader().MoveBy(right)
	tandem2.Trailer().MoveBy(right)
	swap_positions(
		tandem2.Leader().Rotate(geometry.Direction2),
		tandem2.Trailer().Rotate(geometry.Direction2))
	dancer.Reorder(tandem1.Leader(), tandem1.Trailer(),
		tandem2.Leader(), tandem2.Trailer())
	sample := BoxOfFour(&BoxOfFourImpl {
		tandem1: tandem1,
		tandem2: tandem2,
		miniwave1: &MiniWaveImpl{
			dancer1: tandem1.Leader(),
			dancer2: tandem2.Trailer(),
		},
		miniwave2: &MiniWaveImpl{
			dancer1: tandem2.Leader(),
			dancer2: tandem1.Trailer(),
		},
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	// *** NOT PASSING TESTS YET
	RegisterFormationSample(make_BoxOfFour_sample)
}

func rule_BoxOfFour(node rete.Node, mw1, mw2 MiniWave, tandem1, tandem2 Tandem) {
	// *** We are getting a BoxOfFour for each MiniWave symetry
	// but we are also getting each of those twice.
	if mw1 == mw2 {
		return
	}
	// The direction test will also exclude duplicate tandems.
	if !tandem1.Direction().Opposite().Equal(tandem2.Direction()) {
		return
	}
	// Because each Tandem and each MiniWave come in as both of
	// the relevant inputs, the rule doesn't need to consider any
	// combinatorics, it can just test for a single arrangement of
	// the dancers:
	if !mw1.HasDancer(tandem1.Leader()) { return }
	if !mw1.HasDancer(tandem2.Trailer()) { return }
	if !mw2.HasDancer(tandem2.Leader()) { return }
	if !mw2.HasDancer(tandem1.Trailer()) { return }
	node.Emit(BoxOfFour(&BoxOfFourImpl{
		miniwave1: mw1,
		miniwave2: mw2,
		tandem1: tandem1,
		tandem2: tandem2,
	}))
}


type Star interface {
	Formation
	Star()
	MiniWave1() MiniWave       // defimpl:"read miniwave1" fe:"dancers"
	MiniWave2() MiniWave       // defimpl:"read miniwave2" fe:"dancers"
	// Handedness:
	Handedness() Handedness
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
}

func (f *StarImpl) String() string {
	return fmt.Sprintf("Star(%s, %s, %s, %s, %s)",
		f.Handedness(),
		f.MiniWave1().Dancer1(),
		f.MiniWave2().Dancer1(),
		f.MiniWave1().Dancer2(),
		f.MiniWave2().Dancer2())		
}

func (f *StarImpl) Handedness() Handedness {
	return f.MiniWave1().Handedness()
}

func (f *StarImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Beaus(), f.MiniWave2().Beaus())
}

func (f *StarImpl) Belles() dancer.Dancers {
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
	node.Emit(Star(&StarImpl{
		miniwave1: mw1,
		miniwave2: mw2,
	}))
}


type LineOfFour interface {
	Formation
	LineOfFour()
	LeftCouple() Couple      // defimpl:"read leftcouple" fe:"dancers"
	CenterCouple() Couple    // defimpl:"read centercouple"
	RightCouple() Couple     // defimpl:"read rightcouple" fe:"dancers"
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	Centers() dancer.Dancers
	Ends() dancer.Dancers
}

func (f *LineOfFourImpl) String() string {
	return fmt.Sprintf("LineOfFour(%s, %s, %s, %s)",
		f.LeftCouple().Beau(),
		f.LeftCouple().Belle(),
		f.RightCouple().Beau(),
		f.RightCouple().Belle())
}

func (f *LineOfFourImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.LeftCouple().Beaus(), f.RightCouple().Beaus())
}

func (f *LineOfFourImpl) Belles() dancer.Dancers {
	return dancer.Union(f.LeftCouple().Belles(), f.RightCouple().Belles())
}

func (f *LineOfFourImpl) Centers() dancer.Dancers {
	return f.CenterCouple().Dancers()
}

func (f *LineOfFourImpl) Ends() dancer.Dancers {
	return dancer.SetDifference(f.Dancers(), f.Centers())
}

func make_LineOfFour_sample() Formation {
	left := make_Couple_sample().(*CoupleImpl)
	right := make_Couple_sample().(*CoupleImpl)
	right2 := geometry.NewPositionDownLeft(geometry.Down0, -2 * geometry.Left1)
	right.Beau().MoveBy(right2)
	right.Belle().MoveBy(right2)
	center := &CoupleImpl {
		beau: left.Belle(),
		belle: right.Beau(),
	}
	dancer.Reorder(left.Beau(), left.Belle(), right.Beau(), right.Belle())
	sample :=  LineOfFour(&LineOfFourImpl {
		leftcouple: left,
		rightcouple: right,
		centercouple: center,
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_LineOfFour_sample)
}

func rule_LineOfFour(node rete.Node, c1, c2, c3 Couple) {
	if !(c1.Belle() == c2.Beau()) {
		return
	}
	if !(c2.Belle() == c3.Beau()) {
		return
	}
	node.Emit(LineOfFour(&LineOfFourImpl{
		leftcouple: c1,
		centercouple: c2,
		rightcouple: c3,
	}))
}


type WaveOfFour interface {
	Formation
	WaveOfFour()
	CenterMiniWave() MiniWave    // defimpl:"read centerminiwave"
	MiniWave1() MiniWave          // defimpl:"read miniwave1" fe:"dancers"
	MiniWave2() MiniWave          // defimpl:"read miniwave2" fe:"dancers"
	// Handedness:
	Handedness() Handedness
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	Centers() dancer.Dancers
	Ends() dancer.Dancers
}

func (f *WaveOfFourImpl) String() string {
	return fmt.Sprintf("WaveOfFour(%s, %s, %s, %s, %s)",
		f.Handedness(),
		f.MiniWave1().Dancer1(),
		f.MiniWave1().Dancer2(),
		f.MiniWave2().Dancer1(),
		f.MiniWave2().Dancer1())
}

func (f *WaveOfFourImpl) Handedness() Handedness {
	return f.MiniWave1().Handedness()
}

func (f *WaveOfFourImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Beaus(), f.MiniWave2().Beaus())
}

func (f *WaveOfFourImpl) Belles() dancer.Dancers {
	return dancer.Union(f.MiniWave1().Belles(), f.MiniWave2().Belles())
}

func (f *WaveOfFourImpl) Centers() dancer.Dancers {
	return f.CenterMiniWave().Dancers()
}

func (f *WaveOfFourImpl) Ends() dancer.Dancers {
	return dancer.SetDifference(f.Dancers(), f.Centers())
}

func rule_WaveOfFour(node rete.Node, mw1, mw2, mw3 MiniWave) {
	if mw1.Handedness() != mw2.Handedness().Opposite() {
		return
	}
	if mw3.Handedness() != mw2.Handedness().Opposite() {
		return
	}
	// *** We shoulld make sure all of the dancers are in line
	//
	// *** Are we assuming that we will see both symetric
	// MiniWaves of the same two dancers?
	if mw1.HasDancer(mw2.Dancer1()) {
		if !mw2.HasDancer(mw2.Dancer2()) {
			return
		}
	} else if mw1.HasDancer(mw2.Dancer2()) {
		if !mw2.HasDancer(mw2.Dancer1()) {
			return
		}
	}
	node.Emit(WaveOfFour(&WaveOfFourImpl{
		centerminiwave: mw2,
		miniwave1:mw1,
		miniwave2: mw3,
	}))
}


type TwoFacedLine interface {
	Formation
	TwoFacedLine()
	Couple1() Couple            // defimpl:"read couple1" fe:"dancers"
	Couple2() Couple            // defimpl:"read couple2" fe:"dancers"
	CenterMiniWave() MiniWave  // defimpl:"read centerminiwave"
	// Handedness:
	Handedness() Handedness
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	Centers() dancer.Dancers
	Ends() dancer.Dancers
}

func (f *TwoFacedLineImpl) String() string {
	return fmt.Sprintf("TwoFacedLine(%s, %s, %s, %s, %s)",
    	f.Handedness(),
		f.Couple1().Beau(),
		f.Couple1().Belle(),
		f.Couple2().Beau(),
		f.Couple2().Belle())
}

func (f *TwoFacedLineImpl) Handedness() Handedness {
	return f.CenterMiniWave().Handedness()
}

func (f *TwoFacedLineImpl) Beaus() dancer.Dancers {
	return dancer.Union(f.Couple1().Beaus(), f.Couple2().Beaus())
}

func (f *TwoFacedLineImpl) Belles() dancer.Dancers {
	return dancer.Union(f.Couple1().Belles(), f.Couple2().Belles())
}

func (f *TwoFacedLineImpl) Centers() dancer.Dancers {
	return f.CenterMiniWave().Dancers()
}

func (f *TwoFacedLineImpl) Ends() dancer.Dancers {
	return dancer.SetDifference(f.Dancers(), f.Centers())
}

func make_TwoFacedLine_sample() Formation {
	couple1 := make_Couple_sample().(*CoupleImpl)
	couple2 := make_Couple_sample().(*CoupleImpl)
	// Move couple2 to the right and turn them around:
	right2 := geometry.NewPositionDownLeft(geometry.Down0, -2 * geometry.Left1)
	couple2.Beau().MoveBy(right2)
	couple2.Belle().MoveBy(right2)
	swap_positions(
		couple2.Beau().Rotate(geometry.Direction2),
		couple2.Belle().Rotate(geometry.Direction2))
	dancer.Reorder(couple1.Beau(), couple1.Belle(), couple2.Beau(), couple2.Belle())
	sample := TwoFacedLine(&TwoFacedLineImpl {
		couple1: couple1,
		couple2: couple2,
		centerminiwave: &MiniWaveImpl {
			dancer1: couple1.Belle(),
			dancer2: couple2.Belle(),
		},
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	// *** NOT PASSING TESTS YET
	RegisterFormationSample(make_TwoFacedLine_sample)
}

func rule_TwoFacedLine(node rete.Node, c1, c2 Couple, mw MiniWave) {
	// Must have two different Couples:
	if c1.Beau() == c2.Beau() && c1.Belle() == c2.Belle() {
		return
	}
	if !((mw.HasDancer(c1.Beau()) && mw.HasDancer(c2.Beau())) ||
		(mw.HasDancer(c1.Belle()) && mw.HasDancer(c2.Belle()))) {
		return
	}
	node.Emit(TwoFacedLine(&TwoFacedLineImpl{
		couple1: c1,
		couple2: c2,
		centerminiwave: mw,
	}))
}

// Diamonds

// Zs
