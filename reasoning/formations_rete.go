package reasoning

import "fmt"
import "reflect"
import "goshua/rete"
import defimpl_runtime "defimpl/runtime"
import "squaredance/dancer"
import "goshua/rete/rule_compiler/runtime"


type FormationFinder struct {
	rete rete.Node    // The root Node
	typeToBuffer map[reflect.Type]rete.AbstractBufferNode
}


func MakeFormationFinder() *FormationFinder {
	ff := &FormationFinder{
		rete: rete.MakeRootNode(),
		typeToBuffer: make(map[reflect.Type]rete.AbstractBufferNode),
	}
	loadAllRules(ff.rete)
	// Add bufferes where needed.  Index the buffers
	rete.Walk(ff.rete, func(n rete.Node) {
		if ttn, ok := n.(*rete.TypeTestNode); ok {
			if ttn.Type.Implements(reflect.TypeOf(func(TwoDancerSymetric){}).In(0)) {
				ff.typeToBuffer[ttn.Type] = rete.GetUniqueBuffered(n, IsTwoDancerSymetric)
			} else {
				ff.typeToBuffer[ttn.Type] = rete.GetBuffered(n)
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


// DoFormations calls the provided function on each formation that the
// FormationFinder found of the specified FormationType.
func (ff *FormationFinder) DoFormations(formationType reflect.Type, f func(Formation)) {
	if formationType.Kind() != reflect.Interface {
		from := formationType
		formationType = defimpl_runtime.ImplToInterface(formationType)
		if formationType == nil {
			panic(fmt.Sprintf("No interface type corresponding to %v",
				from))
		}
		if formationType.Kind() != reflect.Interface {
			panic(fmt.Sprintf("Not an interface type: %s", formationType))
		}
	}
	bn := ff.typeToBuffer[formationType]
	if bn == nil {
		panic(fmt.Sprintf("no buffer for %s", formationType))
	}
	bn.DoItems(func (item interface{}) {
		f(item.(Formation))
	})
}


