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

// 💡 If the argument to Of is either nil or invalid, return an empty Optional instead.
opt2 := goptional.Of[[]string](nil)
```

`Empty`
> ✍🏼 Note that an empty `Optional` is effectively a `nil` pointer.  
> Regardless, you can **safely** call any methods on it without expecting a nullptr `panic`.

```go
// Create an empty Optional of type string.
opt := goptional.Empty[string]()

fmt.Println(opt)        // Optional.empty
fmt.Println(opt == nil) // true

// This does not panic.
fmt.Println(opt.IsPresent()) // false
```

### Presence Checks

```go
opt := goptional.Of(123)

// Check if opt holds a value.
fmt.Println(opt.IsPresent()) // true

// Check if opt is empty.
fmt.Println(opt.IsEmpty()) // false
```

### Equality

```go
opt := goptional.Of(123)
opt2 := goptional.Of(321)

// Compare opt & opt2 for equality.
// Return true if both contain the same value, or if both are empty.
// Otherwise, return false.
fmt.Println(opt.Equals(opt2)) // false
```

### Value Retrieval

`Get`

```go
opt := goptional.Of(123)

// Retrieve the value held by opt, if any, or panic otherwise.
fmt.Println((opt.Get())) // 123

opt2 := goptional.Empty[int]()

fmt.Println(opt2.Get()) // panics
```

`OrElse`

```go
opt := goptional.Empty[string]()

// Provide a default value if opt is empty.
fmt.Println(opt.OrElse("lfg")) // lfg
```

`OrElseGet`

```go
opt := goptional.Empty[string]()

// Provide a default through a supplier if opt is empty.
v := opt.OrElseGet(func() string {
    return "gm"
})

fmt.Println(v) // gm
```

`OrZero`

```go
opt := goptional.Empty[string]()

// Retrieve the value held by opt, if any, or 
// the zero value of its type otherwise.
fmt.Println(opt.OrZero()) // ""
```

`OrElsePanicWithErr`

```go
opt := goptional.Empty[string]()

// Panic with an error provided by the given supplier if opt is empty.
_ = opt.OrElsePanicWithErr(func() error {
    return errors.New("woops")
}) // panics
```

### Filtering

```go
opt := goptional.Of(123)

// Apply a predicate to the value of opt, if any.
opt = opt.Filter(func(v *int) bool { return *v > 100 })
// Return an empty Optional, as 123 is not even.
opt = opt.Filter(func(v *int) bool { return *v%2 == 0 })

fmt.Println(opt.IsPresent()) // false
fmt.Println(opt.Get())       // panics
```

```go
// The example above can be rewritten in a fluent style.
v := goptional.Of(123).
    Filter(func(v *int) bool { return *v > 100 }).
    Filter(func(v *int) bool { return *v%2 == 0 }).
    OrElse(0)

fmt.Println(v) // 0
```

### Mapping

`Map`

```go
opt := goptional.Of(123)

// Apply the given transformation to the value of opt, if any,
// and return a new Optional of the target type.
strOpt := goptional.Map(opt, func(v *int) string {
    return fmt.Sprintf("%v_mapped", *v)
})

fmt.Println(strOpt.OrZero()) // 123_mapped
```

`MapOr`

```go
opt := goptional.Empty[int]()

// Similar to Map, but returns an Optional holding
// the given default value if opt is empty.
strOpt := goptional.MapOr(opt, func(v *int) string {
    return fmt.Sprintf("%v_mapped", *v)
}, "default")

fmt.Println(strOpt.OrZero()) // default
```

`MapOrElse`

```go
opt := goptional.Empty[int]()

// Similar to Map, but returns an Optional holding
// a default value provided by the given supplier if opt is empty.
strOpt := goptional.MapOrElse(opt, func(v *int) string {
    return fmt.Sprintf("%v_mapped", *v)
}, func() string {
    return "default"
})

fmt.Println(strOpt.OrZero()) // default
```

`FlatMap`

```go
opt := goptional.Of(123)

// FlatMap is similar to Map, but the given supplier returns an Optional instead.
// If you are familiar with Monads, think of it as AndThen.
strOpt := goptional.FlatMap(opt, func(v *int) goptional.Optional[string] {
    return goptional.Of(fmt.Sprintf("%v_mapped", *v))
})

fmt.Println(strOpt.OrZero()) // 123_mapped
```

```go
opt := goptional.Empty[int]()

// Return a new empty Optional of the target type, as opt is empty.
strOpt := goptional.FlatMap(opt, func(v *int) goptional.Optional[string] {
    return goptional.Of(fmt.Sprintf("%v_mapped", v))
})

fmt.Println(strOpt.OrZero()) // ""
```

`Flatten`

```go
opt := goptional.Of(goptional.Of(123))

// Flatten opt by returning the wrapped Optional, if any.
// Transform Optional[Optional[T]] into Optional[T].
fOpt := goptional.Flatten(opt)

fmt.Println(fOpt.Get()) // 123
```

### Peeking

`IfPresent`

```go
opt := goptional.Of(123)

