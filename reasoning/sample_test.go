package reasoning

import "fmt"
import "os"
import "reflect"
import "testing"
import "goshua/rete"

/*
// Required for debug hook tracing:
import "bytes"
import "defimpl/util"
*/


// TestMakeSample makes sure that each sample formation returned by
// MakeSampleFormation matches the rule for that formation.
func TestMakeSampleFormation(t *testing.T) {
	/*
	This test is failing for TwoFacedLine because of MiniWave symetry:
	  TwoFacedLine(RightHanded, Dancer_1, Dancer_2, Dancer_3, Dancer_4)
	  TwoFacedLine(RightHanded, Dancer_3, Dancer_4, Dancer_1, Dancer_2)
	*/
	/*
	// Example of how to get debugging output:
	interesting := util.MemberTypes(reflect.TypeOf(struct {
		Couple
		*CoupleImpl
		Tandem
		*TandemImpl
		TandemCouples
		*TandemCouplesImpl
	}{}))
	rete.DEBUG_EMIT_HOOK = func(n rete.Node, item interface{}) {
		if n.Label() != "root" {
			return
		}
		if interesting.Contains(reflect.TypeOf(item)) {
			t.Logf("      %q Emitting %T(%s)", n.Label(), item, item.(Formation).Dancers().String())
		}
	}
	rete.DEBUG_RULE_PARAMETER_RECEIVE_HOOK = func(n *rete.RuleParameterNode, item interface{}) {
		if interesting.Contains(reflect.TypeOf(item)) {
			t.Logf("%q Receiving %T(%s)", n.Label(), item, item.(Formation).Dancers().String())
		}
	}
	rete.DEBUG_FILL_AND_CALL_ENTRY_HOOK = func(n *rete.RuleParameterNode, item interface{}, rule_node *rete.RuleNode) {
		if rule_node.Label() == "rule TandemCouples" {
			t.Logf("  fill_and_call %q %T(%s) %q", n.Label(), item, item.(Formation).Dancers().String(), rule_node.Label())
		}
	}
	rete.DEBUG_FILL_AND_CALL_RULE_CALL_HOOK = func(rule_node *rete.RuleNode, parameters []interface{}) {
		if rule_node.Label() == "rule TandemCouples" {
			pstring := bytes.NewBufferString("")
			for i, p := range parameters {
				if i > 0 {
					pstring.WriteString(", ")
				}
				fmt.Fprintf(pstring, "%T(%s)", p, p.(Formation).Dancers().String())
			}
			t.Logf("    fill_and_call calling %q on %s", rule_node.Label(), pstring)
		}
	}
	defer func() {
		rete.DEBUG_EMIT_HOOK = nil
		rete.DEBUG_RULE_PARAMETER_RECEIVE_HOOK = nil
		rete.DEBUG_FILL_AND_CALL_ENTRY_HOOK = nil
		rete.DEBUG_FILL_AND_CALL_RULE_CALL_HOOK = nil
	}()
	*/
	for ft, c := range formation_sample_constructors {
		if ft.Kind() != reflect.Interface {
			t.Errorf("Kind of %v is %s, expected Interface", ft, ft.Kind().String())
		}
		sample := c()
		/*
		// ***** Debug a particular formation type.
		if ty := reflect.TypeOf(sample).Elem().Name(); ty != "TandemCouplesImpl" {
			t.Errorf("skipping %s", ty)
			continue
		}
		*/
		f, ff := FindFormations(sample.Dancers(), ft)
		dot_file := fmt.Sprintf("missing-%s.dot", ft.Name())
		if count := len(f); count != 1 {
			t.Errorf("Expected one formation of type %s, got %d", ft, count)
			graph := rete.GraphMissingConclusion(dot_file, ff.rete, ft.(reflect.Type))
			if graph.Err != nil {
				t.Logf("Error writing dot file: %s", graph.Err)
			} else {
				t.Logf("Wrote %s, %d types, %d rules;\n%#v",
					graph.OutputPath,
					graph.TypeCount,
					graph.RuleCount,
					*graph)
			}
			// Show sample's dancers:
			for _, dancer := range sample.Dancers() {
				t.Logf("  %v: %v %v %v\n",
					dancer,
					dancer.Position().Down,
					dancer.Position().Left,
					dancer.Direction())
			}
			// Show contents of all buffer nodes:
			ff.DoAllBuffers(func (bn rete.AbstractBufferNode) {
				t.Logf("rete Node %s:\n", bn.(rete.Node).Label())
				bn.DoItems(func(item interface{}) {
					t.Logf("    %v\n", item)
				})
			})
		} else {
			// Test didn't fail, so delete any lingering
			// dot file from a previous failure:
			_ = os.Remove(dot_file)
		}
	}
}

