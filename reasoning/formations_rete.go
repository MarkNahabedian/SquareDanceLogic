package reasoning

import "goshua/rete"

func MakeFormationsRete() rete.Node {
	root_node := rete.MakeRootNode()
	loadAllRules(root_node)
	bufferAllTypes(root_node)
	return root_node
}

