package errlay

import (
	"fmt"
)

const (
	message = "|layer %s - %d| %s"
)

// New creates a layered error wrapping err. A nil cause creates a leaf error.
func New(layer string, code uint8, msg string, err error) error {
	return &Error{
		layer: layer,
		code:  code,
		msg:   msg,
		err:   err,
	}
}

// Error associates a layer, code, and message with an optional underlying cause.
type Error struct {
	layer string
	code  uint8
	msg   string
	err   error
}

// Error formats the metadata and appends the underlying cause when present.
func (e *Error) Error() string {
	msg := fmt.Sprintf(message, e.layer, e.code, e.msg)
	if e.err != nil {
		msg += " " + e.err.Error()
	}
	return msg
}

/* Expose Methods */
// Layer returns the logical layer associated with the error.
func (e *Error) Layer() string {
	return e.layer
}

// Code returns the numeric code associated with the error.
func (e *Error) Code() uint8 {
	return e.code
}

// Message returns the local message without the metadata or underlying cause.
func (e *Error) Message() string {
	return e.msg
}

// Unwrap returns the underlying cause for standard error-chain traversal.
func (e *Error) Unwrap() error {
	return e.err
}

/* Is Chain */
// Is reports whether this error matches a filter constructed by Target.
// Use errors.Is to search the full error chain.
func (e *Error) Is(target error) bool {
	t, ok := target.(*errTarget)
	if !ok {
		return false
	}

	isLayer := t.layer == nil || *t.layer == e.layer
	isCode := t.code == nil || *t.code == e.code
	return isLayer && isCode
}
