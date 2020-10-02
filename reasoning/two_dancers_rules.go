// Definitions and rules about all two dancer square dance formations.
package reasoning

import "fmt"
import "reflect"
import "goshua/rete"
import "squaredance/dancer"
import "squaredance/geometry"


// Near returns true if the two dancers are near each other.
func Near(dancer1, dancer2 dancer.Dancer) bool {
	// *** Should we add a bit of fudge?
	return dancer1.Position().Distance(dancer2.Position()) <= 1.2 * geometry.CoupleDistance
}


// Pair represents two distinct Dancers.
// Note that rule_PairOfDancers does not filter by Dancer canonical ordering
// and the rete node that Joins one Dancer with another considers the those
// Dancers in both possible orderings, so for every two Dancers, two Pairs
// are made, one with one dancer as Dancer1, and the other with the other
// Danceer as Dancer1.  This should simplify a number of the other two Dancer
// rules, which don't need to consider which Dancer is which in a given Pair
// because there will be another Pair with its Dancers in the other ordering.
type Pair interface {
	// Should Pair be a Formation?
	Pair()
	Dancer1() dancer.Dancer   // defimpl:"read dancer1"
	Dancer2() dancer.Dancer   // defimpl:"read dancer2"
}

func MakePair(dancer1, dancer2 dancer.Dancer) Pair {
/*
	if dancer1.Ordinal() < dancer2.Ordinal() {
		return Pair(&pair{ dancer1: dancer1, dancer2: dancer2 })
	}
*/
	return Pair(&PairImpl{ dancer1: dancer2, dancer2: dancer1 })
}

func (p *PairImpl) Pair() {}

/*
func (p *PairImpl) Ordinal() int {
	return p.dancer1.Ordinal()
}
*/

func (p *PairImpl) String() string {
	return fmt.Sprintf("Pair(%s, %s)", p.dancer1, p.dancer2)
}

// rule_PairOfDancers groups each two Dancers pairwise.
func rule_PairOfDancers(node rete.Node, dancer1, dancer2 dancer.Dancer) {
	// Note that for each pair of dancers we will create two Pair
	// objects: one each for the two possible orders of the dancers.
	if dancer1 != dancer2 {
		node.Emit(Pair(MakePair(dancer1, dancer2)))
	}
}


// TODO: Need to test for nearness or adjacency.


// A Couple consists of two Dancers that are side by side and facing the
// same direction.  Since one Dancer is to the right of the other Dancer
// the same two Dancers can only be in a Couple one way.
type Couple interface {
	Formation
	Couple()
	Beau() dancer.Dancer      // defimpl:"read beau" fe:"dancers"
	Belle() dancer.Dancer     // defimpl:"read belle" fe:"dancers"
	// Roles:
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
	IsNormal() bool
	IsHalfSasheyed()bool
}

func (f *CoupleImpl) Ordinal() int {
	return f.Beau().Ordinal()
}

func (f *CoupleImpl) String() string {
	return fmt.Sprintf("Couple(%s, %s)", f.Beau(), f.Belle())
}

func (f *CoupleImpl) Beaus() dancer.Dancers { return dancer.Dancers{ f.Beau() } }
func (f *CoupleImpl) Belles() dancer.Dancers { return dancer.Dancers{ f.Belle() } }

func (f *CoupleImpl) IsNormal() bool {
	return f.Beau().Gender() == dancer.Guy &&
		f.Belle().Gender() == dancer.Gal
}

func (f *CoupleImpl) IsHalfSasheyed()bool {
	return f.Beau().Gender() == dancer.Gal &&
		f.Belle().Gender() == dancer.Guy
}

func rule_GeneralizedCouple(node rete.Node, p Pair) {
	d1 := p.Dancer1()
	d2 := p.Dancer2()
	if RightOf(d1, d2) && LeftOf(d2, d1) {
		node.Emit(Couple(&CoupleImpl{beau: d1, belle: d2}))
	}
}

func make_Couple_sample() Formation {
	dancers := dancer.MakeSomeDancers(2)
	beau := dancers[0]
	belle := dancers[1]
	beau.Move(geometry.Position{ Left: geometry.Left1, Down: geometry.Down0 },
		geometry.Direction0)
	belle.Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
		geometry.Direction0)
	dancer.Reorder(beau, belle)
	sample := Couple(&CoupleImpl {
		beau: beau,
		belle: belle,
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_Couple_sample)
}


type TwoDancerSymetric interface {
	isSymetricHelper(TwoDancerSymetric) bool
	Dancer1() dancer.Dancer
	Dancer2() dancer.Dancer
}

