package reasoning

import "goshua/rete"
import "goshua/rete/rule_compiler/runtime"


func MakeFormationsRete() rete.Node {
	root_node := rete.MakeRootNode()
	loadAllRules(root_node)
	bufferAllTypes(root_node)
	return root_node
}

func loadAllRules(root rete.Node) {
	for _, rule := range runtime.AllRules {
		rule.Installer()(root)
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

