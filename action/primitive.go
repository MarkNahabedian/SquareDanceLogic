// This file defines simple, primitive actions.
package action

import "fmt"
import "squaredance/geometry"
import "squaredance/dancer"
import "squaredance/reasoning"

func init() {
	// Actions which just change a Dancer's facing direction:

	defineAction("QuarterRight", "QuarterRight turns the dancers one wall to the right.")
	defineFormationAction("QuarterRight", Primitive, LookupFormationType("Dancer"),
		func(f reasoning.Formation) {
			d := f.Dancers()[0]
			d.Move(d.Position(), d.Direction().QuarterRight())
		})
	defineFormationAction("QuarterRight", Primitive, LookupFormationType("Dancers"),
		func(f reasoning.Formation) {
			for _, d := range f.Dancers() {
				d.Move(d.Position(), d.Direction().QuarterRight())
			}
		})

	defineAction("QuarterLeft", "QuarterLeft turns the dancers one wall to the right.")
	defineFormationAction("QuarterLeft", Primitive, LookupFormationType("Dancer"),
		func(f reasoning.Formation) {
			d := f.Dancers()[0]
			d.Move(d.Position(), d.Direction().QuarterLeft())
		})
	defineFormationAction("QuarterLeft", Primitive, LookupFormationType("Dancers"),
		func(f reasoning.Formation) {
			for _, d := range f.Dancers() {
				d.Move(d.Position(), d.Direction().QuarterLeft())
			}
		})

	defineAction("AboutFace", "AboutFace turns the dancers around 180 degrees.")
	defineFormationAction("AboutFace", Primitive, LookupFormationType("Dancer"),
		func(f reasoning.Formation) {
			d := f.Dancers()[0]
			d.Move(d.Position(), d.Direction().Opposite())
		})
	defineFormationAction("AboutFace", Primitive, LookupFormationType("Dancers"),
		func(f reasoning.Formation) {
			for _, d := range f.Dancers() {
				d.Move(d.Position(), d.Direction().Opposite())
			}
		})

	// Fragments of Dosado, Pass Thru and other calls where Dancers
	// approach and pass by each other:

	defineAction("TurnToFace", "Two dancers turn to face each other.")
	turnToFace := func(f reasoning.Formation) {
		dancers := f.Dancers()
		if len(dancers) != 2 {
			panic(fmt.Sprintf("The TurnTwoFace action requires two dancers, not %v", f))
		}
		update := func(this, other dancer.Dancer) {
			this.Move(this.Position(),
				this.Position().Direction(other.Position()))
		}
		update(dancers[0], dancers[1])
		update(dancers[1], dancers[0])
	}
	defineFormationAction("TurnToFace", Primitive, LookupFormationType("Couple"),
		turnToFace)
	defineFormationAction("TurnToFace", Primitive, LookupFormationType("MiniWave"),
		turnToFace)

	defineAction("Meet", "Meet moves FaceToFace Dancers up to meet each other.")
	defineFormationAction("Meet", Primitive, LookupFormationType("FaceToFace"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			distance := center.Distance(f.Dancers()[0].Position()) - 
				geometry.CoupleDistance / 2
			update := func(d dancer.Dancer) {
				dir := d.Position().Direction(center)
				_ = d.Move(d.Position().Add(geometry.NewPosition(dir, distance)), d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

	defineAction("ForwardLeft", "ForwardLeft moves FaceToFace dancers to a RightHanded MiniWave. This is commonly known as 'Touch'.")
	defineFormationAction("ForwardLeft", Primitive, LookupFormationType("FaceToFace"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			update := func(d dancer.Dancer) {
				d.Move(
					center.Add(geometry.NewPosition(d.Direction().QuarterLeft(),
						geometry.CoupleDistance / 2)),
					d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

	defineAction("ForwardRight", "ForwardRight moves FaceToFace dancers to a LeftHanded MiniWave.  This is commonly known as 'Left Touch'.")
	defineFormationAction("ForwardRight", Primitive, LookupFormationType("FaceToFace"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			update := func(d dancer.Dancer) {
				d.Move(
					center.Add(geometry.NewPosition(d.Direction().QuarterRight(),
						geometry.CoupleDistance / 2)),
					d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

	defineAction("PassToBacks", "PassToBacks moves dancers from a MiniWave to being BackToBack.")
	defineFormationAction("PassToBacks", Primitive, LookupFormationType("MiniWave"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			update := func(d dancer.Dancer) {
				d.Move(
					center.Add(geometry.NewPosition(d.Direction(),
						geometry.CoupleDistance / 2)),
					d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

	defineAction("BackwardLeft", "BackwardLeft moves BackToBack dancers to a RightHanded MiniWave.")
	defineFormationAction("BackwardLeft", Primitive, LookupFormationType("BackToBack"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			update := func(d dancer.Dancer) {
				d.Move(
					center.Add(geometry.NewPosition(d.Direction().QuarterLeft(),
						geometry.CoupleDistance / 2)),
					d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

	defineAction("BackwardRight", "BackwardRight moves BackToBack dancers to a LeftHanded MiniWave.")
	defineFormationAction("BackwardRight", Primitive, LookupFormationType("BackToBack"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			update := func(d dancer.Dancer) {
				d.Move(
					center.Add(geometry.NewPosition(d.Direction().QuarterRight(),
						geometry.CoupleDistance / 2)),
					d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

	defineAction("BackToFace", "BackToFace backs Dancers out of a MiniWave to face each other.")
	defineFormationAction("BackToFace", Primitive, LookupFormationType("MiniWave"),
		func(f reasoning.Formation) {
			dancers := f.Dancers()
			center := geometry.Center(dancer.Positions(dancers...)...)
			update := func(d dancer.Dancer) {
				d.Move(
					center.Add(geometry.NewPosition(d.Direction().Opposite(),
						geometry.CoupleDistance / 2)),
					d.Direction())
			}
			update(dancers[0])
			update(dancers[1])
		})

}
