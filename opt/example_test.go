package opt_test

import (
	"encoding/json"
	"fmt"

	"github.com/cohesian-tools/typy/opt"
	"go.mongodb.org/mongo-driver/bson"
)

func ExampleSome() {
	value := opt.Some(0)
	missing := opt.None[int]()
	fmt.Println(value.IsSome(), value.Unwrap())
	fmt.Println(missing.IsSome())
	// Output:
	// true 0
	// false
}

func ExampleOption_MarshalJSON() {
	for _, value := range []opt.Option[int]{opt.Some(42), opt.None[int]()} {
		data, err := json.Marshal(value)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
	}
	// Output:
	// 42
	// null
}

func ExampleOption_UnmarshalJSON() {
	var payload struct {
		Name opt.Option[string] `json:"name"`
	}
	for _, input := range []string{`{"name":"Lucas"}`, `{}`, `{"name":null}`} {
		if err := json.Unmarshal([]byte(input), &payload); err != nil {
			panic(err)
		}
		fmt.Println(payload.Name.IsSome())
	}
	// Output:
	// true
	// true
	// false
}

func ExampleOption_MarshalBSONValue() {
	type document struct {
		Name opt.Option[string] `bson:"name"`
	}
	data, err := bson.Marshal(document{Name: opt.Some("Lucas")})
	if err != nil {
		panic(err)
	}
	var decoded document
	if err := bson.Unmarshal(data, &decoded); err != nil {
		panic(err)
	}
	fmt.Println(decoded.Name.Unwrap())
	// Output: Lucas
}
