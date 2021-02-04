package reasoning

import "fmt"
import "reflect"
import "defimpl/runtime"
import "squaredance/dancer"
import "squaredance/geometry"


// formation_sample_constructors is keyed by interface type.  It is
// filled by calls to RegisterFormationSample in the output of the
// formation expander.
var formation_sample_constructors = map[FormationType] func()Formation{}

// MakeSampleFormation returns nil or a sample square dance formation
// of the specified type.
func MakeSampleFormation(formation_type FormationType) Formation {
	// What about "Dancers"?  I guess we won't ever make a sample
	// for it.
	if formation_type.Kind() == reflect.Slice {
		return nil
	}
	if formation_type.Kind() == reflect.Struct {
		intr := runtime.ImplToInterface(formation_type)
		if intr == nil {
			return nil
		}
		return MakeSampleFormation(intr)
	}
	constructor, ok := formation_sample_constructors[formation_type]
	if ok {
		return constructor()
	}
	return nil
}

// RegisterFormationSample invokes the formation constructor to
// determine the resulting FormationType and then registers the
// constructor for that formation interface type in
// formation_sample_constructors.
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
