package reasoning

import "fmt"
import "reflect"
import "goshua/rete"
import "defimpl/runtime"
import "squaredance/dancer"


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
	// Add buffers where needed.  Index the buffers
	rete.Walk(ff.rete, func(n rete.Node) {
		if ttn, ok := n.(*rete.TypeTestNode); ok {
			if ttn.Type.Implements(reflect.TypeOf(func(TwoDancerSymetric){}).In(0)) {
				ff.typeToBuffer[ttn.Type] = rete.GetUniqueBuffered(n, IsTwoDancerSymetric)
			} else {
				ff.typeToBuffer[ttn.Type] = rete.GetBuffered(n)
			}
			if ff.typeToBuffer[ttn.Type] == nil {
				panic(fmt.Sprintf("Can't find buffered node for %v", ttn.Type))
			}
		}
	})
	return ff
}


func loadAllRules(root rete.Node) {
	for _, rule := range rete.AllRules {
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
	formationType1, err := runtime.InterfaceFor(formationType)
	if formationType1 == nil {
		panic(fmt.Sprintf("Can't find interface type for %s: %s", formationType.String(), err))
	}
	bn := ff.typeToBuffer[formationType1]
	if bn == nil {
		panic(fmt.Sprintf("no buffer for %s, %s", formationType.String(), formationType1.String()))
	}
	bn.DoItems(func (item interface{}) {
		f(item.(Formation))
	})
}


// DoAllBuffers applies the specified function to each AbstractBufferNode.
func (ff *FormationFinder) DoAllBuffers(f func(rete.AbstractBufferNode)) {
	rete.Walk(ff.rete, func(node rete.Node) {
		if n, ok := node.(rete.AbstractBufferNode); ok {
			f(n)
		}
	})
}

