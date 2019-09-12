// Package dancer implements a model of square dancers and the space they
// dance in.
package dancer

import "bytes"
import "fmt"
import "sort"
import "squaredance/geometry"


// Gender represents the gender of a square dancer.
type Gender int

const (
	// Unspecified is used when the Gender is unknown or doesn't matter.
	Unspecified Gender = iota
	Guy
	Gal
)

func (g Gender) String() string {
	switch g {
	case Unspecified:
		return "Unspecified"
	case Guy:
		return "Guy"
	case Gal:
		return "Gal"
	}
	panic(fmt.Sprintf("Unsupported Gender %d", int(g)))
}

// Equal returns true if the Genders are the same.
// Two Unspecified Genders are never the same.
func (g Gender) Equal(g2 Gender) bool {
	if g == Unspecified { return false }
	if g2 == Unspecified { return false }
	return g == g2
}

// Opposite returns the Gender that is opposite to g.
func (g Gender) Opposite() Gender {
	switch g {
	case Unspecified:
		return Unspecified
	case Guy:
		return Gal
	case Gal:
		return Guy
	}
	panic(fmt.Sprintf("Unsupported Gender %v", g))
}

type Dancer interface{
	String() string
	IsDancer() bool
	Set() Set                           // defimpl:"read set"
	// Couple numbers atart at 1 for the #1 head couple (facing down
	// the hall). CoupleNumber is only meaningful for Dancers that
	// started in a squared set.  CoupleNumber <= 0 are invalid.
	CoupleNumber() int                // defimpl:"read coupleNumber"
	Gender() Gender                    // defimpl:"read gender"
	// Each Dancer in a set has a unique Ordinal.  Ordinal is used to
	// avoid duplicating symetric Formations.
	Ordinal() int                      // defimpl:"read ordinal"
	Position() geometry.Position       // defimpl:"read position"
	Direction() geometry.Direction     // defimpl:"read direction"
	OriginalPartner() Dancer          // defimpl:"read originalPartner"
	SetOriginalPartner(Dancer)
	Move(geometry.Position, geometry.Direction) Dancer
	// A single dancer is still a formation so it implements the Formation interface
	NumberOfDancers() int
	Dancers() Dancers
	HasDancer(Dancer) bool
}

type Dancers []Dancer

func (ds Dancers) String() string {
	buf := bytes.NewBufferString("")
	first := true
	for _, d := range ds {
		buf.WriteString(d.String())
		if !first {
			buf.WriteString(", ")
		} else {
			first = false
		}
	}
	return buf.String()
}

// Enable sorting by ordinal"

func (ds Dancers) Len() int { return len(ds) }

func (ds Dancers) Swap(i, j int) {
	ds[i], ds[j] = ds[j], ds[i]
}

func (ds Dancers) Less(i, j int) bool {
	return ds[i].Ordinal() < ds[j].Ordinal()
}

// Ordered sorts the dancers by their Ordinals.  Dancers is modified
// to reflect that ordering.
func (ds Dancers) Ordered() Dancers {
	sort.Sort(ds)
	return ds
}


func (d *DancerImpl) String() string {
	if d.CoupleNumber() <= 0 || d.Gender() == Unspecified {
		return fmt.Sprintf("Dancer_%d", d.Ordinal())
	}
	return fmt.Sprintf("Dancer_%d%s", d.CoupleNumber(), d.Gender())
}

func (d *DancerImpl) IsDancer() bool { return true }

func (d *DancerImpl) SetOriginalPartner(d2 Dancer) {
	d.originalPartner = d2
}

// Modify the Dancer's position and direction.
func (d *DancerImpl) Move(newPosition geometry.Position, newDirection geometry.Direction) Dancer {
	d.position = newPosition
	d.direction = newDirection
	return d
}

func (d1 *DancerImpl) GoshuaEqual(d2 interface{}) (bool, error) {
	// If DancerImpls aren't EQ then they're not EQUAL.
	return false, nil
}


type Set interface {
	FlagpoleCenter() geometry.Position   // defimpl:"read flagpoleCenter"
	Dancers()        Dancers              // defimpl:"read dancers"
	// Support the reasoning/Formation interface:
	NumberOfDancers() int
	HasDancer(d Dancer) bool
}

// NumberOfDancers is part of the reasoning/Formatioon interface.
func (s *SetImpl) NumberOfDancers() int{
	return s.Dancers().NumberOfDancers()
}

// HasDancer is part of the reasoning/Formatioon interface.
func (s *SetImpl) HasDancer(d Dancer) bool {
	return s.Dancers().HasDancer(d)
}


