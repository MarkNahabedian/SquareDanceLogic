// Package geometry implements the geometry of a square dance set or
// floor.  This package defines the coordinate system and provides the
// means for describing the position and direction of each dancer.
//
// The coordinate system provides a Down coordinate and a Left coordinate.
// Down and left are with respect to the caller's point of view.
// Down is a dancer's distance down the floor -- away from the caller.
// Left is the dancer's position from the right hand side of the set
// (from the caller's point of view) toward the caller's left.
//
// If one pictures the caller as being at the left hand edge of one's
// field of view, then Direction, Down and Left form the angle, X axis
// and Y axis of a normal right handed cartesean plane.
package geometry

import "fmt"
import "math"

// Direction is how angles are measured in our coordinate system.  It
// might descibe a direction of motion, the facing direction of a dancer,
// the direction of a dancer from another dancer's point of view, etc.
// Direction can be absolute or relative.
//
// Directions are expressed as fractions of a full circle, so a change
// in direction of 360 degrees is expressed as a change in Direction of 1.0.
// Direction increases in promenade direction.
//
// Direction 0.0 is the direction that the caller is facing and the
// facing direction of couple number one in a squared set.
// In a squared set, the facing direction of couple number two would be 0.25,
// that of couple number three: 0.5, and that of couple number four: 0.75.
type Direction float32

var Direction0 Direction = 0
var Direction1 Direction = Direction0 + FullCircle / 4
var Direction2 Direction = Direction0 + 2 * FullCircle / 4
var Direction3 Direction = Direction0 + 3 * FullCircle / 4

// FullCircle is a change in Direction of 1.0.
const FullCircle Direction = 1.0

// directionTolerance is the amount that the difference between two
// Directions must be within to be considered Equal.
const directionTolerance float64 = float64(FullCircle) / 20.0

func (d Direction) Canonicalize() Direction {
	// Should we also round to the nearest directionTolerance?
	return Direction(math.Remainder(float64(d), float64(FullCircle)))
}

// Opposite returns the direction opposite to d.
func (d Direction) Opposite() Direction {
	return d.Add(Direction(FullCircle / 2))
}

func (d Direction) QuarterRight() Direction {
	return d.Subtract(Direction(FullCircle / 4))
}

func (d Direction) QuarterLeft() Direction {
	return d.Add(Direction(FullCircle / 4))
}

// Inverse returns the arithmetic inverse dirrection to d.
func (d Direction) Inverse() Direction {
	return Direction(-float32(d)).Canonicalize()
}

// Add returns to sum of the two Directions.
func (d1 Direction) Add(d2 Direction) Direction {
	return Direction(float32(d1) + float32(d2)).Canonicalize()
}

// Subtract subtracts the Direction d2 from d1.
func (d1 Direction) Subtract(d2 Direction) Direction {
	return Direction(float32(d1) - float32(d2)).Canonicalize()
}

// Multiply multiplies the Direction d by multiplier.
func (d Direction) MultiplyBy(multiplier float32) Direction {
	return Direction(float32(d) * multiplier).Canonicalize()
}

// Divide divides the Direction d by divisor.
func (d Direction) DivideBy(divisor float32) Direction {
	return Direction(float32(d) / divisor).Canonicalize()
}

// Equal returns true if Directions d1 and d2 are within directionTolerance
// of each other.
func (d1 Direction) Equal(d2 Direction) bool {
	return math.Abs(float64(d1.Subtract(d2))) < directionTolerance
}

// Down is a distance along the Down axis -- away from the caller.
type Down float32

// Left is a distance along the Left axis -- towards the left hand
// wall from the caller's point of view.
type Left float32

const (
	Down0 = Down(0)
	Down1 = Down(CoupleDistance)
	Left0 = Left(0)
	Left1 = Left(CoupleDistance)
)

// Equal returns true if the two Down coordinates are within
// positionTolerance of each other.
func (d1 Down) Equal(d2 Down) bool {
	return math.Abs(float64(d1)-float64(d2)) < positionTolerance
}

// Add returns the sum of two Down coordinates.
func (d1 Down) Add(d2 Down) Down {
	return Down(float32(d1) + float32(d2))
}

