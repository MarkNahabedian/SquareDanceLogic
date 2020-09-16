package reasoning

import "reflect"
import "defimpl/runtime"


// formation_sample_constructors is keyed by implementation type since
// that's what we can get by reflecting on the return value of the
// constructor function.
var formation_sample_constructors = map[reflect.Type] func()Formation{}

// MakeSampleFormation returns nil or a sample square fance formation
// of the specified type.
func MakeSampleFormation(formation_type reflect.Type) Formation {
	if formation_type.Kind() == reflect.Interface {
		impl := runtime.InterfaceToImpl(formation_type)
		return MakeSampleFormation(impl)
	}
	constructor, ok := formation_sample_constructors[formation_type]
	if ok {
		return constructor()
	}
	return nil
}

func RegisterFormationSample(constructor func()Formation) {
	sample := constructor()
	// Panic if constructor doesn't return a Formation:
	_ = sample.(Formation)
	formation_sample_constructors[reflect.TypeOf(sample)] = constructor
}

