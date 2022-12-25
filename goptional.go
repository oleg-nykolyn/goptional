/*
Package goptional implements the Optional type, its methods and functions.

The API is heavily inspired by Java, specifically https://github.com/AdoptOpenJDK/openjdk-jdk11/blob/master/src/java.base/share/classes/java/util/Optional.java –
with some additions inspired by Rust, such as MapOr & MapOrElse.

The family of map-like operations is implemented through functions
due to Go's absence of method-level type parameters. This might change in the future.
*/
package goptional

import (
	"reflect"
)

// Optional represents an optional value.
//
// In essence, it's wrapper around a pointer and therefore can be passed and returned by value without sacrificing performance.
type Optional[T any] struct {
	wrappedValue *valueWrapper[T]
}

type valueWrapper[T any] struct {
	value T
}

const noValueErrMsg = "no value present"

// Empty returns a new empty Optional.
func Empty[T any]() Optional[T] {
	return Optional[T]{}
}

// Of returns a new Optional that holds the given value.
// If value is nil, it returns an empty Optional instead.
func Of[T any](value T) (opt Optional[T]) {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if v.IsNil() {
			return
		}
		fallthrough
	default:
		return Optional[T]{
			wrappedValue: &valueWrapper[T]{value: value},
		}
	}
}

// IsPresent returns true if the Optional holds a value, and false otherwise.
func (o *Optional[T]) IsPresent() bool {
	return o.wrappedValue != nil
}

// IsEmpty returns true if the Optional is empty, and false otherwise.
func (o *Optional[T]) IsEmpty() bool {
	return o.wrappedValue == nil
}

// Get returns the value held by the Optional.
//
// It panics if the Optional is empty.
func (o *Optional[T]) Get() T {
	if o.IsEmpty() {
		panic(noValueErrMsg)
	}
	return o.wrappedValue.value
}

// IfPresent applies the action to the value held by the Optional.
// Does nothing if the Optional is empty.
//
// It panics if action is nil and the Optional is not empty.
func (o *Optional[T]) IfPresent(action func(T)) {
	if o.IsPresent() {
		action(o.Get())
	}
}

// IfPresentOrElse applies the action to the value held by the Optional or calls emptyAction if the Optional is empty.
//
// It panics if one of these is true:
//   - action is nil and the Optional is not empty
//   - emptyAction is nil and the Optional is empty
func (o *Optional[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	if o.IsPresent() {
		action(o.Get())
	} else {
		emptyAction()
	}
}

// Filter returns a copy of self if self is empty or
// if the predicate applied to its value returns false.
//
// It panics if predicate is nil.
func (o *Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.IsEmpty() || predicate(o.Get()) {
		return *o
	}
	return Empty[T]()
}

// Map returns one of the following:
//
//   - an empty Optional if input is empty
//   - a new Optional holding a value that results from the application of the given mapper to the value of input
//
// It panics if mapper is nil and input is not empty.
func Map[X, Y any](input Optional[X], mapper func(X) Y) Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return Of(mapper(input.Get()))
}

// MapOr is similar to Map, but if input is empty, it returns a new Optional holding a default value instead.
//
// It panics if mapper is nil and input is not empty.
func MapOr[X, Y any](input Optional[X], mapper func(X) Y, other Y) Optional[Y] {
	if input.IsEmpty() {
		return Of(other)
	}
	return Of(mapper(input.Get()))
}

// MapOrElse is similar to MapOr, but if input is empty, it returns a new Optional holding the value provided by the given supplier.
//
// It panics if one of these is true:
//   - supplier is nil and input is empty
//   - mapper is nil and input is not empty
func MapOrElse[X, Y any](input Optional[X], mapper func(X) Y, supplier func() Y) Optional[Y] {
	if input.IsEmpty() {
		return Of(supplier())
	}
	return Of(mapper(input.Get()))
}

// FlatMap returns one of the following:
//
//   - an empty Optional if input is empty
//   - a new Optional that results from the application of the given mapper to the value of input
//
// It panics if mapper is nil and input is not empty.
func FlatMap[X, Y any](input Optional[X], mapper func(X) Optional[Y]) Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return mapper(input.Get())
}

// And returns one of the following:
//   - a copy of self if self is empty
//   - a new Optional provided by the given supplier
//
// It panics if the Optional is not empty and supplier is nil.
func (o *Optional[T]) And(supplier func() Optional[T]) Optional[T] {
	if o.IsEmpty() {
		return *o
	}
	return supplier()
}

// Or returns one of the following:
//   - a copy of self if self is not empty
//   - a new Optional provided by the given supplier
//
// It panics if the Optional is not empty and supplier is nil.
func (o *Optional[T]) Or(supplier func() Optional[T]) Optional[T] {
	if o.IsPresent() {
		return *o
	}
	return supplier()
}

// OrElse returns the value held by the Optional if it's not empty, or the given value otherwise.
func (o *Optional[T]) OrElse(other T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return other
}

// OrElseGet returns the value held by the Optional if it's not empty, or a value provided by the given supplier otherwise.
//
// It panics if the Optional is empty and supplier is nil.
func (o *Optional[T]) OrElseGet(supplier func() T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return supplier()
}

// OrElsePanic returns the value held by the Optional if it's not empty, or panics with an error message provided by the given supplier otherwise.
//
// It panics if the Optional is empty and supplier is nil.
func (o *Optional[T]) OrElsePanic(supplier func() string) T {
	if o.IsEmpty() {
		panic(supplier())
	}
	return o.Get()
}
