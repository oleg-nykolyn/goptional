/*
Package goptional implements the Optional type, its methods and functions.
The API is heavily inspired by the Java & Rust implementations.
*/
package goptional

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/davecgh/go-spew/spew"
)

// Optional represents an optional value.
// At any time it can either hold a value or be empty.
type Optional[T any] []T

// ErrNoValue is an error that is returned when attempting to retrieve a value from an empty Optional.
var ErrNoValue = errors.New("no value present")

// Empty returns a new empty Optional.
func Empty[T any]() Optional[T] {
	return nil
}

// Of returns a new Optional that holds the given value.
// It panics if value is nil or invalid.
func Of[T any](value T) (opt Optional[T]) {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		panic("value is invalid")
	}
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if v.IsNil() {
			panic("value is nil")
		}
		fallthrough
	default:
		return []T{value}
	}
}

// OfNillable returns a new Optional that holds the given pointer.
// If value is nil, an empty Optional is returned instead.
func OfNillable[T any](value *T) (opt Optional[*T]) {
	if value == nil {
		return Empty[*T]()
	}
	return Of(value)
}

// IsPresent returns true if this instance holds a value, and false otherwise.
func (o Optional[T]) IsPresent() bool {
	return o != nil
}

// IsEmpty returns true if this instance is empty, and false otherwise.
func (o Optional[T]) IsEmpty() bool {
	return o == nil
}

// Get returns the value held by this instance.
//
// It panics if this instance is empty.
func (o Optional[T]) Get() T {
	if o.IsEmpty() {
		panic(ErrNoValue)
	}
	return o[0]
}

// IfPresent applies the action to the value held by this instance.
// Does nothing if this instance is empty.
//
// It panics if action is nil and this instance is not empty.
func (o Optional[T]) IfPresent(action func(T)) {
	if o.IsPresent() {
		action(o.Get())
	}
}

// IfPresentOrElse applies the action to the value held by this instance or calls emptyAction if this instance is empty.
//
// It panics if one of these is true:
//   - action is nil and this instance is not empty
//   - emptyAction is nil and this instance is empty
func (o Optional[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	if o.IsPresent() {
		action(o.Get())
	} else {
		emptyAction()
	}
}

// Filter returns self if self is empty or
// if the predicate applied to its value returns false.
//
// It panics if predicate is nil and this instance is not empty.
func (o Optional[T]) Filter(predicate func(T) bool) Optional[T] {
	if o.IsEmpty() || predicate(o.Get()) {
		return o
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
//   - self if self is empty
//   - a new Optional provided by the given supplier
//
// It panics if this instance is not empty and supplier is nil.
func (o Optional[T]) And(supplier func() Optional[T]) Optional[T] {
	if o.IsEmpty() {
		return o
	}
	return supplier()
}

// Or returns one of the following:
//   - self if self is not empty
//   - a new Optional provided by the given supplier
//
// It panics if this instance is empty and supplier is nil.
func (o Optional[T]) Or(supplier func() Optional[T]) Optional[T] {
	if o.IsPresent() {
		return o
	}
	return supplier()
}

// Xor returns one of the following:
//   - an empty Optional if both are either non-empty or empty
//   - the first non-empty Optional between this instance & opt
func (o Optional[T]) Xor(opt Optional[T]) Optional[T] {
	if (o.IsPresent() && opt.IsPresent()) || (o.IsEmpty() && opt.IsEmpty()) {
		return Empty[T]()
	}
	if o.IsPresent() {
		return o
	}
	return opt
}

// OrZero returns the value held by this instance, if there is any, or the zero value of T otherwise.
func (o Optional[T]) OrZero() T {
	if o.IsEmpty() {
		var zero T
		return zero
	}
	return o.Get()
}

// OrElse returns the value held by this instance, if there is any, or the given value otherwise.
func (o Optional[T]) OrElse(fallback T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return fallback
}

// OrElseGet returns the value held by this instance, if there is any, or a value provided by the given supplier otherwise.
//
// It panics if this instance is empty and supplier is nil.
func (o Optional[T]) OrElseGet(supplier func() T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return supplier()
}

// OrElsePanicWithErr returns the value held by this instance, if there is any, or panics with an error provided by the given supplier otherwise.
//
// It panics if this instance is empty and supplier is nil.
func (o Optional[T]) OrElsePanicWithErr(supplier func() error) T {
	if o.IsEmpty() {
		err := supplier()
		if err == nil {
			panic(ErrNoValue)
		} else {
			panic(err)
		}
	}
	return o.Get()
}

// Equals compares two Optionals for equality.
// It returns true if both Optionals contain the same value, or if both Optionals are empty.
// Otherwise, it returns false.
func (o Optional[T]) Equals(o2 Optional[T]) bool {
	if !o.IsPresent() && !o2.IsPresent() {
		return true
	}
	if o.IsPresent() && o2.IsPresent() {
		return reflect.DeepEqual(o.Get(), o2.Get())
	}
	return false
}

var nilAsJSON = []byte("null")

// MarshalJSON returns the JSON representation of this instance.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return nilAsJSON, nil
	}
	return json.Marshal(o.Get())
}

// UnmarshalJSON attempts to populate this instance with the given JSON data.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, nilAsJSON) {
		*o = Empty[T]()
		return nil
	}

	var val T
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	*o = Of(val)

	return nil
}

