package formationpredicates

import "squaredance/dancer"

type Pair interface {
	Dancer1() dancer.Dancer
	Dancer2() dancer.Dancer
}

func MakePair(dancer1, dancer2 dancer.Dancer) Pair {
	return &pair{ dancer1: dancer1, dancer2: dancer2 }
}

type pair struct {
	dancer1 dancer.Dancer
	dancer2 dancer.Dancer
}

func (p *pair) Dancer1() dancer.Dancer {
	return p.dancer1
}

func (p *pair) Dancer2() dancer.Dancer {
	return p.dancer2
}

