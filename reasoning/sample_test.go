package reasoning

import "testing"
import "defimpl/runtime"


// TestMakeSample makes sure that each sample formation returned by
// MakeSampleFormation matches the rule for that formation.
func TestMakeSampleFormation(t *testing.T) {
	runtime.Dump()
	ff := MakeFormationFinder()
	for bft, bn := range ff.typeToBuffer {
		t.Logf("Buffer for %s: \t%T", bft, bn)
	}
	for ft, _ := range formation_sample_constructors {
		t.Logf("Constructor for %s", ft)
	}

	for ft, c := range formation_sample_constructors {
		sample := c()
		t.Logf("%s %v", ft, sample)
		f := FindFormations(sample.Dancers(), ft)
		if len(f) != 1 {
			t.Errorf("Expected one formation of type %s", ft)
		}
	}
}

