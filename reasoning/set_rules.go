package reasoning

import "goshua/rete"
import "squaredance/dancer"

func rule_SetHasDancers(node rete.Node, s dancer.Set) {
	for _, d := range s.Dancers() {
		node.Emit(d)
	}
}
