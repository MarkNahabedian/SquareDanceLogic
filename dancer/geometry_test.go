package dancer

import "testing"


func TestGeometry(t *testing.T) {
	set := NewSquaredSet(4)
	if want, got := set.FlagpoleCenter(), set.Dancers().Center(); got != want {
		t.Errorf("Center failed: got %v, want %v", got, want)
	}
	leftmost, rightmost, downmost, upmost := set.Dancers().Bounds()
	// Since a squared set is symetric around its center:
	fpc := set.FlagpoleCenter()
	if (leftmost + rightmost) / 2 != fpc.Left {
		t.Errorf("Average of %f and %f != %f", leftmost, rightmost, fpc.Left)
	}
	if (downmost + upmost) / 2 != fpc.Down {
		t.Errorf("Average of %f and %f != %f", downmost, upmost, fpc.Down)
	}
}

