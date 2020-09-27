// Package action is responsible for moving Dancers around.
package action

import "fmt"
import "reflect"
import "squaredance/reasoning"

type Level int

const (
	Primitive Level = iota
	Basic1
	Basic2
	Mainstream
	Plus
	A_1
	A_2
	C_1
	C_2
	C_3A
	C_3B
	C_4
	NotOnList
)


// FormationType is a reflect.Type identifying the interface type of a reasoning.Formation.
type FormationType reflect.Type

func LookupFormationType(name string) FormationType {
	ft, ok := reasoning.AllFormationTypes[name]
	if !ok {
		panic(fmt.Sprintf("No formation named %q", name))
	}
	return FormationType(ft)
}


// Action represents some snipet of square dance choreography -- the
// simplest action a Dancer or Formation of Dancers can perform.  We
// use them as primitive choreography when defining how a square dance
// call should be interpreted.
// How an action is performed might depend on the starting formation.
// Each Action maps a Formation to a FormationAction which represents
// the details of how the action should be performed from that Formation.
type Action interface {
	Name() string              // defimpl:"read name"
	Description() string      // defimpl:"read description"
	AddFormationAction(...FormationAction)     // defimpl:"append formationActions"
	DoFormationActions(func(FormationAction) bool)         // defimpl:"iterate formationActions"
	GetFormationAction(FormationType) FormationAction
	GetFormationActionFor(f reasoning.Formation) FormationAction
}


func (a *ActionImpl) GetFormationAction(ft FormationType) FormationAction {
	var found FormationAction
	a.DoFormationActions(func(fa FormationAction) bool {
		if fa.ApplicableToFormationType(ft) {
			found = fa
			return false
		}
		return true
	})
	return found
}

func (a *ActionImpl) GetFormationActionFor(f reasoning.Formation) FormationAction {
	t := reflect.TypeOf(f)
	var found FormationAction
	a.DoFormationActions(func(fa FormationAction) bool {
		if t.ConvertibleTo(fa.FormationType()) {
			found = fa
			return false
		}
		return true
	})
	return found
}


var AllActions []Action = []Action{}

func FindAction(actionName string) Action {
	for _, a := range AllActions {
    	if a.Name() == actionName {
			return a
		}
	}
	return nil
}

func defineAction(name string, description string) {
	if FindAction(name) != nil {
		panic(fmt.Sprintf("Attempt to redefine action %s", name))
	}
	AllActions = append(AllActions, &ActionImpl{
		name: name,
		description: description,
		formationActions: []FormationAction{},
	}	)
}


// FormationAction represents how an Action should be performed from a
// specific Formation.
type FormationAction interface {
	Action() Action                        // defimpl:"read action"
	Level() Level                          // defimpl:"read level"
	FormationType() FormationType         // defimpl:"read formationType"
	// DoItFunc is a function that will perform the action.
	DoItFunc() func(reasoning.Formation)   // defimpl:"read doItFunc"
	DoIt(reasoning.Formation)
	String() string
	ApplicableTo(reasoning.Formation) bool
	ApplicableToFormationType(FormationType) bool
	// IdString returns a string suitable for use as an HTML ID
	IdString() string
}


func (fa *FormationActionImpl) String() string {
	return fmt.Sprintf("%s on %s", fa.Action().Name(), fa.FormationType().Name())
}

func (fa *FormationActionImpl) IdString() string {
	return fmt.Sprintf("%s-%s", fa.Action().Name(), fa.FormationType().Name())
}

func (fa *FormationActionImpl) ApplicableToFormationType(ft FormationType) bool {
	return ft.AssignableTo(fa.formationType)
}

func (fa *FormationActionImpl) ApplicableTo(f reasoning.Formation) bool {
	return fa.ApplicableToFormationType(FormationType(reflect.TypeOf(f)))
}

func (fa *FormationActionImpl) DoIt(f reasoning.Formation) {
	if !fa.ApplicableTo(f) {
		panic(fmt.Sprintf("%s doesn't apply to %#v", fa, f))
	}
	fa.doItFunc(f)
}


func defineFormationAction(actionName string, level Level, formationType FormationType, doit func(reasoning.Formation)) {
	a := FindAction(actionName)
	if a == nil {
		a = &ActionImpl{
			name: actionName,
			formationActions: []FormationAction{} }
		AllActions = append(AllActions, a)
	}
	a.AddFormationAction(&FormationActionImpl{
		action: a,
		level: level,
		formationType: formationType,
		doItFunc: doit,
	})
}