// IsTwoDancerSymetric returns true if the two Formations would be
// identical to each other if the dancers were swapped.
func IsTwoDancerSymetric(formation1, formation2 interface{}) bool {
	f1, ok1 := formation1.(TwoDancerSymetric)
	f2, ok2 := formation2.(TwoDancerSymetric)
	if !(ok1 && ok2) {
		return false
	}
	if reflect.TypeOf(formation1) != reflect.TypeOf(formation2) {
		return false
	}
	return (f1.isSymetricHelper(f2) &&
		f1.Dancer1() == f2.Dancer2() &&
                f1.Dancer2() == f2.Dancer1())
}


// A MiniWave consists of two Dancers facing in opposite directions.
// There isn't anything inate to a MiniWave that would restrict which
// of the same two Dancers is Dancer1 versus which is Dancer2 unless
// we resort to Dancer.Ordinal.
type MiniWave interface{
	Formation
	TwoDancerSymetric
	MiniWave()
	Dancer1() dancer.Dancer   // defimpl:"read dancer1" fe:"dancers"
	Dancer2() dancer.Dancer   // defimpl:"read dancer2" fe:"dancers"
    // Handedness
	Handedness() Handedness
	// Roles
	Beaus() dancer.Dancers
	Belles() dancer.Dancers
}

func MakeMiniWave(dancer1, dancer2 dancer.Dancer) MiniWave {
	if dancer1.Ordinal() < dancer2.Ordinal() {
		return &MiniWaveImpl{ dancer1, dancer2 }
	}
	return &MiniWaveImpl{ dancer2, dancer1 }
}

func (mw *MiniWaveImpl) isSymetricHelper(other TwoDancerSymetric) bool {
	o, ok := other.(*MiniWaveImpl)
	return ok && mw.Handedness() == o.Handedness()
}

func (mw *MiniWaveImpl) String() string {
	return fmt.Sprintf("MiniWave(%s, %s, %s)", mw.Dancer1(), mw.Dancer2(), mw.Handedness())
}

func (mw *MiniWaveImpl) Handedness() Handedness {
	if RightOf(mw.Dancer1(), mw.Dancer2()) {
		return RightHanded
	}
	return LeftHanded
}

func (mw *MiniWaveImpl) Beaus() dancer.Dancers {
	switch mw.Handedness() {
		case RightHanded:
			return dancer.Dancers{
				mw.Dancer1(), 
				mw.Dancer2(), 
			}
		case LeftHanded:
			return dancer.Dancers{}
		default:
			panic("MiniWave is neither right nor left handed.")
	}
}

func (mw *MiniWaveImpl) Belles() dancer.Dancers {
	switch mw.Handedness() {
		case RightHanded:
			return dancer.Dancers{}
		case LeftHanded:
			return dancer.Dancers{
				mw.Dancer1(), 
				mw.Dancer2(), 
			}
		default:
			panic("MiniWave is neither right nor left handed.")
	}
}

func rule_MiniWave(node rete.Node, p Pair) {
	d1 := p.Dancer1()
	d2 := p.Dancer2()
	if d1 == d2 {
		return
	}
	if !Near(d1, d2) {
		return
	}
	if d1.Ordinal() >= d2.Ordinal() {   // Huh?  de-dup?
		return
	}
	if RightOf(d1, d2) && RightOf(d2, d1) {
		node.Emit(MakeMiniWave(d1, d2))
		return
	}
	if LeftOf(d1, d2) && LeftOf(d2, d1) {
		node.Emit(MakeMiniWave(d1, d2))
	}
}

func make_MiniWave_sample() Formation {
	dancers := dancer.MakeSomeDancers(2)
	// Right handed
	dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
		geometry.FullCircle / 2)
	dancers[1].Move(geometry.Position{ Left: geometry.Left1, Down: geometry.Down0 },
		geometry.Direction0)
	dancer.Reorder(dancers...)
	sample := MiniWave(&MiniWaveImpl {
		dancer1: dancers[0],
		dancer2: dancers[1],
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_MiniWave_sample)
}


// FaceToFace represents two dancewrs that are facing each other.
type FaceToFace interface {
	Formation
	TwoDancerSymetric
	FaceToFace()
	Dancer1() dancer.Dancer     // defimpl:"read dancer1" fe:"dancers" 
	Dancer2() dancer.Dancer     // defimpl:"read dancer2" fe:"dancers"
	// Roles
	Leaders() dancer.Dancers
	Trailers() dancer.Dancers
}

func (mw *FaceToFaceImpl)  isSymetricHelper(other TwoDancerSymetric) bool {
	return true
}

func (f *FaceToFaceImpl) String() string {
	return fmt.Sprintf("FaceToFace(%s, %s)", f.Dancer1(), f.Dancer2())
}

