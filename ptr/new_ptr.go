package ptr

// New returns a pointer to a copy of t. Reference values retain their
// underlying references; New does not perform a deep copy.
func New[T any](t T) *T {
	return &t
}
