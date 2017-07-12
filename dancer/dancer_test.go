package dancer

import "testing"

func TestSquaredSet(t *testing.T) {
	s := NewSquaredSet(4)
	for i, dancer := range s.Dancers {
		partner := dancer.OriginalPartner
		if s != dancer.Set {
			t.Errorf("Dancer's Set is wrong, %d", i)
		}
		if partner.CoupleNumber != dancer.CoupleNumber {
			t.Errorf("CoupleNumbers don't match: %d:\n  %v\n  %v", i, dancer, partner)
		}
		if !partner.Direction.Equal(dancer.Direction) {
			t.Errorf("Directions don't match: %d:\n  %v\n  %v", i, dancer, partner)
		}

	}
}
