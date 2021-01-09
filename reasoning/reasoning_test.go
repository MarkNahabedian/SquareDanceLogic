package reasoning

import "fmt"
import "os"
import "reflect"
import "strings"
import "testing"
// import "squaredance/dancer"
import "goshua/rete"


func showAllAssertions(t *testing.T, root rete.Node) {
	rete.Walk(root, func(n rete.Node) {
		if n, ok := n.(*rete.BufferNode); ok {
			t.Logf("Dump: %s", n.Label())
			n.DoItems(func (item interface{}) {
				t.Logf("    %s", item)
			})
		}
	})
}

func TestWriteDiagrams(t *testing.T) {
	WriteFormationDiagrams()
}

func TestShowFullRete(t *testing.T) {
	formation_finder := MakeFormationFinder()
	// Show the rete
	rete.Walk(formation_finder.rete, func(n rete.Node) {
		t.Logf("node %T %s", n, n.Label())
	})
	graph, err := rete.MakeGraph(formation_finder.rete)
	if err != nil {
		t.Fatal(err)
	}
	rete.WriteGraphvizFile(graph, "formations_rete.dot")
}


func TestAllRules(t *testing.T) {
	typesPretty := func(ts []reflect.Type) string {
		s := make([]string, len(ts))
		for i, t1 := range ts {
			s[i] = t1.String()
		}
		return strings.Join(s, ", ")
	}
	filename := "all_rules.txt"
	out, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't open %s: %s", filename, err)
	}
	defer out.Close()
	for _, r := range rete.AllRules {
		out.WriteString(fmt.Sprintf("%s \n\t%s\n\t%s\n\n",
			r.Name(),
			typesPretty(r.ParamTypes()),
			typesPretty(r.EmitTypes())))
	}
}

func TestAllFormations(t *testing.T) {
	filename := "formation_types.txt"
	out, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't open %s: %s", filename, err)
	}
	defer out.Close()
	for name, typ := range AllFormationTypes {
		out.WriteString(fmt.Sprintf("%s\t  %v\n", name, typ))
	}
}

/*
func TestCouple(t *testing.T) {
	root_node := MakeFormationsRete()

	set := dancer.NewSquaredSet(4)
	root_node.Receive(set)

	showAllAssertions(t, root_node)
	t.Errorf("foo")
}
*/

/*
func TestTwoFacedLines(t *testing.T) {
	root_node := MakeFormationsRete()

}
*/

