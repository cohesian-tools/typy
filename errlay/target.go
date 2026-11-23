package errlay

// Target returns a field filter for errors.Is. Nil fields are wildcards;
// two nil fields match any *Error. Both specified fields must match the same
// error in the chain. Target retains the supplied pointers.
func Target(layer *string, code *uint8) error {
	return &errTarget{
		layer: layer,
		code:  code,
	}
}

type errTarget struct {
	layer *string
	code  *uint8
}

func (t *errTarget) Error() string {
	return ""
}
