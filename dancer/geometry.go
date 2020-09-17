package dancer

import "squaredance/geometry"


func (dancers Dancers) Positions() []geometry.Position {
	positions := []geometry.Position{}
	for _, dancer := range dancers {
		positions = append(positions, dancer.Position())
	}
	return positions
}


func (dancers Dancers) Center() geometry.Position {
	return geometry.Center(dancers.Positions()...)
}


func (dancers Dancers) Bounds() (leftmost, rightmost geometry.Left, downmost, upmost geometry.Down) {
	p := dancers.Center()
	leftmost = rightmost
	rightmost = p.Left
	downmost = upmost
	upmost = p.Down
	for _, dancer := range dancers {
		if left := dancer.Position().Left; left > leftmost {
			leftmost = left
		} else if left < rightmost {
			rightmost = left
		}
		if down := dancer.Position().Down; down  > downmost {
			downmost = down
		} else if down < upmost {
			upmost = down
		}
	}
	return leftmost, rightmost, upmost, downmost
}


func Distance (dancer1, dancer2 Dancer) float32 {
	return dancer1.Position().Subtract(dancer2.Position()).Magnitude()
}