// Subtract subtracts the Down coorsinate d2 from d1.
func (d1 Down) Subtract(d2 Down) Down {
	return Down(float32(d1) - float32(d2))
}

// Equal returns true if the two Left coordinates are within
// positionTolerance of each other.
func (d1 Left) Equal(d2 Left) bool {
	return math.Abs(float64(d1)-float64(d2)) < positionTolerance
}

// Add returns the sum of two Left coordinates.
func (d1 Left) Add(d2 Left) Left {
	return Left(float32(d1) + float32(d2))
}

// Subtract subtracts the Left coorsinate d2 from d1.
func (d1 Left) Subtract(d2 Left) Left {
	return Left(float32(d1) - float32(d2))
}

// CoupleDistance is the distance between (the center reference points
// of) two dancers standing side by side.
const CoupleDistance float32 = 1.0

const positionTolerance float64 = float64(CoupleDistance) / 1000.0

// Position represents a position on the floor or in a square dance set,
// or the relationship between one position and another.
type Position struct {
	Down Down
	Left Left
}

func (p Position) String() string {
	return fmt.Sprintf("{%f, %f}", p.Down, p.Left)
}

// NewPositionDownLeft returns a new Position with the given Down and Left
// values.
func NewPositionDownLeft(down Down, left Left) Position {
	return Position{
		Down: down,
		Left: left,
	}
}

// NewPosition returns a new Position computed from the given direction and
// distance.
func NewPosition(direction Direction, distance float32) Position {
	angle := float64(direction) * 2 * math.Pi
	return Position{
		Down: Down(float32(float64(distance) * math.Cos(angle))),
		Left: Left(float32(float64(distance) * math.Sin(angle))),
	}
}

// Equal returns true if the two positions have Equal Down coordinates
// and Equal Left coordinates.
func (p1 Position) Equal(p2 Position) bool {
	return p1.Down.Equal(p2.Down) && p1.Left.Equal(p2.Left)
}

// Distance returns the distance between two Positions.
func (p1 Position) Distance(p2 Position) float32 {
	d := float32(p2.Down.Subtract(p1.Down))
	l := float32(p2.Left.Subtract(p1.Left))
	return float32(math.Sqrt(float64(d*d + l*l)))
}

// Direction returns the Direction p2 is in from the perspective of p1.
func (p1 Position) Direction(p2 Position) Direction {
	d := float64(p2.Down.Subtract(p1.Down))
	l := float64(p2.Left.Subtract(p1.Left))
	return Direction(math.Atan2(l, d) / (2 * math.Pi))
}

// Add returns a new Position that's the result of adding p1 and p2.
func (p1 Position) Add(p2 Position) Position {
	return Position{
		Down: p1.Down.Add(p2.Down),
		Left: p1.Left.Add(p2.Left),
	}
}

// Subtract returns a new Position that's the result of subtracting p2 from p1.
func (p1 Position) Subtract(p2 Position) Position {
	return Position{
		Down: p1.Down.Subtract(p2.Down),
		Left: p1.Left.Subtract(p2.Left),
	}
}

func (p Position) Scale(scale float32) Position {
	return Position {
		Down: Down(float32(p.Down) * scale),
		Left: Left(float32(p.Left) * scale),
	}
}

// Magnitude returns the magnitude of the Position interpreted as a
// vector.
func (p Position) Magnitude() float32 {
	return float32(math.Sqrt(float64(p.Down * p.Down) + float64(p.Left * p.Left)))
}

func (p Position) Angle() Direction {
	d := float64(p.Down)
	l := float64(p.Left)
	return Direction(math.Atan2(l, d) / (2 * math.Pi))
}

// Center returns a new Position that's at the center of the specified
// Positions.
func Center(positions ...Position) Position {
	down := float32(0.0)
	left := float32(0.0)
	count := 0
	for _, p := range positions {
		down += float32(p.Down)
		left += float32(p.Left)
		count += 1
	}
	return Position{
		Down: Down(down / float32(count)),
		Left: Left(left / float32(count)),
	}
}
