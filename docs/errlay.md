# Layered errors

[Back to the toolbox](../readme.md)

`errlay` attaches a logical layer, a numeric code, and a message to an error.
The underlying cause stays available through the standard error chain.

```go
import "github.com/cohesian-tools/typy/errlay"
```

## Create and inspect errors

```go
cause := errors.New("connection refused")
err := errlay.New("database", 10, "query failed", cause)

fmt.Println(err)
// |layer database - 10| query failed connection refused

fmt.Println(errors.Is(err, cause)) // true

var detail *errlay.Error
if errors.As(err, &detail) {
    fmt.Println(detail.Layer(), detail.Code(), detail.Message())
    // database 10 query failed
}
```

Import `errors` and `fmt` from the standard library. A nil cause creates a leaf
error. Codes are `uint8` values, from 0 through 255; layers are arbitrary strings.

`Message()` returns only the error's own message. `Error()` formats the metadata
and appends the underlying cause, separated by a space. `Unwrap()` returns the
cause, or nil for a leaf.

## Match by metadata

```go
layer := "database"
code := uint8(10)

byLayer := errlay.Target(&layer, nil)
byCode := errlay.Target(nil, &code)
exact := errlay.Target(&layer, &code)

fmt.Println(errors.Is(err, byLayer)) // true
fmt.Println(errors.Is(err, byCode))  // true
fmt.Println(errors.Is(err, exact))   // true
```

| Target | Match condition |
| --- | --- |
| Layer only | Same layer |
| Code only | Same code |
| Both fields | Both match the same error in the chain |
| `Target(nil, nil)` | Any `*errlay.Error` in the chain |

Nil fields are wildcards. Empty strings and code 0 are explicit values when
their pointers are supplied. Layer comparisons are case-sensitive, and messages
do not participate in matching.

`Target` returns one error value and accepts an empty filter.

Targets retain the supplied pointers, so changing a pointed-to variable changes
the filter. Keep filter values unchanged while matching; concurrent mutation
requires synchronization. The `ptr` package is convenient for fixed filters:

```go
target := errlay.Target(ptr.New("database"), ptr.New(uint8(10)))
```

## Error-chain behavior

`errors.Is` follows wrapped causes and applies the filter to each error. A
combined filter cannot match the layer of one error and the code of another:

```go
err := errlay.New("service", 10, "request failed",
    errlay.New("database", 20, "query failed", nil),
)
layer, code := "database", uint8(10)
fmt.Println(errors.Is(err, errlay.Target(&layer, &code))) // false
```

Standard wrappers such as `fmt.Errorf("context: %w", err)` preserve traversal.
`errors.As` retrieves the first `*errlay.Error`, starting from the outermost
error. A filter constructed by `Target` is intended for `errors.Is`;
`errors.As` extracts values by type and does not use these field filters.

Two independently created errors with identical metadata do not match each
other by metadata. Use `Target` for filtering. Standard error identity still
works when the target is the same error instance or an underlying sentinel.

## API

```go
func New(layer string, code uint8, msg string, err error) error
func Target(layer *string, code *uint8) error

type Error struct { /* unexported fields */ }

func (e *Error) Error() string
func (e *Error) Layer() string
func (e *Error) Code() uint8
func (e *Error) Message() string
func (e *Error) Unwrap() error
func (e *Error) Is(target error) bool
```

Calling `e.Is(target)` checks only that error's fields. Use `errors.Is` to search
the entire chain. Supply ordinary non-nil error instances and targets; the
implementation does not add special handling for typed nil `*Error` values.

## Checks

From the module root:

```sh
go test ./errlay -cover
```