func (f *FaceToFaceImpl) Leaders() dancer.Dancers {
	return dancer.Dancers {}
}

func (f *FaceToFaceImpl) Trailers() dancer.Dancers {
	return dancer.Dancers {
		f.Dancer1(),
		f.Dancer2(),
	}
}

func rule_FaceToFace(node rete.Node, p Pair) {
	d1 := p.Dancer1()
	d2 := p.Dancer2()
	// Remove the duplication that's inherent in Pair:
	if d1.Ordinal() >= d2.Ordinal() {
		return
	}
	if InFrontOf(d1, d2) && InFrontOf(d2, d1) {
		node.Emit(FaceToFace(&FaceToFaceImpl{dancer1: d1, dancer2: d2}))
	}
}

func make_FaceToFace_sample() Formation  {
	dancers := dancer.MakeSomeDancers(2)
	dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
		geometry.Direction0)
	dancers[1].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down1 },
		geometry.FullCircle / 2)
	dancer.Reorder(dancers...)
	sample := FaceToFace(&FaceToFaceImpl{
		dancer1: dancers[0],
		dancer2: dancers[1],
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_FaceToFace_sample)
}


// BackToBack represents two dancers with their backs to each other.
type BackToBack interface {
	Formation
	TwoDancerSymetric
	BackToBack()
	Dancer1() dancer.Dancer    // defimpl:"read dancer1" fe:"dancers"
	Dancer2() dancer.Dancer    // defimpl:"read dancer2" fe:"dancers"
	// Roles:
	Leaders() dancer.Dancers
}

func (mw *BackToBackImpl)  isSymetricHelper(other TwoDancerSymetric) bool {
	return true
}

func (f *BackToBackImpl) String() string {
	return fmt.Sprintf("BackToBack(%s, %s)", f.Dancer1(), f.Dancer2())
}

func (f *BackToBackImpl) Leaders() dancer.Dancers {
	return dancer.Dancers {
		f.Dancer1(),
		f.Dancer2(),
	}
}

func rule_BackToBack(node rete.Node, p Pair) {
	d1 := p.Dancer1()
	d2 := p.Dancer2()
	if d1.Ordinal() >= d2.Ordinal() {
		return
	}
	if Behind(d1, d2) && Behind(d2, d1) {
		node.Emit(BackToBack(&BackToBackImpl{dancer1: d1, dancer2: d2}))
	}
}

func make_BackToBack_sample() Formation {
	dancers := dancer.MakeSomeDancers(2)
	dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
		geometry.FullCircle / 2)
	dancers[1].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down1 },
		geometry.Direction0)
	dancer.Reorder(dancers...)
	sample := BackToBack(&BackToBackImpl{
		dancer1: dancers[0],
		dancer2: dancers[1],
	})
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_BackToBack_sample)
}


// Tandem represents two dancers facing the same direction with the Leader
// in front of the Trailer.
type Tandem interface {
	Formation
	Tandem()
	Leader() dancer.Dancer          // defimpl:"read leader" fe:"dancers"
	Trailer() dancer.Dancer         // defimpl:"read trailer" fe:"dancers"
	Direction() geometry.Direction
	// Roles:
	Leaders() dancer.Dancers
	Trailers() dancer.Dancers
}

func (f *TandemImpl) String() string {
	return fmt.Sprintf("Tandem(%s, %s)", f.Leader(), f.Trailer())
}

func (f *TandemImpl) Leaders() dancer.Dancers {
	return dancer.Dancers{ f.Leader() }
}

func (f *TandemImpl) Trailers() dancer.Dancers {
	return dancer.Dancers{ f.Trailer() }
}

func (t *TandemImpl) Direction() geometry.Direction {
	return t.Leader().Direction()
}

func rule_Tandem(node rete.Node, p Pair) {
	d1 := p.Dancer1()
	d2 := p.Dancer2()
	if !d1.Direction().Equal(d2.Direction()) {
		return
	}
	if Behind(d1, d2) && InFrontOf(d2, d1) {
		node.Emit(Tandem(&TandemImpl{leader: d1, trailer: d2}))
	}
}

func make_Tandem_sample() Formation {
	dancers := dancer.MakeSomeDancers(2)
	leader, trailer := dancers[0], dancers[1]
	leader.Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down1 },
		geometry.Direction0)
	trailer.Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
		geometry.Direction0)
	dancer.Reorder(leader, trailer)
	sample := &TandemImpl {
		leader: leader,
		trailer: trailer,
	}
	sample.Dancers().Recenter0()
	return sample
}

func init() {
	RegisterFormationSample(make_Tandem_sample)
}

