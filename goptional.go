package goptional

import "reflect"

// Optional represents an optional value.
// Every Optional is either populated with a value or empty.
type Optional[T any] struct {
	wrappedValue *valueWrapper[T]
}

type valueWrapper[T any] struct {
	value T
}

const noSuchElementErrMsg = "no value present"

// Empty returns a new empty Optional.
func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

// Of returns a new Optional that holds the given value.
// If it is the zero value of its type, it returns an empty Optional instead.
func Of[T any](value T) *Optional[T] {
	if reflect.ValueOf(&value).Elem().IsZero() {
		return Empty[T]()
	}
	return &Optional[T]{
		wrappedValue: &valueWrapper[T]{value: value},
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
		panic(noSuchElementErrMsg)
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

// Filter returns self if the Optional is empty or
// if the predicate applied to its value returns false.
//
// It panics if predicate is nil.
func (o *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	if o.IsEmpty() {
		return o
	}
	if predicate(o.Get()) {
		return o
	}
	return Empty[T]()
}

// Map returns one of the following:
//
//   - input if the Optional is empty
//   - a new Optional holding a value that results from the application of the given mapper to the input Optional value
//
// It panics if one of these is true:
//   - input is nil
//   - mapper is nil and the input Optional is not empty
func Map[X, Y any](input *Optional[X], mapper func(X) Y) *Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return Of(mapper(input.Get()))
}

// MapOr is similar to Map, but if the input Optional is empty, it returns a new Optional holding a default value instead.
//
// It panics if one of these is true:
//   - input is nil
//   - mapper is nil and the input Optional is not empty
func MapOr[X, Y any](input *Optional[X], mapper func(X) Y, other Y) *Optional[Y] {
	if input.IsEmpty() {
		return Of(other)
	}
	return Of(mapper(input.Get()))
}

// MapOrElse is similar to MapOr, but if the input Optional is empty, it returns a new Optional holding the value provided by the given supplier.
//
// It panics if one of these is true:
//   - input is nil
//   - supplier is nil and the input Optional is empty
//   - mapper is nil and the input Optional is not empty
func MapOrElse[X, Y any](input *Optional[X], mapper func(X) Y, supplier func() Y) *Optional[Y] {
	if input.IsEmpty() {
		return Of(supplier())
	}
	return Of(mapper(input.Get()))
}

// FlatMap returns one of the following:
//
//   - an empty Optional if the input Optional is empty
//   - a new Optional that results from applying the mapper to the input Optional value
//
// It panics if one of these is true:
//   - input is nil
//   - mapper is nil and the input Optional is not empty
func FlatMap[X, Y any](input *Optional[X], mapper func(X) *Optional[Y]) *Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return mapper(input.Get())
}

// And returns one of the following:
//   - self if the Optional is empty
//   - a new Optional provided by the given supplier
//
// It panics if the Optional is not empty and supplier is nil.
func (o *Optional[T]) And(supplier func() *Optional[T]) *Optional[T] {
	if o.IsEmpty() {
		return o
	}
	return supplier()
}

// Or returns one of the following:
//   - self if the Optional is not empty
//   - a new Optional provided by the given supplier
//
// It panics if the Optional is not empty and supplier is nil.
func (o *Optional[T]) Or(supplier func() *Optional[T]) *Optional[T] {
	if o.IsPresent() {
		return o
	}
	return supplier()
}

// OrElse returns the value held by the Optional if it not empty, or the given value otherwise.
func (o *Optional[T]) OrElse(other T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return other
}

// OrElseGet returns the value held by the Optional if it not empty, or a value provided by the given supplier otherwise.
//
// It panics if the Optional is empty and supplier is nil.
func (o *Optional[T]) OrElseGet(supplier func() T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return supplier()
}

// OrPanic returns the value held by the Optional if it not empty, or panics with an error message provided by the given supplier.
//
// It panics if the Optional is empty and supplier is nil.
func (o *Optional[T]) OrElsePanicWithErr(supplier func() error) T {
	if o.IsEmpty() {
		panic(supplier().Error())
	}
	return o.Get()
}
