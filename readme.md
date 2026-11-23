# typy

My general-purpose Go toolbox: small packages for pointers, layered errors,
and optional values. Each tool has its own import path under one module.

Requires **Go 1.25 or newer**.

## Tools and documentation

| Package | Purpose | Documentation |
| --- | --- | --- |
| `ptr` | Create a pointer initialized from a value | [Pointer helpers](docs/ptr.md) |
| `errlay` | Wrap errors with layer/code metadata and match error chains | [Layered errors](docs/errlay.md) |
| `opt` | Represent optional values, with JSON and BSON encoding | [Optional values](docs/opt.md) |
| `opt/nullable` | Represent an explicit null state separately from a zero value | [Nullable values](docs/opt.md#nullable-values) |

## Installation

```sh
go get github.com/cohesian-tools/typy@v0.1.0
```

## Usage

Import the tools you need:

```go
import (
    "github.com/cohesian-tools/typy/ptr"
    "github.com/cohesian-tools/typy/errlay"
    "github.com/cohesian-tools/typy/opt"
)
```

For example, pointer construction also makes error filters convenient:

```go
err := errlay.New("database", 10, "query failed", nil)
target := errlay.Target(ptr.NewPtr("database"), ptr.NewPtr(uint8(10)))
matched := errors.Is(err, target) // true; import "errors"

name := opt.Some("Lucas")
if name.IsSome() {
    fmt.Println(name.Unwrap()) // Lucas; import "fmt"
}
```

For local development, a consumer module can use the absolute checkout path:

```sh
go mod edit -replace=github.com/cohesian-tools/typy=/absolute/path/to/typy
go mod tidy
```

## Layout

| Path | Contents |
| --- | --- |
| `go.mod`, `go.sum` | Shared module and dependency versions |
| `ptr/` | Pointer implementation, tests, and examples |
| `errlay/` | Layered error implementation, tests, and examples |
| `opt/` | Optional implementation, encoders, tests, and examples |
| `opt/nullable/` | Nullable implementation, tests, and examples |
| `docs/` | Per-tool guides and API references |

The module root is the documentation entry point. Import a tool's subpackage
rather than the module root.

## Migrating existing imports

| Previous import | Toolbox import | Package name |
| --- | --- | --- |
| `github.com/cohesian-tools/typy` (`NewPtr`) | `github.com/cohesian-tools/typy/ptr` | `ptr` |
| `github.com/cohesian-tools/errlay` | `github.com/cohesian-tools/typy/errlay` | `errlay` |
| `github.com/cohesian-tools/optional` | `github.com/cohesian-tools/typy/opt` | `opt` |
| `github.com/cohesian-tools/optional/nullable` | `github.com/cohesian-tools/typy/opt/nullable` | `nullable` |

Change `typy.NewPtr(...)` to `ptr.NewPtr(...)` and `optional.Some(...)` to
`opt.Some(...)`, and likewise for other optional identifiers. An import alias
such as `optional "github.com/cohesian-tools/typy/opt"` can keep existing call
sites intact while migrating. Run `go mod tidy` in each consumer afterward.

The implementations and existing tests were brought over from the local
`errlay` and `optional` repositories. Those repositories remain independent
and have not been deleted, redirected, or deprecated. Their types and the
toolbox types have different package identities: migrate related imports
together rather than mixing old and new types or error targets.

## Dependencies and checks

`ptr`, `errlay`, and `opt/nullable` use only the standard library.
`opt` retains its BSON codecs and the original MongoDB Go driver dependency
(`go.mongodb.org/mongo-driver` v1.12.1). Tests retain `testify` v1.8.3.
All dependencies live in the root module; the subpackages are not modules.

```sh
go test ./... -cover
go vet ./...
```

Tests sit beside their implementations. Executable examples also exercise the
public import paths and common combinations of tools.

## License

MIT. See [LICENSE](LICENSE). The original 2023 optional and 2026 errlay
copyright notices are retained.
