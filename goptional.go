package goptional

import "reflect"

// Optional represents an optional value.
type Optional[T any] struct {
	wrappedValue *valueWrapper[T]
}

type valueWrapper[T any] struct {
	value T
}

const noSuchElementErrMsg = "no value present"

// Empty returns a new Optional instance that does not hold a value.
func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

// Of returns a new Optional instance that holds a value.
// If the value passed to it is the zero value of its type, it returns an empty Optional instance instead.
func Of[T any](value T) *Optional[T] {
	if reflect.ValueOf(&value).Elem().IsZero() {
		return Empty[T]()
	}
	return &Optional[T]{
		wrappedValue: &valueWrapper[T]{value: value},
	}
}

// IsPresent returns true if the Optional instance holds a value, and false if it is empty.
func (o *Optional[T]) IsPresent() bool {
	return o.wrappedValue != nil
}

// IsEmpty returns true if the Optional instance is empty, and false if it holds a value.
func (o *Optional[T]) IsEmpty() bool {
	return o.wrappedValue == nil
}

// Get returns the value held by the Optional instance.
// If the instance is empty, it will cause a panic with the message "no value present".
func (o *Optional[T]) Get() T {
	if o.IsEmpty() {
		panic(noSuchElementErrMsg)
	}
	return o.wrappedValue.value
}

// IfPresent takes a function as input and calls it with the value held by the Optional instance if it holds a value.
// If the Optional instance is empty, it does nothing.
func (o *Optional[T]) IfPresent(action func(T)) {
	if o.IsPresent() {
		action(o.Get())
	}
}

// IfPresentOrElse takes a function and an "empty action" function as input.
// It calls the first function with the value held by the Optional instance if it holds a value,
// or the "empty action" function if it is empty.
func (o *Optional[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	if o.IsPresent() {
		action(o.Get())
	} else {
		emptyAction()
	}
}

// Filter returns a new Optional instance holding the value of the original instance if it holds a value and the value satisfies the given predicate.
// If the original instance is empty or the value does not satisfy the predicate, it returns an empty Optional instance.
func (o *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	if o.IsEmpty() {
		return o
	}
	if predicate(o.Get()) {
		return o
	}
	return Empty[T]()
}

// Map takes an Optional instance and a function as input,
// and returns a new Optional instance holding the result of applying the function to the value held by the input instance (if it holds a value).
// If the input instance is empty, it returns an empty Optional instance.
func Map[X, Y any](input *Optional[X], mapper func(X) Y) *Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return Of(mapper(input.Get()))
}

// MapOr is similar to Map, but if the input Optional instance is empty, it returns a new Optional instance holding a default value instead.
func MapOr[X, Y any](input *Optional[X], mapper func(X) Y, other Y) *Optional[Y] {
	if input.IsEmpty() {
		return Of(other)
	}
	return Of(mapper(input.Get()))
}

// MapOrElse is similar to MapOr, but allows you to specify a function to supply the default value to use if the input Optional instance is empty.
func MapOrElse[X, Y any](input *Optional[X], mapper func(X) Y, supplier func() Y) *Optional[Y] {
	if input.IsEmpty() {
		return Of(supplier())
	}
	return Of(mapper(input.Get()))
}

// FlatMap is similar to Map, but allows the function it takes as input to return an Optional instance.
// If the input Optional instance is empty, it returns an empty Optional instance.
// If the input Optional instance holds a value and the function returns an Optional instance that holds a value,
// the returned Optional instance holds the value returned by the function.
// If the input Optional instance holds a value but the function returns an empty Optional instance, the returned Optional instance is also empty.
func FlatMap[X, Y any](input *Optional[X], mapper func(X) *Optional[Y]) *Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return mapper(input.Get())
}

// And returns an empty Optional if the original Optional is empty, otherwise it returns the Optional returned by the provided supplier.
// It panics if the provided supplier is nil.
func (o *Optional[T]) And(supplier func() *Optional[T]) *Optional[T] {
	if o.IsEmpty() {
		return o
	}
	return supplier()
}

// Xor returns an Optional instance holding the value of the original instance if it holds a value and the other Optional instance is empty,
// or if the original instance is empty and the other Optional instance holds a value.
// If both Optional instances are empty or both hold a value, it returns an empty Optional instance.
func (o *Optional[T]) Xor(o2 *Optional[T]) *Optional[T] {
	if (o.IsEmpty() && o2.IsEmpty()) || (o.IsPresent() && o2.IsPresent()) {
		return Empty[T]()
	}
	if o.IsPresent() {
		return o
	}
	return o2
}

// Or returns the original Optional instance if it holds a value, or a new Optional instance returned by a given supplier if the original instance is empty.
func (o *Optional[T]) Or(supplier func() *Optional[T]) *Optional[T] {
	if o.IsPresent() {
		return o
	}
	return supplier()
}

// OrElse returns the value held by the original Optional instance if it holds a value, or a default value if it is empty.
func (o *Optional[T]) OrElse(other T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return other
}

// OrElseGet returns the value held by the original Optional instance if it holds a value, or a value supplied by a given function if it is empty.
func (o *Optional[T]) OrElseGet(supplier func() T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return supplier()
}

// OrPanic returns the value held by the original Optional instance if it holds a value, or panics with an error message supplied by a given function if it is empty.
func (o *Optional[T]) OrElsePanicWithErr(supplier func() error) T {
	if o.IsEmpty() {
		panic(supplier().Error())
	}
	return o.Get()
}
