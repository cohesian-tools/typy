package ptr_test

import (
	"testing"

	"github.com/cohesian-tools/typy/ptr"
)

func TestNewValues(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		checkNew(t, 42)
		checkNew(t, 0)
		checkNew(t, -42)
	})
	t.Run("string", func(t *testing.T) {
		checkNew(t, "hello")
		checkNew(t, "")
	})
	t.Run("bool", func(t *testing.T) {
		checkNew(t, true)
		checkNew(t, false)
	})
	t.Run("named type", func(t *testing.T) {
		type UserID int
		checkNew(t, UserID(42))
	})
	t.Run("struct", func(t *testing.T) {
		type User struct{ Name string }
		checkNew(t, User{Name: "Lucas"})
		checkNew(t, User{})
	})
}

func checkNew[T comparable](t *testing.T, value T) {
	t.Helper()
	var got *T = ptr.New(value)
	if got == nil {
		t.Fatal("New returned nil")
	}
	if *got != value {
		t.Fatalf("*New(%#v) = %#v; want %#v", value, *got, value)
	}
}

func TestNewIndependentStorage(t *testing.T) {
	value := 42
	first := ptr.New(value)
	second := ptr.New(value)
	if first == &value || second == &value || first == second {
		t.Fatal("New must allocate independent storage for each call")
	}

	value = 100
	if *first != 42 || *second != 42 {
		t.Fatal("changing the input changed a returned value")
	}

	*first = 7
	if value != 100 || *second != 42 {
		t.Fatal("changing one result changed the input or another result")
	}
}

func TestNewReferenceValue(t *testing.T) {
	value := []int{1, 2, 3}
	got := ptr.New(value)
	if got == nil {
		t.Fatal("New returned nil")
	}
	if len(*got) != len(value) || &(*got)[0] != &value[0] {
		t.Fatal("New must preserve the slice's backing array")
	}

	(*got)[0] = 99
	if value[0] != 99 {
		t.Fatal("New must copy the slice value without deep copying its elements")
	}
}

func TestNewNilValues(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		checkNew[*int](t, nil)
	})
	t.Run("interface", func(t *testing.T) {
		checkNew[any](t, nil)
	})
	t.Run("slice", func(t *testing.T) {
		got := ptr.New[[]int](nil)
		if got == nil || *got != nil {
			t.Fatal("want a non-nil pointer to a nil slice")
		}
	})
	t.Run("map", func(t *testing.T) {
		got := ptr.New[map[string]int](nil)
		if got == nil || *got != nil {
			t.Fatal("want a non-nil pointer to a nil map")
		}
	})
}