// String returns the string representation of this instance.
func (o Optional[T]) String() string {
	if o.IsEmpty() {
		return "Optional.empty"
	}
	return spew.Sprintf("Optional[%#+v]", o.Get())
}

// Take takes the value out of this instance, if there is any, leaving an empty Optional in its place.
func (o *Optional[T]) Take() Optional[T] {
	if o.IsEmpty() {
		return *o
	}
	v := o.Get()
	*o = nil
	return Of(v)
}

// Replace replaces the value in this instance with the given value,
// returning the old value if present, leaving a non-empty Optional in its place
// without deinitializing either one.
func (o *Optional[T]) Replace(value T) Optional[T] {
	inOpt := Of(value)
	if o.IsEmpty() {
		*o = inOpt
		return Empty[T]()
	}
	v := o.Get()
	*o = inOpt
	return Of(v)
}

// Pair is your usual generic pair.
type Pair[X, Y any] struct {
	// First is the first element of the pair.
	First X
	// Second is the second element of the pair.
	Second Y
}

// Zip zips o1 with o2.
// If o1 and o2 are both non-empty, it returns an optional pair holding the value of o1 & o2.
//
// Otherwise, an empty Optional is returned.
func Zip[X, Y any](o1 Optional[X], o2 Optional[Y]) Optional[*Pair[X, Y]] {
	if o1.IsPresent() && o2.IsPresent() {
		return Of(&Pair[X, Y]{First: o1.Get(), Second: o2.Get()})
	}
	return Empty[*Pair[X, Y]]()
}

// Unzip unzips o containing a tuple of two Optionals.
// If o is empty, it returns the unwrapped pair. Otherwise, two empty Optionals are returned.
func Unzip[X, Y any](o Optional[*Pair[Optional[X], Optional[Y]]]) (Optional[X], Optional[Y]) {
	if o.IsPresent() {
		pair := o.Get()
		return pair.First, pair.Second
	}
	return Empty[X](), Empty[Y]()
}

// ZipWith zips o1 with o2.
// If o1 and o2 are both non-empty, it returns an Optional with a value
// that results from the application of the given mapper to the value of o1 & o2.
// Otherwise, an empty Optional is returned.
//
// It panics if o1 & o2 are both non-empty and mapper is nil.
func ZipWith[X, Y, Z any](o1 Optional[X], o2 Optional[Y], mapper func(X, Y) Z) Optional[Z] {
	if o1.IsPresent() && o2.IsPresent() {
		return Of(mapper(o1.Get(), o2.Get()))
	}
	return Empty[Z]()
}
