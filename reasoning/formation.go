package reasoning

import "fmt"
import "reflect"
import "squaredance/dancer"

// Formation represents a square dance formation.
type Formation interface {
	// NumberOfDancers returns the number of dancers in the Formation.
	NumberOfDancers() int

	// Dancers returns a slice containing the Dancers in the Formation.
	Dancers() dancer.Dancers

	// HasDancer returns true of the specified Dancer is in the Formation.
	HasDancer(dancer.Dancer) bool
}

// FormationType is a reflect.Type identifying the interface type of a Formation.
type FormationType = reflect.Type

// AllFormationTypes maps from a formation name to the corresponding
// FormationType.  That FormationType is an interface type, not an
// implementation type.
var AllFormationTypes map[string] FormationType = make(map[string] FormationType)

func LookupFormationType(name string) FormationType {
	ft, ok := AllFormationTypes[name]
	if !ok {
		panic(fmt.Sprintf("No formation named %q", name))
	}
	return FormationType(ft)
}

func init() {
	// Fudge the AllFormationTypes entries for Formations that aren't
	// automatically expanded.
	f1 := func(f dancer.Dancer) {}
	t1 := reflect.TypeOf(f1).In(0)
	AllFormationTypes[t1.Name()] = t1
	f2 := func(f dancer.Dancers) {}
	t2 := reflect.TypeOf(f2).In(0)
	AllFormationTypes[t2.Name()] = t2
}

// HasDancers returns true if f contains all of the specified Dancers.
func HasDancers(f Formation, dancers ...dancer.Dancer) bool {
	for _, d := range dancers {
		if !f.HasDancer(d) {
			return false
		}
	}
	return true
}

// OrderedDancers returns true if all of the Dancers are in
// ascending order by Ordinal.
func OrderedDancers(dancers ...dancer.Dancer) bool {
	for i := 0; i < len(dancers) - 1; i++ {
		if dancers[i].Ordinal() >= dancers[i+1].Ordinal() {
			return false
		}
	}
	return true
}

