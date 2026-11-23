# Pointer helpers

[Back to the toolbox](../readme.md)

```go
import "github.com/cohesian-tools/typy/ptr"
```

## API

```go
func NewPtr[T any](t T) *T
```

`NewPtr` returns a non-nil pointer initialized with a copy of its argument.
Type inference usually determines `T`:

```go
name := ptr.NewPtr("Lucas") // *string
count := ptr.NewPtr(0)      // *int, pointing to zero
enabled := ptr.NewPtr(true) // *bool
```

Named types are preserved. For example, passing a `UserID` returns a `*UserID`,
rather than a pointer to its underlying type.

## Copy behavior

For a scalar or a struct containing only value fields, changing the result
does not change the input or another call's result:

```go
original := 42
copy := ptr.NewPtr(original)
*copy = 7
// original is still 42.
```

This is a shallow copy. Slices retain their backing array, maps retain their
entries, and pointers retain their pointees. A new pointer to a slice does not
clone its elements. The same applies to reference fields inside structs.

Do not rely on distinct pointer addresses for zero-sized types: Go permits
pointers to distinct zero-sized variables to compare equal.

## Nil values

Specify a type when the argument is untyped `nil`:

```go
p := ptr.NewPtr[*int](nil) // **int: p != nil, *p == nil
s := ptr.NewPtr[[]int](nil)
// s != nil, *s == nil
```

The returned pointer exists even when the value stored through it is nil.

## Using pointers as filters

```go
target := errlay.Target(
    ptr.NewPtr("database"),
    ptr.NewPtr(uint8(10)),
)
```

Import `github.com/cohesian-tools/typy/errlay` for this example. Explicitly
convert the numeric code to `uint8`, since an untyped integer would otherwise
infer `int`.

## Checks

From the module root:

```sh
go test ./ptr -cover
```
