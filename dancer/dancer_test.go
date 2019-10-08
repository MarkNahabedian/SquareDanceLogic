package dancer

import "reflect"
import "sort"
import "testing"
import "goshua/goshua"
import "goshua/equality"

func TestDancerCanEqual(t *testing.T) {
	s := NewSquaredSet(4)
	if _, ok := s.Dancers()[0].(equality.CanEqual); !ok {
		t.Errorf("Dancer not CanEqual")
	}
}

func TestDancerEqual(t *testing.T) {
	s := NewSquaredSet(4)
	if eq, err := goshua.Equal(s.Dancers()[0], s.Dancers()[0]); err != nil {
		t.Errorf("Equal error %s", err)
	} else if !eq {
		t.Errorf("Dancer not equal to itself.")
	}
	if eq, err := goshua.Equal(s.Dancers()[0], s.Dancers()[2]); err != nil {
		t.Errorf("Equal error %s", err)
	} else if eq {
		t.Errorf("Different Dancers are equal.")
	}
}

func TestSquaredSet(t *testing.T) {
	s := NewSquaredSet(4)
	for i, dancer := range s.Dancers() {
		partner := dancer.OriginalPartner()
		if s != dancer.Set() {
			t.Errorf("Dancer's Set is wrong, %d", i)
		}
		if partner.CoupleNumber() != dancer.CoupleNumber() {
			t.Errorf("CoupleNumbers don't match: %d:\n  %v\n  %v", i, dancer, partner)
		}
		if !partner.Direction().Equal(dancer.Direction()) {
			t.Errorf("Directions don't match: %d:\n  %v\n  %v", i, dancer, partner)
		}

	}
}

func TestUnion(t *testing.T) {
	s := NewSquaredSet(4)
	dancers1 := s.Dancers()[1:4]
	dancers2 := s.Dancers()[2:6]
	got := Union(dancers1, dancers2).Ordered()
	if want := s.Dancers()[1:6].Ordered(); !reflect.DeepEqual(got, want) {
		t.Errorf("Union failed: want: %s, got: %s",
			want.String(), got.String())
	}
}

func TestIntersection(t *testing.T) {
	s := NewSquaredSet(4)
	dancers1 := s.Dancers()[1:4]
	dancers2 := s.Dancers()[2:6]
	got := Intersection(dancers1, dancers2).Ordered()
	if want := s.Dancers()[2:4].Ordered(); !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection failed: want: %s, got: %s",
			want.String(), got.String())
	}
}

func TestSetDifference(t *testing.T) {
	s := NewSquaredSet(4)
	dancers1 := s.Dancers()[1:4]
	dancers2 := s.Dancers()[3:6]
	got := SetDifference(dancers1, dancers2).Ordered()
	if want := s.Dancers()[1:3].Ordered(); !reflect.DeepEqual(got, want) {
		t.Errorf("Intersection failed: want: %s, got: %s",
			want.String(), got.String())
	}
}

func TestSpatiallyOrderedDancers(t *testing.T) {
	s := NewSquaredSet(4)
	ordinals := func (ds Dancers) []int {
		result := []int{}
		for _, d := range ds {
			result = append(result, d.Ordinal())
		}
		return result
	}
	dancers := s.Dancers()
	sort.Sort(SpatiallyOrderedDancers(dancers))
	gots := ordinals(dancers)
	wants := []int{ 2, 1, 3, 0, 4, 7, 5, 6 }
	for i, got := range gots {
		want := wants[i]
		if got != want {
			t.Errorf("Spatial ordering failed at %d: want: %d, got: %d",
				i, got, want)
		}
	}
}
