package reasoning

import "testing"
import "squaredance/dancer"
import "goshua/rete"

func init() {
	rete.EnsureTypeTestRegistered("Couple", func(i interface{}) bool { _, ok := i.(Couple); return ok })
	rete.EnsureTypeTestRegistered("MiniWave", func(i interface{}) bool { _, ok := i.(MiniWave); return ok })
}

func TestCouple(t *testing.T) {
	root_node := rete.MakeRootNode()
	SetHasDancersRule(root_node)
	PairOfDancersRule(root_node)
	GeneralizedCoupleRule(root_node)
	MiniWaveRule(root_node)

	// Make sure the assertions get buffered so we can dump them:
	rete.GetBuffered(rete.GetTypeTestNode(root_node, "Couple"))
	rete.GetBuffered(rete.GetTypeTestNode(root_node, "MiniWave"))

	// Show the rete
	rete.Walk(root_node, func(n rete.Node) {
		t.Logf("node %T %s", n, n.Label())
	})
	graph, err := rete.MakeGraph(root_node)
	if err != nil {
		t.Fatal(err)
	}
	rete.WriteGraphvizFile(graph, "TestCouple_rete.dot")

	set := dancer.NewSquaredSet(4)
	root_node.Receive(set)

	rete.Walk(root_node, func(n rete.Node) {
		if n, ok := n.(*rete.BufferNode); ok {
			t.Logf("Dump: %s", n.Label())
			c := n.GetCursor()
			for item, present := c.Next(); present; item, present = c.Next() {
				t.Logf("    %v", item)
			}
		}
	})

	t.Errorf("foo")
}

