package ptr

// NewPtr returns a pointer to a copy of t. Reference values retain their
// underlying references; NewPtr does not perform a deep copy.
func NewPtr[T any](t T) *T {
	return &t
}
