package ptr_test

import (
	"fmt"

	"github.com/cohesian-tools/typy/ptr"
)

func ExampleNewPtr() {
	original := 42
	value := ptr.NewPtr(original)
	*value = 7
	fmt.Println(original, *value)
	// Output: 42 7
}
