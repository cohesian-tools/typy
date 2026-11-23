package errlay_test

import (
	"errors"
	"fmt"

	"github.com/cohesian-tools/typy/errlay"
	"github.com/cohesian-tools/typy/ptr"
)

func ExampleNew() {
	cause := errors.New("connection refused")
	err := errlay.New("database", 10, "query failed", cause)
	fmt.Println(err)
	fmt.Println(errors.Is(err, cause))
	// Output:
	// |layer database - 10| query failed connection refused
	// true
}

func ExampleTarget() {
	err := errlay.New("service", 1, "request failed",
		errlay.New("database", 10, "query failed", nil),
	)
	target := errlay.Target(ptr.NewPtr("database"), ptr.NewPtr(uint8(10)))
	fmt.Println(errors.Is(err, target))
	fmt.Println(errors.Is(err, errlay.Target(nil, nil)))
	fmt.Println(errors.Is(errors.New("ordinary error"), errlay.Target(nil, nil)))
	// Output:
	// true
	// true
	// false
}

func ExampleError() {
	err := fmt.Errorf("context: %w", errlay.New("database", 10, "query failed", nil))
	var detail *errlay.Error
	if errors.As(err, &detail) {
		fmt.Println(detail.Layer(), detail.Code(), detail.Message())
	}
	// Output: database 10 query failed
}
