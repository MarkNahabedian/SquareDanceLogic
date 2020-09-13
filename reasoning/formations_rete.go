package reasoning

import "reflect"
import "goshua/rete"
import "squaredance/dancer"
import "goshua/rete/rule_compiler/runtime"


type FormationFinder struct {
	rete rete.Node    // The root Node
	typeToBuffer map[reflect.Type]rete.AbstractBufferNode
}


func MakeFormationFinder() *FormationFinder {
	ff := &FormationFinder{}
	ff.rete = rete.MakeRootNode()
	ff.typeToBuffer = make(map[reflect.Type]rete.AbstractBufferNode)
	loadAllRules(ff.rete)
	// Add bufferes where needed.  Index the buffers
	rete.Walk(ff.rete, func(n rete.Node) {
		if ttn, ok := n.(*rete.TypeTestNode); ok {
			if ttn.Type.Implements(reflect.TypeOf(func(TwoDancerSymetric){}).In(0)) {
				ff.typeToBuffer[ttn.Type] = rete.GetUniqueBuffered(n, IsTwoDancerSymetric)
			}
		}
	})
	return ff
}


func loadAllRules(root rete.Node) {
	for _, rule := range runtime.AllRules {
		rule.Installer()(root)
	}
}


func (ff *FormationFinder) Injest(dancers dancer.Dancers) {
	for _, dancer := range dancers {
		ff.rete.Receive(dancer)
	}
}


func (ff *FormationFinder) Clear() {
	rete.Walk(ff.rete, rete.Node.Clear)
}


func (ff *FormationFinder) DoFormations(formationType reflect.Type, f func(Formation)) {
	bn := ff.typeToBuffer[formationType]
	bn.DoItems(func (item interface{}) {
		f(item.(Formation))
	})
}