// Execute the given action on the value of opt, if any.
// Do nothing otherwise.
opt.IfPresent(func(v *int) {
    fmt.Println(v) // 123
})
```

`IfPresentOrElse`

```go
opt := goptional.Empty[int]()

// Similar to IfPresent, but executes a fallback action if opt is empty.
opt.IfPresentOrElse(func(v *int) {
    // ...
}, func() {
    // This block will execute, as opt is empty.
})
```

### Boolean Operators

> Think of an *empty* `Optional` as `false` and `true` otherwise.  
> The suppliers used in the examples below are lazily-evaluated.  
> If a boolean expression can be *short-circuited*, the supplier is *ignored*.

`And`

```go
opt := goptional.Empty[int]()

// AND between opt & the supplied Optional.
opt = opt.And(func() goptional.Optional[int] {
    return goptional.Of(123)
})

fmt.Println(opt.OrZero()) // 0
```

`Or`

```go
opt := goptional.Empty[int]()

// OR between opt & the supplied Optional.
opt = opt.Or(func() goptional.Optional[int] {
    return goptional.Of(123)
})

fmt.Println(opt.OrZero()) // 123
```

`Xor`

```go
opt := goptional.Empty[int]()

// XOR between opt & the given Optional.
opt = opt.Xor(goptional.Of(321))

fmt.Println(opt.OrZero()) // 321
```

### Mutations

`Take`

```go
opt1 := goptional.Of(123)

// Take the value from opt1, if any,
// and transfer it to opt2 by leaving opt1 empty.
opt2 := opt1.Take()

fmt.Println(opt1.IsEmpty()) // true
fmt.Println(opt2.Get())     // 123
```

`Replace`

```go
opt1 := goptional.Of(123)

// Transfer the value of opt1, if any, to opt2 and replace it with 789.
opt2 := opt1.Replace(789)

fmt.Println(opt1.Get()) // 789
fmt.Println(opt2.Get()) // 123
```

```go
opt1 := goptional.Empty[int]()
opt2 := opt1.Replace(789)

fmt.Println(opt1.Get()) // 789
fmt.Println(opt2.IsEmpty()) // true
```

### JSON

`MarshalJSON`

```go
opt := goptional.Of(123)

// Get the JSON representation of opt.
// It returns []byte("null") if called on an empty Optional.
jsonBytes, err := opt.MarshalJSON()

fmt.Println(err == nil)        // true
fmt.Println(string(jsonBytes)) // 123

// Make opt empty.
opt = goptional.Empty[int]()
jsonBytes, err = opt.MarshalJSON()

fmt.Println(err == nil)        // true
fmt.Println(string(jsonBytes)) // null
```

`UnmarshalJSON`

```go
opt := goptional.Empty[int]()
numAsJSON := "123"

// Populate opt with the given JSON.
err := opt.UnmarshalJSON([]byte(numAsJSON))

fmt.Println(err == nil) // true
fmt.Println(opt.Get())  // 123
```

### Zipping

`Zip`

```go
opt1 := goptional.Of(123)
opt2 := goptional.Of("gm")

// Zip opt1 & opt2 and return a non-empty Optional Pair.
optPair := goptional.Zip(opt1, opt2)

fmt.Println(optPair.IsPresent()) // true

pair := optPair.Get()
fmt.Println(pair.First)  // 123
fmt.Println(pair.Second) // gm

// Return an empty Optional Pair if at least one of the arguments to Zip is an empty Optional.
optPair = goptional.Zip(opt1, goptional.Empty[string]())

fmt.Println(optPair.IsEmpty()) // true
```

`ZipWith`

```go
opt1 := goptional.Of(123)
opt2 := goptional.Of("gm")
mapper := func(x *int, y *string) string {
    return fmt.Sprintf("%v_%v", *x, *y)
}
// Zip opt1 & opt2 with the given mapper and return a non-empty Optional of the target type.
opt3 := goptional.ZipWith(opt1, opt2, mapper)

fmt.Println(opt3.IsPresent()) // true
fmt.Println(opt3.Get())       // 123_gm

// Return an empty Optional as one of the arguments to ZipWith is an empty Optional.
opt3 = goptional.ZipWith(opt1, goptional.Empty[string](), mapper)

fmt.Println(opt3.IsEmpty()) // true
```

`Unzip`

```go
// Create a Pair of Optionals.
pair := goptional.Pair[goptional.Optional[int], goptional.Optional[string]]{
    First:  goptional.Of(123),
    Second: goptional.Of("gm"),
}

// Unwrap the given Optional Pair.
opt1, opt2 := goptional.Unzip(goptional.Of(&pair))

fmt.Println(opt1.Get()) // 123
fmt.Println(opt2.Get()) // gm

// Create an empty Optional Pair.
emptyPair := goptional.Empty[*goptional.Pair[goptional.Optional[int], goptional.Optional[string]]]()

// Return two empty Optionals if the given Optional is empty.
opt1, opt2 = goptional.Unzip(emptyPair)

fmt.Println(opt1.IsEmpty()) // true
fmt.Println(opt2.IsEmpty()) // true
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