// SquaredSet returns a new squared set with the specified number of couples.
func NewSquaredSet(couples int) Set {
	circleFraction := geometry.FullCircle.DivideBy(float64(couples))
	s := SetImpl{
		flagpoleCenter: geometry.NewPositionDownLeft(0, 0),
		dancers:        make(Dancers, couples*2),
	}
	for couple := 0; couple < couples; couple++ {
		facing := circleFraction.MultiplyBy(float64(couple))
		for _, gender := range []Gender{Guy, Gal} {
			var side int
			var ordinalAdjustment int
			switch gender {
			case Guy:
				side = 1
				ordinalAdjustment = 0
			case Gal:
				side = -1
				ordinalAdjustment = 1
			}
			index := 2*couple + ordinalAdjustment
			s.dancers[index] = &DancerImpl{
				set:          &s,
				ordinal:      index,
				gender:       gender,
				coupleNumber: couple + 1,
				position: s.flagpoleCenter.
					Add(geometry.NewPosition(
						facing.Opposite(), 1.5*geometry.CoupleDistance)).
					Add(geometry.NewPosition(
						facing.Add(geometry.FullCircle.DivideBy(4.0)),
						float64(side)*geometry.CoupleDistance/2.0)),
				direction: geometry.Direction(facing),
			}
		}
	}
	for index, dancer := range s.dancers {
		dancer.SetOriginalPartner(s.dancers[index^1])
	}
	return &s
}

// MakeSomeDancers returns the specified number of Gender neutral Dancers.
func MakeSomeDancers(count int) Dancers {
	dancers := Dancers{}
	for ordinal := 0; ordinal < count; ordinal++ {
		dancers = append(dancers,&DancerImpl{
			set:          nil,
			ordinal:      ordinal,
			gender:       Unspecified,
			coupleNumber: -1,
		})
	}
	return dancers
}


// Positions returns the Position of each Dancer
func Positions(dancers ...Dancer) []geometry.Position {
	length := len(dancers)
	positions := make([]geometry.Position, length, length)
	for i := 0; i < length; i++ {
		positions[i] = dancers[i].Position()
	}
	return positions
}


// Union returns the dancers that are present in any of the slices.
func Union(dancerSets ...Dancers) Dancers {
	got := map[Dancer]bool{}
	for _, set := range dancerSets {
		for _, d := range set {
	    	got[d] = true
		}
	}
	result := Dancers{}
	for d, keep := range got {
		if keep {
			result = append(result, d)
		}
	}
	return result
}


// Intersection returns the Dancers that are present in all of the slices.
func Intersection(dancerSets ...Dancers) Dancers {
	attendance := make(map[Dancer] []bool)
	get := func (d Dancer) []bool {
		b, found := attendance[d]
		if found {
			return b
		}
		b = make([]bool, len(dancerSets), len(dancerSets))
		attendance[d] = b
		return b
	}
	for i, s := range dancerSets {
		for _, d := range s {
			get(d)[i] = true
		}
	}
	result := Dancers{}
	for d, b := range attendance {
		for i := 0; i < len(b); i++ {
			if !b[i] {
				goto skip
			}
		}
		result = append(result, d)
		skip:
	}
	return result
}


// SetDifference returns a slice of those Dancers that are in universe
// but not in minus.
func SetDifference(universe Dancers, minus Dancers) Dancers {
	skip := func(d Dancer) bool {
		for _, d2 := range minus {
			if d == d2 {
				return true
			}
		}
		return false
	}
	result := Dancers{}
	for _, d := range universe {
		if !skip(d) {
			result = append(result, d)
		}
	}
	return result
}


// Since a single dancer can be considered a Formation, Dancer should
// implement the reasoning.Formation interface:

// NumberOfDancers is part of the reasoning.Formation interface.
func (d *DancerImpl) NumberOfDancers() int { return 1 }

// Dancers is part of the reasoning.Formation interface.
func (d *DancerImpl) Dancers() Dancers {
	return Dancers { d }
}

// HasDancer is part of the reasoning.Formation interface.
func (d *DancerImpl) HasDancer(d2 Dancer) bool {
	if d2, ok := d2.(*DancerImpl); ok {
		return d == d2
	}
	return false
}


// It's convenient to have Dancers implement the Formation interface too.

// NumberOfDancers is part of the reasoning.Formation interface.
func (f Dancers) NumberOfDancers() int { return len(f) }

// Dancers is part of the reasoning.Formation interface.
func (f Dancers) Dancers() Dancers {
	return f
}

// HasDancer is part of the reasoning.Formation interface.
func (f Dancers) HasDancer(d2 Dancer) bool {
	for _, d := range f {
		if d == d2 {
			return true
		}
	}
	return false
}
