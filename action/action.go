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


type Action interface {
	Name() string
	Description() string
	GetFormationAction(FormationType) FormationAction
	GetFormationActionFor(f reasoning.Formation) FormationAction
	AddFormationAction(FormationAction)
}

type implAction struct {
	name string
	description string
	formationActions []FormationAction
}

func (a *implAction) Name() string { return a.name }
func (a *implAction) Description() string { return a.description }

func (a *implAction) GetFormationAction(ft FormationType) FormationAction {
	for _, fa := range a.formationActions {
		if fa.ApplicableToFormationType(ft) {
			return fa
		}
	}
	return nil
}

func (a *implAction) GetFormationActionFor(f reasoning.Formation) FormationAction {
	t := reflect.TypeOf(f)
	for _, fa := range a.formationActions {
		if t.Implements(fa.FormationType()) {
			return fa
		}
	}
	return nil
}

func (a *implAction) AddFormationAction(fa FormationAction) {
	if a.GetFormationAction(fa.FormationType()) != nil {
		panic(fmt.Sprintf("Attempt to redefine %s", fa))
	}
	a.formationActions = append(a.formationActions, fa)
}


func (a *implAction) FormationActions() []FormationAction {
	return a.formationActions
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
	AllActions = append(AllActions, &implAction{
		name: name,
		description: description,
		formationActions: []FormationAction{},
	}	)
}


type FormationAction interface {
	Action() Action
	FormationType() FormationType
	ApplicableTo(reasoning.Formation) bool
	ApplicableToFormationType(FormationType) bool
	DoIt(reasoning.Formation)
}


type implFormationAction struct {
	action Action
	formationType FormationType
	doItFunc func(reasoning.Formation)
}

func (fa *implFormationAction) Action() Action { return fa.action }
func (fa *implFormationAction) FormationType() FormationType { return fa.formationType }

func (fa *implFormationAction) String() string {
	return fmt.Sprintf("%s on %s", fa.Action().Name(), fa.FormationType().Name())
}

func (fa *implFormationAction) ApplicableToFormationType(ft FormationType) bool {
	return ft.Implements(fa.formationType)
}

func (fa *implFormationAction) ApplicableTo(f reasoning.Formation) bool {
	return fa.ApplicableToFormationType(FormationType(reflect.TypeOf(f)))
}

func (fa *implFormationAction) DoIt(f reasoning.Formation) {
	if !fa.ApplicableTo(f) {
		panic(fmt.Sprintf("%s doesn't apply to %#v", fa, f))
	}
	fa.doItFunc(f)
}


func defineFormationAction(actionName string, formationType FormationType, doit func(reasoning.Formation)) {
	a := FindAction(actionName)
	if a == nil {
		a = &implAction{ name: actionName, formationActions: []FormationAction{} }
		AllActions = append(AllActions, a)
	}
	a.AddFormationAction(&implFormationAction{
		action: a,
		formationType: formationType,
		doItFunc: doit,
	})
}

