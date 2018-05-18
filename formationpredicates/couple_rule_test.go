package formationpredicates

// Experiments to define rules and hand translate them into a working
// rete.

import "testing"
import "reflect"

import "goshua/goshua"
import "goshua/rete"

import _ "goshua/variables"
import _ "goshua/bindings"
import _ "goshua/equality"
import _ "goshua/unification"
import _ "goshua/query"

import "squaredance/dancer"


var tDancer reflect.Type

func init() {
	// var d dancer.Dancer
	d := dancer.NewSquaredSet(4).Dancers[0]
	// Dancer is already a pointer type.
	tDancer = reflect.TypeOf(d)
}


func MakeRootNode() rete.Node {
	return rete.MakeFunctionNode("root", func(n rete.Node, item interface{}) {
		n.Emit(item)
	})
}

func MakeQueryNode(label string, q goshua.Query) rete.Node {
	return rete.MakeFunctionNode(label, func(n rete.Node, item interface{}) {
		goshua.Unify(q, item, goshua.EmptyBindings(),
			func(b goshua.Bindings) {
				n.Emit(b)
			})
	})
}

// ***** FunctionNode might replace both TestNode and ActionNode.


func TestDancerQuery(t *testing.T) {
	s := goshua.NewScope()
	v := func(name string) goshua.Variable {
		return s.Lookup(name)
	}
	q := goshua.NewQuery(tDancer, v("dancer1"), map[string]interface{}{
		// "Set": v("set"),
		"Gender": dancer.Guy,
		"Direction": v("Direction"),
		// "Position": v("Position1"),
	})
	set := dancer.NewSquaredSet(4)
	var b1 goshua.Bindings
	goshua.Unify(q, set.Dancers[0], goshua.EmptyBindings(), func(b goshua.Bindings) {
		b1 = b
	})
	if b1 == nil {
		t.Errorf("Failed to find GUY")
	}
	var b2 goshua.Bindings
	goshua.Unify(q, set.Dancers[1], goshua.EmptyBindings(), func(b goshua.Bindings) {
		b2 = b
	})
	if b2 != nil {
		t.Errorf("Found GAL")
	}
}


func Test1(t *testing.T) {
	s := goshua.NewScope()
	v := func(name string) goshua.Variable {
		return s.Lookup(name)
	}

	root_node := MakeRootNode()
	guy := MakeQueryNode("guy",
		goshua.NewQuery(tDancer, v("dancer1"), map[string]interface{}{
	    	// "Set": v("set"),
			"Gender": dancer.Guy,
			"Direction": v("Direction"),
			// "Position": v("Position1"),
		}))

	gal := MakeQueryNode("gal",
		goshua.NewQuery(tDancer, v("dancer2"), map[string]interface{}{
	    	// "Set": v("set"),               // Same set
			"Gender": dancer.Gal,
			"Direction": v("Direction"),   // Same facing direction
			// "Position": v("Position2"),
		}))

	rete.Connect(root_node, guy)
	rete.Connect(root_node, gal)
	jn := rete.Join("join dancers", guy, gal)

	makePair := rete.MakeFunctionNode("Make Pair", func(n rete.Node, joined interface{}) {
		j := joined.([2]interface{})
		dancer1, _ := j[0].(goshua.Bindings).Get(v("dancer1"))
        dancer2, _ := j[1].(goshua.Bindings).Get(v("dancer2"))
		eq, _ := goshua.Equal(dancer1, dancer2)
		if !eq {
			p:= MakePair(dancer1.(dancer.Dancer), dancer2.(dancer.Dancer))
			// t.Logf("Pair %#v", p)
			n.Emit(p)
		}
	})
	rete.Connect(jn, makePair)

	near := rete.MakeFunctionNode("Near", func(n rete.Node, paired interface{}) {
		pair := paired.(Pair)
        if Near(pair.Dancer1(), pair.Dancer2()) {
			n.Emit(pair)
		}
	})
	rete.Connect(makePair, near)

	// The functions in the normal and sasheyed nodes below assume that
	// dancer genders are as specified above.

	normal := rete.MakeFunctionNode("normal couple", func(n rete.Node, paired interface{}) {
		pair := paired.(Pair)
		if RightOf(pair.Dancer1(), pair.Dancer2()) &&
			LeftOf(pair.Dancer2(), pair.Dancer1()) {
			n.Emit(&NormalCouple{ Beau: pair.Dancer1(), Belle: pair.Dancer2() })
		}
	})
	rete.Connect(near, normal)

	sasheyed := rete.MakeFunctionNode("sasheyed couple", func(n rete.Node, paired interface{}) {
		pair := paired.(Pair)
		if RightOf(pair.Dancer2(), pair.Dancer1()) &&
			LeftOf(pair.Dancer1(), pair.Dancer2()) {
			n.Emit(&SasheyedCouple{ Beau: pair.Dancer2(), Belle: pair.Dancer1() })
		}
	})
	rete.Connect(near, sasheyed)

	show := rete.MakeFunctionNode("show conclusions", func(n rete.Node, item interface{}) {
		t.Logf("Concluded %s", item)
	})
	rete.Connect(normal, show)
	rete.Connect(sasheyed, show)

/*
	rete.Walk(root_node, func(n rete.Node) {
		t.Logf(`node "%s" %T`, n.Label(), n)
	})
*/

	Walk(root_node, Node.IsValid)

	if g, err := rete.MakeGraph(root_node); err != nil {
		t.Errorf("Can't graphviz: %s", err)
	} else {
		err = rete.WriteGraphvizFile(g, "Test1.dot")
		if err != nil {
        	t.Errorf("Couldn't write dot file.")
		}
	}

	squared_set := dancer.NewSquaredSet(4)
	for _, dancer := range squared_set.Dancers {
		// t.Logf("injecting %s", dancer)
		root_node.Receive(dancer)
	}

	t.Logf("%d %d",
		jn.Inputs()[0].Inputs()[0].(*rete.BufferNode).Count(),
		jn.Inputs()[1].Inputs()[0].(*rete.BufferNode).Count())
	t.Errorf("force show")
}

