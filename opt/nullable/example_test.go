package nullable_test

import (
	"fmt"

	"github.com/cohesian-tools/typy/opt/nullable"
)

func ExampleSome() {
	value := nullable.Some(0)
	missing := nullable.None[int]()
	fmt.Println(value.IsNull(), value.Unwrap())
	fmt.Println(missing.IsNull())
	// Output:
	// false 0
	// true
}
