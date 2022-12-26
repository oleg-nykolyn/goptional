# goptional

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

The definitive implementation of the `Optional` type in Go 🚀

## Features

- A recognizablle API that is *heavily* inspired by the Java & Rust implementations
- Compatibility with all *value* and *reference* types
- *Chainable* generic methods and functions

## Installation

> ❗️ *goptional* requires **go 1.19**

```bash
go get -u github.com/nykolynoleg/goptional
```

## Usage

> 👨🏻‍💻 Refer to the [documentation](https://pkg.go.dev/github.com/nykolynoleg/goptional) for a complete reference.

```go
// Import goptional into your code and refer to it as `goptional`.
import "github.com/nykolynoleg/goptional"
```

### Creation

```go
// Create an Optional of type int that holds 123.
// All value and reference types are supported.
opt := goptional.Of(123)
```

```go
// Create an empty Optional of type string.
opt := goptional.Empty[string]()
```

```go
// 'Of' returns an empty Optional if its argument is nil.
opt := goptional.Of[[]string](nil)
```

```go
// Note that if the argument is the zero value of a value type,
// such as "", false, 0 then a non-empty Optional is returned instead.
opt := goptional.Of("")
```

### Presence Checks

```go
opt := goptional.Of(123)

// Is true, as opt holds 123.
if opt.IsPresent() {
    // ...
}

// Is false.
if opt.IsEmpty() {
    // ...
}
```

### Value Retrieval

`Get`

```go
opt := goptional.Of(123)

// Retrieve the value held by opt.
// Panic otherwise.
v := opt.Get()
```

`OrElse`

```go
opt := goptional.Empty[string]()

// Provide a default value if opt is empty.
v := opt.OrElse("default")
```

`OrElseGet`

```go
opt := goptional.Empty[string]()

// Provide a default through a supplier if opt is empty.
v := opt.OrElseGet(func() string {
    return "default"
})
```

`OrElsePanicWithErr`

```go
opt := goptional.Empty[string]()

// Panic with an error provided by the given supplier if opt is empty.
v := opt.OrElsePanicWithErr(func() error {
    return errors.New("woops")
})
```

### Filtering

```go
opt := goptional.Of(123)

// Apply a predicate to the value of opt, if there is any.
opt = opt.Filter(func(v int) bool { return v > 100 })
// Returns an empty Optional, as 123 is not even.
opt = opt.Filter(func(v int) bool { return v%2 == 0 })

v := 0

// Is false.
if opt.IsPresent() {
    v = opt.Get()
}
```

```go
// The example above can be rewritten in a fluent style.
v := goptional.Of(123).
    Filter(func(v int) bool { return v > 100 }).
    Filter(func(v int) bool { return v%2 == 0 }).
    OrElse(0)
```

### Mapping

`Map`

```go
opt := goptional.Of(123)

// Apply the given transformation to the value of opt, if there is any,
// and return a new Optional of the target type.
strOpt := goptional.Map(opt, func(v int) string {
    return fmt.Sprintf("%v_mapped", v)
})

// v is "123_mapped"
v := strOpt.OrElse("")
```

`MapOr`

```go
opt := goptional.Empty[int]()

// Similar to Map, but returns an Optional holding
// the given default value if opt is empty.
strOpt := goptional.MapOr(opt, func(v int) string {
    return fmt.Sprintf("%v_mapped", v)
}, "default")

// v is "default"
v := strOpt.OrElse("")
```

`MapOrElse`

```go
opt := goptional.Empty[int]()

// Similar to Map, but returns an Optional holding
// a default value provided by the given supplier if opt is empty.
strOpt := goptional.MapOrElse(opt, func(v int) string {
    return fmt.Sprintf("%v_mapped", v)
}, func() string {
    return "default"
})

// v is "default"
v := strOpt.OrElse("")
```

`FlatMap`

```go
opt := Of(123)

// FlatMap is similar to Map, but the given supplier returns an Optional instead.
// If you are familiar with Monads, think of it as AndThen.
strOpt := FlatMap(opt, func(v int) *Optional[string] {
    return Of(fmt.Sprintf("%v_mapped", v))
})

// v is "123_mapped"
v := strOpt.OrElse("")
```

```go
opt := Empty[int]()

// Returns a new empty Optional of the target type, as opt is empty.
strOpt := FlatMap(opt, func(v int) *Optional[string] {
    return Of(fmt.Sprintf("%v_mapped", v))
})

// v is ""
v := strOpt.OrElse("")
```

### Peeking

`IfPresent`

```go
opt := goptional.Of(123)

// Execute the given action on the value of opt, if there is any.
// Do nothing otherwise.
opt.IfPresent(func(v int) {
    fmt.Println(v) // Prints '123'
})
```

`IfPresentOrElse`

```go
opt := goptional.Empty[int]()

// Similar to IfPresent, but executes a fallback action if opt is empty.
opt.IfPresentOrElse(func(v int) {
    // ...
}, func() {
    // This will be executed, as 'opt' is empty.
})
```

### Boolean Operators

> Think of an *empty* `Optional` as `false` and `true` otherwise.  
> The suppliers used in the examples below are lazily-evaluated.  
> If a boolean expression can be *short-circuited*, the supplier is ignored.

`And`

```go
opt := goptional.Empty[int]()

// AND between opt & the supplied Optional.
opt = opt.And(func() *Optional[int] {
    return goptional.Of(123)
})

// v is 0
v := opt.OrElse(0)
```

`Or`

```go
opt := goptional.Empty[int]()

// OR between opt & the supplied Optional.
opt = opt.Or(func() *Optional[int] {
    return goptional.Of(123)
})

// v is 123
v := opt.OrElse(0)
```

## Testing

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
