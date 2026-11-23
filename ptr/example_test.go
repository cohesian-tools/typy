package ptr_test

import (
	"fmt"

	"github.com/cohesian-tools/typy/ptr"
)

func ExampleNew() {
	original := 42
	value := ptr.New(original)
	*value = 7
	fmt.Println(original, *value)
	// Output: 42 7
}
