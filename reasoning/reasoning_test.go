package reasoning

import "fmt"
import "os"
import "strings"
import "testing"
// import "squaredance/dancer"
import "goshua/rete"
import "goshua/rete/rule_compiler/runtime"


func loadAllRules(root rete.Node) {
	for _, rule := range runtime.AllRules {
		rule.Inserter()(root)
	}
}

func bufferAllTypes(root rete.Node) {
	// Make sure the assertions get buffered so we can dump them:
	rete.Walk(root, func(n rete.Node) {
		if _, ok := n.(*rete.TypeTestNode); ok {
			_ = rete.GetBuffered(n)
		}
	})
}

func showAllAssertions(t *testing.T, root rete.Node) {
	rete.Walk(root, func(n rete.Node) {
		if n, ok := n.(*rete.BufferNode); ok {
			// Skip this Node it its input is a Join
			if _, ok := n.Inputs()[0].(*rete.JoinNode); ok {
				return
			} 
			t.Logf("Dump: %s", n.Label())
			c := n.GetCursor()
			for item, present := c.Next(); present; item, present = c.Next() {
				t.Logf("    %s", item)
			}
		}
	})
}


func init() {
	rete.EnsureTypeTestRegistered("Couple", func(i interface{}) bool { _, ok := i.(Couple); return ok })
	rete.EnsureTypeTestRegistered("MiniWave", func(i interface{}) bool { _, ok := i.(MiniWave); return ok })
}

func TestShowFullRete(t *testing.T) {
	root_node := rete.MakeRootNode()
	loadAllRules(root_node)
	bufferAllTypes(root_node)
	// Show the rete
	rete.Walk(root_node, func(n rete.Node) {
		t.Logf("node %T %s", n, n.Label())
	})
	graph, err := rete.MakeGraph(root_node)
	if err != nil {
		t.Fatal(err)
	}
	rete.WriteGraphvizFile(graph, "formations_rete.dot")
}


func TestAllRules(t *testing.T) {
	filename := "all_rules.txt"
	out, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Can't open %s: %s", filename, err)
	}
	defer out.Close()
	for _, r := range runtime.AllRules {
		out.WriteString(fmt.Sprintf("%s \n\t%s\n\t%s\n\n",
			r.Name(),
			strings.Join(r.ParamTypes(), ", "),
			strings.Join(r.EmitTypes(), ", ")))
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
	root_node := rete.MakeRootNode()
	loadAllRules(root_node)
	bufferAllTypes(root_node)

	set := dancer.NewSquaredSet(4)
	root_node.Receive(set)

	showAllAssertions(t, root_node)
	t.Errorf("foo")
}
*/

/*
func TestTwoFacedLines(t *testing.T) {
	root_node := rete.MakeRootNode()
	loadAllRules(root_node)

}
*/

