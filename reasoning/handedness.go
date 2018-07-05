package reasoning


// Handedness represents the handedness of a square dance formation.
type Handedness int

const (
	// Unspecified is used when the Gender is unknown or doesn't matter.
	NoHanded Handedness = iota
	RightHanded
	LeftHanded
)

func (h Handedness) String() string {
	switch h {
		case NoHanded: return "NoHanded"
		case RightHanded: return "RightHanded"
		case LeftHanded: return "LeftHanded"
	}
	panic("Unsupported handedness")
}

func (h Handedness) Opposite() Handedness {
	switch h {
		case NoHanded:		return NoHanded
		case RightHanded:	return LeftHanded
		case LeftHanded:	return RightHanded
	}
	panic("Unsupported handedness")
}

type HasHandedness interface {
	Handedness() Handedness
}
