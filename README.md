# goptional

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

A generics-based implementation of the `Optional` type in Go 🚀

## Features

- A recognizablle API that is *heavily* inspired by the Java & Rust implementations
- Compatibility with all *value* and *reference* types
- *Chainable* generic methods and functions

## Installation

> ❗️ *goptional* requires **go 1.19**

```bash
go get -u github.com/nykolynoleg/goptional
```

## Examples

> 👨🏻‍💻 Refer to the [documentation](https://pkg.go.dev/github.com/nykolynoleg/goptional) for a complete reference.

```go
// Import goptional into your code and refer to it as `goptional`.
import "github.com/nykolynoleg/goptional"
```

### Creation

```go
// Creates an Optional that holds an int value.
intOpt := goptional.Of(123)

// Creates an empty Optional of type string.
strOpt := goptional.Empty[string]()

// 'Of' returns an empty Optional if its argument is nil.
strSliceOpt := goptional.Of[[]string](nil)

// Note that if Of's argument is the zero value of a value type,
// such as ("", false, 0) then a non-empty Optional is returned.
emptyStrOpt := goptional.Of("")
```

### Inspection

```go
strOpt := goptional.Of("hello goptional")
intOpt := goptional.Empty[int]()

// Returns true if Optional holds a value.
if strOpt.IsPresent() {
    // ...
}

// Returns true if Optional is empty.
if intOpt.IsEmpty() {
    // ...
}

// Method chaining is supported.
if goptional.Of(123).IsPresent() {
    // ...
}
```

### Retrieval

```go
// TODO
```

### Filtering

```go
// TODO
```

### Mapping

```go
// TODO
```

### Peeking

```go
// TODO
```

## Testing

Without coverage:

```bash
go test ./... -v
```

With coverage & HTML output:

```bash
go test ./... -v -race -covermode=atomic -coverprofile=coverage.txt
go tool cover -html=coverage.txt
```

## Contributing

Any kind of support is more than welcome 🤝  
Refer to the [contribution guide](CONTRIBUTING.md) and the [code of conduct](CODE_OF_CONDUCT.md) for details.

## Creator

[Oleg Nykolyn](https://linktr.ee/lgnk)

---

Released under the [MIT license](LICENSE.txt).

[doc-img]: https://pkg.go.dev/badge/github.com/nykolynoleg/goptional
[doc]: https://pkg.go.dev/github.com/nykolynoleg/goptional
[ci-img]: https://github.com/nykolynoleg/goptional/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/nykolynoleg/goptional/actions/workflows/go.yml
[cov-img]: https://codecov.io/gh/nykolynoleg/goptional/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/nykolynoleg/goptional
