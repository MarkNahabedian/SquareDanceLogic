package reasoning

import "fmt"
import "reflect"
import "defimpl/runtime"
import "squaredance/dancer"
import "squaredance/geometry"


// formation_sample_constructors is keyed by implementation type since
// that's what we can get by reflecting on the return value of the
// constructor function.
var formation_sample_constructors = map[FormationType] func()Formation{}

// MakeSampleFormation returns nil or a sample square dance formation
// of the specified type.
func MakeSampleFormation(formation_type FormationType) Formation {
	if formation_type.Kind() == reflect.Slice {
		return nil
	}
	if formation_type.Kind() == reflect.Interface {
		impl := runtime.InterfaceToImpl(formation_type)
		if impl == nil {
			return nil
		}
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
	i, err := runtime.InterfaceFor(reflect.TypeOf(sample))
	if i == nil {
		panic(fmt.Sprintf("InterfaceFor returned nil for %T: %s", sample, err))
	}
	formation_sample_constructors[i] = constructor
}


func init() {
	// Register a constructor for Dancer:
	RegisterFormationSample(func() Formation {
		dancers := dancer.MakeSomeDancers(1)
		dancers[0].Move(geometry.Position{ Left: geometry.Left0, Down: geometry.Down0 },
			geometry.Direction0)
		return dancers[0]
	})
}
