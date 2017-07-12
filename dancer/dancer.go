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
	panic(fmt.Sprintf("Unsupported Gender #v", g))
}

// Opposite returns the Gender that is oooopposite to g.
func (g Gender) Opposite() Gender {
	switch g {
	case Unspecified:
		return Unspecified
	case Guy:
		return Gal
	case Gal:
		return Guy
	}
	panic(fmt.Sprintf("Unsupported Gender #v", g))
}

type dancer struct {
	// The set of dancers that this dancer is dancing with.
	Set          Set
	CoupleNumber int
	Gender       Gender
	// Each dancer in a set has a unique Ordinal.
	Ordinal         int
	Position        geometry.Position
	Direction       geometry.Direction
	OriginalPartner Dancer
}

type Dancer *dancer

type set struct {
	FlagpoleCenter geometry.Position
	Dancers        []Dancer
}

type Set *set

// SquaredSet returns a new squared set with the specified number of couples.
func NewSquaredSet(couples int) Set {
	circleFraction := geometry.FullCircle.DivideBy(float64(couples))
	s := set{
		FlagpoleCenter: geometry.NewPositionDownLeft(0, 0),
		Dancers:        make([]Dancer, couples*2),
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
			s.Dancers[index] = &dancer{
				Set:          &s,
				Ordinal:      index,
				Gender:       gender,
				CoupleNumber: couple + 1,
				Position: s.FlagpoleCenter.
					Add(geometry.NewPosition(
						facing.Opposite(), 1.5*geometry.CoupleDistance)).
					Add(geometry.NewPosition(
						facing.Add(geometry.FullCircle.DivideBy(4.0)),
						float64(side)*geometry.CoupleDistance/2.0)),
				Direction: geometry.Direction(facing),
			}
		}
	}
	for index, dancer := range s.Dancers {
		dancer.OriginalPartner = s.Dancers[index^1]
	}
	return &s
}
