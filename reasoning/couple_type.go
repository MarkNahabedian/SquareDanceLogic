package reasoning

import "squaredance/dancer"

// ??? Should these be methods on CoupleImpl and part of the Couple
// interface instead?

// IsNormal returns true of a normal Couple.
// A couple is normal if the Beau (dancer on the left) is a Guy and
// the Belle (dancer on the right) is a Guy.
func IsNormal(c Couple) bool {
	if c.Beau().Gender() == dancer.Guy &&
		c.Belle().Gender() == dancer.Gal {
		return true
	}
	return false
}

// IsSasheyed returns true if Couple is half-sasheyed.
// A half-sasheyed couple has a Gal in Beau position and a Guy in
// Belle position.
func IsSasheyed(c Couple) bool {
	if c.Beau().Gender() == dancer.Gal &&
		c.Belle().Gender() == dancer.Guy {
		return true
	}
	return false
}

// IsSameGender returns true if both dancers of the Couple
// have the same Gender.
func IsSameGender(c Couple) bool {
	return c.Beau().Gender().Equal(c.Belle().Gender())
}
