// Package dancer implements a model of square dancers and the space they
// dance in.
package dancer

import "fmt"
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
	IsDancer() bool
	Set() Set
	CoupleNumber() int
	Gender() Gender
	Ordinal() int
	Position() geometry.Position
	Direction() geometry.Direction
	OriginalPartner() Dancer
	SetOriginalPartner(Dancer)
}

type dancer struct {
	// The set of dancers that this dancer is dancing with.
	set          Set
	coupleNumber int
	gender       Gender
	// Each dancer in a set has a unique Ordinal.
	ordinal         int
	position        geometry.Position
	direction       geometry.Direction
	originalPartner Dancer
}


func (d *dancer) String() string {
	return fmt.Sprintf("Dancer_%d%s", d.CoupleNumber(), d.Gender())
}

func (d *dancer) IsDancer() bool { return true }

func (d *dancer) Set() Set { return d.set }

func (d *dancer) CoupleNumber() int { return d.coupleNumber }

func (d *dancer) Gender() Gender { return d.gender }

func (d *dancer) Ordinal() int { return d.ordinal }

func (d *dancer) Position() geometry.Position { return d.position }

func (d *dancer) Direction() geometry.Direction { return d.direction }

func (d *dancer) OriginalPartner() Dancer { return d.originalPartner }

func (d *dancer) SetOriginalPartner(d2 Dancer) {
	d.originalPartner = d2
}

func (d1 *dancer) GoshuaEqual(d2 interface{}) (bool, error) {
	// If Dancers aren't EQ then they're not EQUAL.
	return false, nil
}


type Set interface {
	FlagpoleCenter() geometry.Position
	Dancers()        []Dancer

}

type set struct {
	flagpoleCenter geometry.Position
	dancers        []Dancer
}

func (s *set) FlagpoleCenter() geometry.Position {
	return s.flagpoleCenter
}

func (s *set) Dancers() []Dancer {
	return s.dancers
}

// SquaredSet returns a new squared set with the specified number of couples.
func NewSquaredSet(couples int) Set {
	circleFraction := geometry.FullCircle.DivideBy(float64(couples))
	s := set{
		flagpoleCenter: geometry.NewPositionDownLeft(0, 0),
		dancers:        make([]Dancer, couples*2),
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
			s.dancers[index] = &dancer{
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
