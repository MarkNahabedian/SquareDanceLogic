package action

import "math"
import "squaredance/dancer"
import "squaredance/geometry"


// Breathe finds dancers that are too close together and spreads them
// apart.
func Breathe(dancers dancer.Dancers) {
	any := true
	for any {
		any = false
		for i, d1 := range dancers {
			for j := i+1; j < len(dancers); j++ {
				d2 := dancers[j]
				distance := dancer.Distance(d1, d2)
				spread_by := geometry.CoupleDistance - distance
				if  spread_by <= 0  {
					continue
				}
				any = true
				midpoint := dancer.Dancers{d1, d2}.Center()
				spread(dancers, midpoint,
					geometry.NewPosition(midpoint.Direction(d1.Position()),
						spread_by))
				break
			}
			if any {
				break
			}
		}
	}
}

// spread moves dancers away from midpoiint by + or - vector, as appropriate.
func spread(dancers dancer.Dancers, midpoint geometry.Position, vector geometry.Position) {
	vdir := vector.Angle()
	vodir := vdir.Opposite()
	ovector := geometry.NewPositionDownLeft(0.0, 0.0).Subtract(vector)
	for _, dancer := range dancers {
		ddir := midpoint.Direction(dancer.Position())
		if math.Abs(float64(ddir - vdir)) < math.Abs(float64(ddir - vodir)) {
			dancer.MoveBy(vector)
		} else if math.Abs(float64(ddir - vdir)) > math.Abs(float64(ddir - vodir)){
			dancer.MoveBy(ovector)
		}
	}
}

