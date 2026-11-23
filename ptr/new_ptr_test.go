package ptr_test

import (
	"testing"

	"github.com/cohesian-tools/typy/ptr"
)

func TestNewPtrValues(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		checkNewPtr(t, 42)
		checkNewPtr(t, 0)
		checkNewPtr(t, -42)
	})
	t.Run("string", func(t *testing.T) {
		checkNewPtr(t, "hello")
		checkNewPtr(t, "")
	})
	t.Run("bool", func(t *testing.T) {
		checkNewPtr(t, true)
		checkNewPtr(t, false)
	})
	t.Run("named type", func(t *testing.T) {
		type UserID int
		checkNewPtr(t, UserID(42))
	})
	t.Run("struct", func(t *testing.T) {
		type User struct{ Name string }
		checkNewPtr(t, User{Name: "Lucas"})
		checkNewPtr(t, User{})
	})
}

func checkNewPtr[T comparable](t *testing.T, value T) {
	t.Helper()
	var got *T = ptr.NewPtr(value)
	if got == nil {
		t.Fatal("NewPtr returned nil")
	}
	if *got != value {
		t.Fatalf("*NewPtr(%#v) = %#v; want %#v", value, *got, value)
	}
}

func TestNewPtrIndependentStorage(t *testing.T) {
	value := 42
	first := ptr.NewPtr(value)
	second := ptr.NewPtr(value)
	if first == &value || second == &value || first == second {
		t.Fatal("NewPtr must allocate independent storage for each call")
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

func TestNewPtrReferenceValue(t *testing.T) {
	value := []int{1, 2, 3}
	got := ptr.NewPtr(value)
	if got == nil {
		t.Fatal("NewPtr returned nil")
	}
	if len(*got) != len(value) || &(*got)[0] != &value[0] {
		t.Fatal("NewPtr must preserve the slice's backing array")
	}

	(*got)[0] = 99
	if value[0] != 99 {
		t.Fatal("NewPtr must copy the slice value without deep copying its elements")
	}
}

func TestNewPtrNilValues(t *testing.T) {
	t.Run("pointer", func(t *testing.T) {
		checkNewPtr[*int](t, nil)
	})
	t.Run("interface", func(t *testing.T) {
		checkNewPtr[any](t, nil)
	})
	t.Run("slice", func(t *testing.T) {
		got := ptr.NewPtr[[]int](nil)
		if got == nil || *got != nil {
			t.Fatal("want a non-nil pointer to a nil slice")
		}
	})
	t.Run("map", func(t *testing.T) {
		got := ptr.NewPtr[map[string]int](nil)
		if got == nil || *got != nil {
			t.Fatal("want a non-nil pointer to a nil map")
		}
	})
}
