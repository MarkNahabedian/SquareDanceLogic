// Definitions and rules about all two dancer square dance formations.
package reasoning

import "fmt"
import "goshua/rete"
import "squaredance/dancer"
import "squaredance/geometry"


// Pair represents two distinct Dancers.
// Note that rule_PairOfDancers does not filter by Dancer canonical ordering
// and the rete node that Joins one Dancer with another considers the those
// Dancers in both possible orderings, so for every two Dancers, two Pairs
// are made, one with one dancer as Dancer1, and the other with the other
// Danceer as Dancer1.  This should simplify a number of the other two Dancer
// rules, which dont need to consider which Dancer is which in a given Pair
// because there will be another Pair with its Dancers in the other ordering.
type Pair interface {
	// Should Pair be a Formation?
	Pair()
	Dancer1() dancer.Dancer
	Dancer2() dancer.Dancer
}

func MakePair(dancer1, dancer2 dancer.Dancer) Pair {
	if dancer1.Ordinal() < dancer2.Ordinal() {
		return Pair(&pair{ dancer1: dancer1, dancer2: dancer2 })
	}
	return Pair(&pair{ dancer1: dancer2, dancer2: dancer1 })
}

type pair struct {
	dancer1 dancer.Dancer
	dancer2 dancer.Dancer
}

func (p *pair) Pair() {}

// Dancer1 returns one dancer of a Pair.
func (p *pair) Dancer1() dancer.Dancer {
	return p.dancer1
}

// Dancer2 returns the Dancer of the Pair that is not returned by Dancer1.
func (p *pair) Dancer2() dancer.Dancer {
	return p.dancer2
}

func (p *pair) Ordinal() int {
	return p.dancer1.Ordinal()
}

func (p *pair) String() string {
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
}

func (f *CoupleImpl) Ordinal() int {
	return f.Beau().Ordinal()
}

func (f *CoupleImpl) String() string {
	return fmt.Sprintf("Couple(%s, %s)", f.Beau(), f.Belle())
}

func (f *CoupleImpl) Beaus() dancer.Dancers { return dancer.Dancers{ f.Beau() } }
func (f *CoupleImpl) Belles() dancer.Dancers { return dancer.Dancers{ f.Belle() } }

func rule_GeneralizedCouple(node rete.Node, p Pair) {
	d1 := p.Dancer1()
	d2 := p.Dancer2()
	if RightOf(d1, d2) && LeftOf(d2, d1) {
		node.Emit(Couple(&CoupleImpl{beau: d1, belle: d2}))
	}
}


// A MiniWave consists of two Dancers facing in opposite directions.
// There isn't anything inate to a MiniWave that would restrict which
// of the same two Dancers is Dancer1 versus which is Dancer2 unless
// we resort to Dancer.Ordinal.
type MiniWave interface{
	Formation
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


// FaceToFace represents two dancewrs that are facing each other.
type FaceToFace interface {
	Formation
	FaceToFace()
	Dancer1() dancer.Dancer     // defimpl:"read dancer1" fe:"dancers" 
	Dancer2() dancer.Dancer     // defimpl:"read dancer2" fe:"dancers"
	// Roles
	Leaders() dancer.Dancers
	Trailers() dancer.Dancers
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
	if d1.Ordinal() >= d2.Ordinal() {  // Huh?  de-dup?
		return
	}
	if InFrontOf(d1, d2) && InFrontOf(d2, d1) {
		node.Emit(FaceToFace(&FaceToFaceImpl{dancer1: d1, dancer2: d2}))
	}
}


// BackToBack represents two dancers with their backs to each other.
type BackToBack interface {
	Formation
	BackToBack()
	Dancer1() dancer.Dancer    // defimpl:"read dancer1" fe:"dancers"
	Dancer2() dancer.Dancer    // defimpl:"read dancer2" fe:"dancers"
	// Roles:
	Leaders() dancer.Dancers
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

