# Optional values

[Back to the toolbox](../readme.md)

`opt.Option[T]` represents either a present value (`Some`) or an absent value
(`None`). Presence is independent of the zero value of `T`.

```go
import "github.com/cohesian-tools/typy/opt"
```

## Construct, inspect, and unwrap

```go
present := opt.Some(0)
absent := opt.None[int]()
var zero opt.Option[int] // also None

fmt.Println(present.IsSome()) // true
fmt.Println(absent.IsSome())  // false
fmt.Println(zero.IsSome())    // false

if present.IsSome() {
    fmt.Println(present.Unwrap()) // 0
}
```

Import `fmt` for printing. `Some` accepts any type, including slices, maps,
pointers, and structs. `Some[*int](nil)` is present even though its inner pointer
is nil.

`Unwrap()` on None panics with the string `"unable to unwrap None"`. The exported
`ErrUnwrapNone` variable is retained from the standalone package, but it is not
the panic value and is not returned by any current method.

`IsSome` and `Unwrap` have pointer receivers. Store constructor results in a
variable before calling them; `opt.Some(42).Unwrap()` is not a valid call.

## Copying and mutation

`Some` and `Unwrap` make shallow copies of `T`. Reference values keep their
underlying references. Copying an Option also copies its internal pointer,
so decoding into one copy can affect another copy's stored value. Options
should not be treated as deeply immutable or safe for concurrent mutation.
Synchronize shared mutation, including JSON and BSON decoding.

## JSON

```go
type Payload struct {
    Name opt.Option[string] `json:"name"`
}

var payload Payload
err := json.Unmarshal([]byte(`{"name":"Lucas"}`), &payload)
if err != nil {
    return err
}
fmt.Println(payload.Name.Unwrap()) // Lucas

data, err := json.Marshal(payload)
// data is {"name":"Lucas"} when err == nil
```

This snippet belongs inside a function returning `error`; import `encoding/json`
and `fmt`. Marshal accepts an Option value or pointer, while unmarshal needs a
pointer to update the value.

| Value or input | Behavior |
| --- | --- |
| Marshal Some | Encode its inner value |
| Marshal None | Encode `null` |
| Decode a non-null value | Set Some, if decoding succeeds |
| Decode `null` | Set None, including when replacing Some |
| A missing struct field | Leave that field unchanged |

For a freshly zero-initialized struct, a missing field and an explicit `null`
both result in None. When reusing a struct, a missing field preserves an existing
Some. None does not itself request omission of the containing struct field;
the codec emits `null`. Encoding a present nil pointer also emits `null`, so
presence of that value is lost on a round trip.

The inner type controls which JSON values are accepted. Codec errors propagate
to the caller; decoding is not guaranteed to leave the option unchanged on error.

## BSON

The migrated encoders retain `go.mongodb.org/mongo-driver` v1.12.1, including
its `bson` and `bson/bsontype` packages.

```go
type Document struct {
    Name opt.Option[string] `bson:"name"`
}

original := Document{Name: opt.Some("Lucas")}
data, err := bson.Marshal(original)
if err != nil {
    return err
}
var decoded Document
if err := bson.Unmarshal(data, &decoded); err != nil {
    return err
}
fmt.Println(decoded.Name.Unwrap()) // Lucas
```

This snippet also belongs inside an error-returning function. Import
`go.mongodb.org/mongo-driver/bson` and `fmt` in addition to `opt`.
Some encodes its inner value, None encodes BSON null, and decoding null sets
None. No running MongoDB server is needed to use or test these codecs.

## Option API

```go
var ErrUnwrapNone = errors.New("cannot unwrap None")

type Option[T any] struct { /* unexported fields */ }

func Some[T any](value T) Option[T]
func None[T any]() Option[T]
func (o *Option[T]) IsSome() bool
func (o *Option[T]) Unwrap() T
func (o *Option[T]) UnmarshalJSON(data []byte) error
func (o Option[T]) MarshalJSON() ([]byte, error)
func (o *Option[T]) UnmarshalBSONValue(t bsontype.Type, data []byte) error
func (o Option[T]) MarshalBSONValue() (bsontype.Type, []byte, error)
```

## Nullable values

The original nullable package is preserved under `opt/nullable`:

```go
import "github.com/cohesian-tools/typy/opt/nullable"
```

`nullable.Null[T]` stores a value and a validity flag. Its zero value is null.

```go
amount := nullable.Some(0)
missing := nullable.None[int]()

fmt.Println(amount.IsNull())  // false: zero is a valid value
fmt.Println(amount.Unwrap())  // 0
fmt.Println(missing.IsNull()) // true
```

`Unwrap` on a null value panics with `"unable to unwrap Null"`. Its methods have
value receivers. Like Option, values are copied shallowly.

```go
type Null[T any] struct { /* unexported fields */ }

func Some[T any](value T) Null[T]
func None[T any]() Null[T]
func (n Null[T]) IsNull() bool
func (n Null[T]) Unwrap() T
```

| Type | Representation | Inspection | Encoding support |
| --- | --- | --- | --- |
| `opt.Option[T]` | Internal `*T` | `IsSome()` | JSON and BSON |
| `nullable.Null[T]` | Inline `T` plus validity flag | `IsNull()` | No custom codecs |

`nullable.Null` does not implement JSON/BSON codecs or SQL `Scanner`/`Valuer`.
It is an in-memory type; it is not automatically a replacement for `sql.Null*`.
Nesting it inside Option can model states in memory, but the current codecs
do not implement a missing/null/value wire format for that combination.

## Checks

From the module root:

```sh
go test ./opt/... -cover
```
