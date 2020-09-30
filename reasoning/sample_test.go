package reasoning

import "testing"
import "defimpl/runtime"
import "goshua/rete"


// TestMakeSample makes sure that each sample formation returned by
// MakeSampleFormation matches the rule for that formation.
func TestMakeSampleFormation(t *testing.T) {
	runtime.Dump()
	/*
	ff := MakeFormationFinder()
	for bft, bn := range ff.typeToBuffer {
		t.Logf("Buffer for %s: \t%T", bft, bn)
	}
	for ft, _ := range formation_sample_constructors {
		t.Logf("Constructor for %s", ft)
	}
	*/
	for ft, c := range formation_sample_constructors {
		sample := c()
		// t.Logf("%s %v", ft, sample)
		f, ff := FindFormations(sample.Dancers(), ft)
		if count := len(f); count != 1 {
			t.Errorf("Expected one formation of type %s, got %d", ft, count)
			for _, dancer := range sample.Dancers() {
				t.Logf("  %d: %v %v %v\n",
					dancer.Ordinal(),
					dancer.Position().Down,
					dancer.Position().Left,
					dancer.Direction())
			}
			ff.DoAllBuffers(func (bn rete.AbstractBufferNode) {
				t.Logf("rete Node %s:\n", bn.(rete.Node).Label())
				bn.DoItems(func(item interface{}) {
					t.Logf("    %v\n", item)
				})
			})
		}
	}
}

