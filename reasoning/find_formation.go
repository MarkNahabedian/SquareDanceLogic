package reasoning

import "reflect"
import "squaredance/dancer"


var formationFinder *FormationFinder = nil

func FindFormations(dancers dancer.Dancers, formation_type reflect.Type) []Formation {
	if formationFinder == nil {
		formationFinder = MakeFormationFinder()
	}
	formationFinder.Clear()
	formationFinder.Injest(dancers)
	result := []Formation{}
	formationFinder.DoFormations(formation_type, func (f Formation) {
		if f == nil {
			panic("nil Formation")
		}
		result = append(result, f)
	})
	return result
}

