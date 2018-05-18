package formationpredicates

import "fmt"
import "squaredance/dancer"
// import "squaredance/geometry"

// NormalCouple is a couple with the guy on the left and the gal on the right.

// SasheyedCouple is a half-sasheyed couple.

// Couple is a generalized couple -- two dancers next to each other with
// the same facing direction.


// Do we want to consider all two-dancer relationships in this file?

type NormalCouple struct {
	Beau dancer.Dancer
	Belle dancer.Dancer
}

func (c *NormalCouple) String() string {
	return fmt.Sprintf("NormalCouple(%s, %s)", c.Beau, c.Belle)
}

type SasheyedCouple struct {
	Beau dancer.Dancer
	Belle dancer.Dancer
}

func (c *SasheyedCouple) String() string {
	return fmt.Sprintf("SasheyedCouple(%s, %s)", c.Beau, c.Belle)
}
