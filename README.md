# goptional

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

🦄 A *comprehensive* implementation of the `Optional` type in Go

## Features

- A recognizablle API that is *heavily* inspired by the [Java](https://github.com/AdoptOpenJDK/openjdk-jdk11/blob/master/src/java.base/share/classes/java/util/Optional.java) & [Rust](https://doc.rust-lang.org/std/option/enum.Option.html) implementations
- Compatible with all *value* and *reference* types
- Chainable and expressive [generic methods and functions](https://go.dev/doc/tutorial/generics)
- Minimal overhead - `Optional` is just a singleton slice
- Null safety *by design*

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

`Of`

```go
// Create an Optional of type int that holds 123.
opt := goptional.Of(123)

// If the argument to Of is either nil or invalid, an empty Optional is returned instead.
opt2 := goptional.Of[[]string](nil)

// Is true.
if opt2.IsEmpty() {
    // ...
}
```

`Empty`
> ✍🏼 Note that an empty `Optional` is effectively a `nil` pointer.  
> Regardless, you can **safely** call any methods on it without expecting a nullptr `panic`.

```go
// Create an empty Optional of type string.
opt := goptional.Empty[string]()

// Is true.
if opt == nil {

}

// Will not panic.
if opt.IsPresent() {
    // ...
}
```

### Presence

```go
opt := goptional.Of(123)

// Check if opt holds a value.
// In this specific instance, it returns true.
if opt.IsPresent() {
    // ...
}

// Check if opt has no value.
if opt.IsEmpty() {
    // ...
}
```

### Equality

```go
opt := goptional.Of(123)
opt2 := goptional.Of(321)

// Compare opt & opt2 for equality.
// It returns true if both contain the same value, or if both are empty.
// Otherwise, it returns false.
if opt.Equals(opt2) {
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

`OrZero`

```go
opt := goptional.Empty[string]()

// Retrieve the value held by opt or 
// the zero value of its type otherwise.
//
// Some zero values:
// string -> ""
// bool -> false
// int -> 0
// ptr -> nil
v := opt.OrZero()

// Is true.
if v == "" {
    // ...
}
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
opt := goptional.Of(123)

// FlatMap is similar to Map, but the given supplier returns an Optional instead.
// If you are familiar with Monads, think of it as AndThen.
strOpt := FlatMap(opt, func(v int) Optional[string] {
    return goptional.Of(fmt.Sprintf("%v_mapped", v))
})

// v is "123_mapped"
v := strOpt.OrElse("")
```

```go
opt := goptional.Empty[int]()

// Returns a new empty Optional of the target type, as opt is empty.
strOpt := FlatMap(opt, func(v int) Optional[string] {
    return goptional.Of(fmt.Sprintf("%v_mapped", v))
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
    // This block will execute, as 'opt' is empty.
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
opt = opt.And(func() Optional[int] {
    return goptional.Of(123)
})

// v is 0
v := opt.OrElse(0)
```

`Or`

```go
opt := goptional.Empty[int]()

// OR between opt & the supplied Optional.
opt = opt.Or(func() Optional[int] {
    return goptional.Of(123)
})

// v is 123
v := opt.OrElse(0)
```

`Xor`

```go
opt := goptional.Empty[int]()

// XOR between opt & the given Optional.
opt = opt.Xor(goptional.Of(321))

// v is 321
v := opt.OrElse(0)
```

### Mutations

`Take`

```go
opt1 := goptional.Of(123)

// Take the value from opt1, if any,
// and transfer it to opt2 by leaving opt1 empty.
opt2 := opt1.Take()

// Is true.
if opt1.IsEmpty() {
    // ...
}

// v is 123
v := opt2.Get()
```

`Replace`

```go
opt1 := goptional.Of(123)

// opt1 -> 789, opt2 -> 123
opt2 := opt1.Replace(789)
```

```go
opt1 := goptional.Empty[int]()

// opt1 -> 789, opt2 -> Empty
opt2 := opt1.Replace(789)
```

### JSON

`MarshalJSON`

```go
opt := goptional.Of(123)

// Get the JSON representation of opt.
// Marshal returns []byte("null") if an Optional is empty.
jsonBytes, err := opt.MarshalJSON()
```

`UnmarshalJSON`

```go
opt := goptional.Empty[int]()
intAsJSON := "123"

// Populate opt with the given JSON.
err := opt.UnmarshalJSON([]byte(intAsJSON))
if err == nil {
    // v is 123
    v := opt.Get()
}
```

### String Representation

`Optional` implements the `Stringer` interface and relies on [spew](https://github.com/davecgh/go-spew).

##  FAQ

1. **Why are `Map`, `MapOr`, etc. implemented as functions and not methods?**  
As of now, Go does **not** support method-level type parameters. This might change in the future.

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
