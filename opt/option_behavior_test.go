package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/cohesian-tools/typy/opt"
)

func TestOptionZeroValueAndNilPointer(t *testing.T) {
	var empty opt.Option[int]
	if empty.IsSome() {
		t.Fatal("a zero Option must be None")
	}

	presentNil := opt.Some[*int](nil)
	if !presentNil.IsSome() || presentNil.Unwrap() != nil {
		t.Fatal("Some of a nil pointer must be present and unwrap to nil")
	}

	data, err := json.Marshal(presentNil)
	if err != nil || string(data) != "null" {
		t.Fatalf("Marshal(Some(nil)) = %s, %v; want null, nil", data, err)
	}
	if err := json.Unmarshal(data, &presentNil); err != nil {
		t.Fatal(err)
	}
	if presentNil.IsSome() {
		t.Fatal("JSON null must decode as None")
	}
}

func TestOptionNoneUnwrapPanics(t *testing.T) {
	defer func() {
		if got := recover(); got != "unable to unwrap None" {
			t.Fatalf("Unwrap panic = %#v; want string %q", got, "unable to unwrap None")
		}
	}()
	value := opt.None[int]()
	value.Unwrap()
}
