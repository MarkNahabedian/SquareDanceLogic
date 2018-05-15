package dancer

import "testing"
import "goshua/goshua"
import "goshua/equality"

func TestDancerCanEqual(t *testing.T) {
	s := NewSquaredSet(4)
	if _, ok := s.Dancers[0].(equality.CanEqual); !ok {
		t.Errorf("Dancer not CanEqual")
	}
}

func TestDancerEqual(t *testing.T) {
	s := NewSquaredSet(4)
	if eq, err := goshua.Equal(s.Dancers[0], s.Dancers[0]); err != nil {
		t.Errorf("Equal error %s", err)
	} else if !eq {
		t.Errorf("Dancer not equal to itself.")
	}
	if eq, err := goshua.Equal(s.Dancers[0], s.Dancers[2]); err != nil {
		t.Errorf("Equal error %s", err)
	} else if eq {
		t.Errorf("Different Dancers are equal.")
	}
}

func TestSquaredSet(t *testing.T) {
	s := NewSquaredSet(4)
	for i, dancer := range s.Dancers {
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
