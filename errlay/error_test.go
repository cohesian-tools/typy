package errlay

import (
	"errors"
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	cause := errors.New("connection refused")
	tests := []struct {
		name  string
		layer string
		code  uint8
		msg   string
		cause error
	}{
		{name: "leaf", layer: "database", code: 42, msg: "query failed"},
		{name: "composite", layer: "service", code: 7, msg: "operation failed", cause: cause},
		{name: "zero values", layer: "", code: 0, msg: ""},
		{name: "maximum code", layer: "api", code: 255, msg: "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			layer, code, msg, cause := tt.layer, tt.code, tt.msg, tt.cause

			// Act
			result := New(layer, code, msg, cause)

			// Assert
			got, ok := result.(*Error)
			if !ok {
				t.Fatalf("New() returned %T, want *Error", result)
			}
			if got.layer != layer || got.code != code || got.msg != msg || got.err != cause {
				t.Errorf("New() = %+v, want layer=%q code=%d msg=%q cause=%v", got, layer, code, msg, cause)
			}
		})
	}
}

func TestError(t *testing.T) {
	cause := errors.New("connection refused")
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "leaf",
			err:  &Error{layer: "database", code: 42, msg: "query failed"},
			want: "|layer database - 42| query failed",
		},
		{
			name: "composite",
			err:  &Error{layer: "service", code: 7, msg: "operation failed", err: cause},
			want: "|layer service - 7| operation failed connection refused",
		},
		{
			name: "nested",
			err: &Error{layer: "api", code: 1, msg: "request failed", err: &Error{
				layer: "database", code: 42, msg: "query failed",
			}},
			want: "|layer api - 1| request failed |layer database - 42| query failed",
		},
		{
			name: "zero values",
			err:  &Error{},
			want: "|layer  - 0| ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err := tt.err

			// Act
			got := err.Error()

			// Assert
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLayer(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{name: "nonempty", err: &Error{layer: "database"}, want: "database"},
		{name: "empty", err: &Error{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err := tt.err
			// Act
			got := err.Layer()
			// Assert
			if got != tt.want {
				t.Errorf("Layer() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCode(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want uint8
	}{
		{name: "zero", err: &Error{code: 0}, want: 0},
		{name: "nonzero", err: &Error{code: 42}, want: 42},
		{name: "maximum", err: &Error{code: 255}, want: 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err := tt.err
			// Act
			got := err.Code()
			// Assert
			if got != tt.want {
				t.Errorf("Code() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{name: "nonempty", err: &Error{msg: "query failed"}, want: "query failed"},
		{name: "empty", err: &Error{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err := tt.err
			// Act
			got := err.Message()
			// Assert
			if got != tt.want {
				t.Errorf("Message() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	cause := errors.New("underlying failure")
	inner := &Error{layer: "database", code: 42, err: cause}
	tests := []struct {
		name string
		err  *Error
		want error
	}{
		{name: "leaf", err: &Error{}, want: nil},
		{name: "standard cause", err: &Error{err: cause}, want: cause},
		{name: "nested errlay cause", err: &Error{err: inner}, want: inner},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err := tt.err
			// Act
			got := err.Unwrap()
			// Assert
			if got != tt.want {
				t.Errorf("Unwrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIs(t *testing.T) {
	database, service := "database", "service"
	code42, code7, zeroCode := uint8(42), uint8(7), uint8(0)
	tests := []struct {
		name   string
		err    *Error
		target error
		want   bool
	}{
		{name: "layer matches", err: &Error{layer: database, code: 7}, target: &errTarget{layer: &database}, want: true},
		{name: "layer differs", err: &Error{layer: database}, target: &errTarget{layer: &service}},
		{name: "code matches", err: &Error{layer: service, code: 42}, target: &errTarget{code: &code42}, want: true},
		{name: "code differs", err: &Error{code: 7}, target: &errTarget{code: &code42}},
		{name: "both match", err: &Error{layer: database, code: 42}, target: &errTarget{layer: &database, code: &code42}, want: true},
		{name: "layer mismatch with matching code", err: &Error{layer: service, code: 42}, target: &errTarget{layer: &database, code: &code42}},
		{name: "code mismatch with matching layer", err: &Error{layer: database, code: 7}, target: &errTarget{layer: &database, code: &code42}},
		{name: "both differ", err: &Error{layer: service, code: 7}, target: &errTarget{layer: &database, code: &code42}},
		{name: "zero code is a valid filter", err: &Error{layer: service, code: 0}, target: &errTarget{code: &zeroCode}, want: true},
		{name: "other code exists", err: &Error{code: 42}, target: &errTarget{code: &code7}},
		{name: "ordinary error target", err: &Error{layer: database, code: 42}, target: errors.New("database")},
		{name: "another errlay error is not a pattern", err: &Error{layer: database, code: 42}, target: &Error{layer: database, code: 42}},
		{name: "nil target", err: &Error{layer: database}, target: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err, target := tt.err, tt.target
			// Act
			got := err.Is(target)
			// Assert
			if got != tt.want {
				t.Errorf("Is(%T) = %t, want %t", target, got, tt.want)
			}
		})
	}
}

func TestErrorsIs(t *testing.T) {
	database, service := "database", "service"
	code42, code7 := uint8(42), uint8(7)
	sentinel := errors.New("underlying failure")
	identity := New("service", 7, "identity", nil)

	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "outer layer",
			err:    New("service", 7, "failed", nil),
			target: &errTarget{&service, nil},
			want:   true,
		},
		{
			name:   "inner layer",
			err:    New("api", 1, "failed", New("database", 42, "failed", sentinel)),
			target: &errTarget{&database, nil},
			want:   true,
		},
		{
			name:   "inner code",
			err:    New("api", 1, "failed", New("database", 42, "failed", sentinel)),
			target: &errTarget{nil, &code42},
			want:   true,
		},
		{
			name:   "both fields match the same error",
			err:    New("api", 1, "failed", New("database", 42, "failed", nil)),
			target: &errTarget{&database, &code42},
			want:   true,
		},
		{
			name:   "fields match different errors - no match",
			err:    New("service", 42, "failed", New("database", 7, "failed", nil)),
			target: &errTarget{&database, &code42},
		},
		{
			name:   "unmatched layer",
			err:    New("service", 7, "failed", nil),
			target: &errTarget{&database, nil},
		},
		{
			name:   "unmatched code",
			err:    New("service", 7, "failed", nil),
			target: &errTarget{nil, &code42},
		},
		{
			name:   "existing sentinel in cause chain",
			err:    New("service", 7, "failed", New("database", 42, "failed", sentinel)),
			target: sentinel,
			want:   true,
		},
		{
			name:   "wrapped errlay error",
			err:    fmt.Errorf("outer context: %w", New("database", 42, "failed", nil)),
			target: &errTarget{&database, &code42},
			want:   true,
		},
		{
			name:   "different error with identical fields",
			err:    New("database", 42, "failed", nil),
			target: New("database", 42, "failed", nil),
		},
		{
			name:   "exact same error instance",
			err:    identity,
			target: identity,
			want:   true,
		},
		{
			name:   "nil target",
			err:    New("service", 7, "failed", nil),
			target: nil,
		},
		{
			name:   "nil error",
			err:    nil,
			target: &errTarget{&service, &code7},
		},
		{
			name:   "nil error and target",
			err:    nil,
			target: nil,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err, target := tt.err, tt.target
			// Act
			got := errors.Is(err, target)
			// Assert
			if got != tt.want {
				t.Errorf("errors.Is(%v, %v) = %t, want %t", err, target, got, tt.want)
			}
		})
	}
}

func TestErrorsAs(t *testing.T) {
	direct := &Error{layer: "database", code: 42}
	inner := &Error{layer: "database", code: 42}
	outer := &Error{layer: "api", code: 1, err: inner}

	tests := []struct {
		name string
		err  error
		want *Error
	}{
		{name: "direct errlay error", err: direct, want: direct},
		{name: "wrapped error", err: fmt.Errorf("context: %w", direct), want: direct},
		{name: "outer error takes precedence", err: outer, want: outer},
		{name: "ordinary error", err: errors.New("other")},
		{name: "nil error", err: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			err := tt.err
			var got *Error

			// Act
			ok := errors.As(err, &got)

			// Assert
			if ok != (tt.want != nil) {
				t.Fatalf("errors.As() = %t, want %t", ok, tt.want != nil)
			}
			if got != tt.want {
				t.Errorf("errors.As() captured %p, want %p", got, tt.want)
			}
		})
	}
}
