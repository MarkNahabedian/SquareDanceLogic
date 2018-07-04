package reasoning


// Handedness represents the handedness of a square dance formation.
type Handedness int

const (
	// Unspecified is used when the Gender is unknown or doesn't matter.
	NoHanded Handedness = iota
	RightHanded
	LeftHanded
)

func (h Handedness) Opposite() Handedness {
	switch h {
		case NoHanded:		return NoHanded
		case RightHanded:	return LeftHanded
		case LeftHanded:	return RightHanded
	}
	return NoHanded
}

type HasHandedness interface {
	Handedness() Handedness
}
