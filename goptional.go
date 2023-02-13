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
type Optional[T any] struct {
	value        T
	isValueValid bool
}

// ErrNoValue is returned when attempting to retrieve a value from an empty Optional.
var ErrNoValue = errors.New("no value present")

// ErrMutationOnNil is returned when attempting to mutate a nil Optional instance.
var ErrMutationOnNil = errors.New("cannot mutate nil Optional instance")

// Empty returns a new empty Optional.
func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

// Of attempts to return a new non-empty Optional wrapping the given value.
// If such value is either invalid or nil, it returns an empty Optional instead.
func Of[T any](value T) *Optional[T] {
	// There are two sensible implementation choices here: reflection and a Nillable interface.
	// The third _undesirable_ choice would be to skip the IsNil check and
	// trust the caller: needless to say, this could inevitably lead to very unpredictable API results.
	//
	// goptional adopts reflection to avoid the hassle of forcing
	// the caller to implement Nillable on all types.

	if isValueValid(value) {
		return &Optional[T]{value: value, isValueValid: true}
	}

	return Empty[T]()
}

// IsPresent returns true if this instance holds a value, and false otherwise.
func (o *Optional[T]) IsPresent() bool {
	return o != nil && o.isValueValid
}

// IsEmpty returns true if this instance does not hold a value, and false otherwise.
func (o *Optional[T]) IsEmpty() bool {
	return o == nil || !o.isValueValid
}

// Unwrap returns the value held by this instance, if any, or _panics_ otherwise.
//
// Use it only if you _know_ what you are doing.
// Usage of OrDefault / OrElse is preferred.
func (o *Optional[T]) Unwrap() T {
	if o.IsEmpty() {
		panic(ErrNoValue)
	}

	return o.value
}

// IfPresent applies the action to the value held by this instance.
// Does nothing if this instance is empty. If action is nil, nothing is done.
func (o *Optional[T]) IfPresent(action func(T)) {
	if o.IsPresent() && action != nil {
		action(o.Unwrap())
	}
}

// IfPresentOrElse applies the action to the value held by this instance or calls emptyAction if this instance is empty.
// If action or emptyAction are nil, nothing is done.
func (o *Optional[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	if o.IsPresent() {
		if action != nil {
			action(o.Unwrap())
		}
	} else {
		if emptyAction != nil {
			emptyAction()
		}
	}
}

// Filter returns this instance if it is empty or
// if the predicate applied to its value returns false.
// If this instance is not empty and predicate is nil, it returns an empty Optional.
func (o *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	if o.IsEmpty() {
		return o
	}

	if predicate == nil {
		return Empty[T]()
	}

	if predicate(o.Unwrap()) {
		return o
	}

	return Empty[T]()
}

// Map returns one of the following:
//   - an empty Optional if input is empty
//   - a new Optional holding a value that results from the application of the given mapper to the value of input
//
// If input is not empty and mapper is nil, it returns an empty Optional of the target type.
func Map[X, Y any](input *Optional[X], mapper func(X) Y) *Optional[Y] {
	if input.IsEmpty() || mapper == nil {
		return Empty[Y]()
	}

	return Of(mapper(input.Unwrap()))
}

// MapOr is similar to Map, but if input is empty, it returns a new Optional holding a default value instead.
// If input is not empty and mapper is nil, it returns an empty Optional of the target type.
func MapOr[X, Y any](input *Optional[X], mapper func(X) Y, other Y) *Optional[Y] {
	if input.IsEmpty() {
		return Of(other)
	}

	if mapper == nil {
		return Empty[Y]()
	}

	return Of(mapper(input.Unwrap()))
}

// MapOrElse is similar to MapOr, but if input is empty, it returns a new Optional holding the value provided by the given supplier.
//
// It returns an empty Optional of the target type if one of these is true:
//   - input is empty and supplier is nil
//   - input holds a value and mapper is nil
func MapOrElse[X, Y any](input *Optional[X], mapper func(X) Y, supplier func() Y) *Optional[Y] {
	if input.IsEmpty() {
		if supplier != nil {
			return Of(supplier())
		}
		return Empty[Y]()
	}

	if mapper == nil {
		return Empty[Y]()
	}

	return Of(mapper(input.Unwrap()))
}

// FlatMap returns one of the following:
//   - an empty Optional if input is empty
//   - a new Optional that results from the application of the given mapper to the value of input
//
// If input is not empty and mapper is nil, it returns an empty Optional of the target type.
func FlatMap[X, Y any](input *Optional[X], mapper func(X) *Optional[Y]) *Optional[Y] {
	if input.IsEmpty() || mapper == nil {
		return Empty[Y]()
	}

	return mapper(input.Unwrap())
}

// And returns one of the following:
//   - this instance if it is empty
//   - a new Optional provided by the given supplier
//
// If this instance is not empty and supplier is nil, it returns an empty Optional.
func (o *Optional[T]) And(supplier func() *Optional[T]) *Optional[T] {
	if o.IsEmpty() {
		return o
	}

	if supplier == nil {
		return Empty[T]()
	}

	return supplier()
}

// Or returns one of the following:
//   - this instance if it holds a value
//   - a new Optional provided by the given supplier
//
// It returns this instance if it is empty and supplier is nil.
func (o *Optional[T]) Or(supplier func() *Optional[T]) *Optional[T] {
	if o.IsPresent() || supplier == nil {
		return o
	}

	return supplier()
}

// Xor returns one of the following:
//   - an empty Optional if both are either non-empty or empty
//   - the first non-empty Optional between this instance & o2
func (o *Optional[T]) Xor(o2 *Optional[T]) *Optional[T] {
	if (o.IsPresent() && o2.IsPresent()) || (o.IsEmpty() && o2.IsEmpty()) {
		return Empty[T]()
	}

	if o.IsPresent() {
		return o
	}

	return o2
}

// OrDefault returns the value held by this instance, if any, or the zero value of T otherwise e.g.
// string -> "", bool -> false, int -> 0, ptr -> nil, etc.
func (o *Optional[T]) OrDefault() T {
	if o.IsEmpty() {
		return getZeroOfType[T]()
	}

	return o.Unwrap()
}

// OrElse returns the value held by this instance, if any, or the given value otherwise.
func (o *Optional[T]) OrElse(fallback T) T {
	if o.IsPresent() {
		return o.Unwrap()
	}

	return fallback
}

// OrElseGet returns the value held by this instance, if any, or a value provided by the given supplier otherwise.
//
// If this instance is empty and supplier is nil, it returns the zero value of T.
func (o *Optional[T]) OrElseGet(supplier func() T) T {
	if o.IsPresent() {
		return o.Unwrap()
	}

	if supplier == nil {
		return getZeroOfType[T]()
	}

	return supplier()
}

// UnwrapOr returns the value held by this instance, if any, or panics with an error provided by the given supplier otherwise.
//
// If this instance is empty and supplier is either nil or returns a nil error, it panics with ErrNoValue instead.
func (o *Optional[T]) UnwrapOr(supplier func() error) T {
	if o.IsEmpty() {
		if supplier == nil {
			panic(ErrNoValue)
		}

		if err := supplier(); err != nil {
			panic(err)
		} else {
			panic(ErrNoValue)
		}
	}

	return o.Unwrap()
}

// Equals compares two Optionals for deep equality.
// It returns true if both Optionals contain the same value, or if both Optionals are empty.
// It returns false otherwise.
func (o *Optional[T]) Equals(o2 *Optional[T]) bool {
	if !o.IsPresent() && !o2.IsPresent() {
		return true
	}

	if o.IsPresent() && o2.IsPresent() {
		return reflect.DeepEqual(o.Unwrap(), o2.Unwrap())
	}

	return false
}

// EqualsBy compares two Optionals for equality through a custom predicate.
// It returns true if both Optionals are not empty and the given predicate applied to their values holds,
// or if both Optionals are empty.
// It returns false otherwise.
//
// If the given predicate is nil, it is substituted with reflect.DeepEqual
func (o *Optional[T]) EqualsBy(o2 *Optional[T], predicate func(v1, v2 T) bool) bool {
	if !o.IsPresent() && !o2.IsPresent() {
		return true
	}

	if o.IsPresent() && o2.IsPresent() {
		if predicate != nil {
			return predicate(o.Unwrap(), o2.Unwrap())
		}
		return reflect.DeepEqual(o.Unwrap(), o2.Unwrap())
	}

	return false
}

var nilAsJSON = []byte("null")

// MarshalJSON returns the JSON representation of this instance.
func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return nilAsJSON, nil
	}

	return json.Marshal(o.Unwrap())
}

// UnmarshalJSON attempts to populate this instance with the given JSON data.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if o == nil {
		return ErrMutationOnNil
	}

	if len(data) == 0 || bytes.Equal(data, nilAsJSON) {
		o.unsetValue()
		return nil
	}

	var value T
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	o.setValue(value)

	return nil
}

// String returns the string representation of this instance.
func (o *Optional[T]) String() string {
	if o.IsEmpty() {
		return "Optional.empty"
	}

	return spew.Sprintf("Optional[%#+v]", o.Unwrap())
}

// Take takes the value out of this instance, if any, leaving an empty Optional in its place.
func (o *Optional[T]) Take() *Optional[T] {
	if o.IsEmpty() {
		return o
	}

	v := o.Unwrap()
	o.unsetValue()
	return Of(v)
}

// Replace replaces the value in this instance with the given value,
// returning the old value if present, leaving a non-empty Optional in its place
// without deinitializing either one.
// It returns an ErrMutationOnNil error otherwise.
func (o *Optional[T]) Replace(value T) (*Optional[T], error) {
	if o == nil {
		return nil, ErrMutationOnNil
	}

	if !o.isValueValid {
		o.setValue(value)
		return Empty[T](), nil
	}

	v := o.Unwrap()
	o.setValue(value)
	return Of(v), nil
}

// Pair is your usual generic pair.
type Pair[X, Y any] struct {
	// First is the first element of the pair.
	First X
	// Second is the second element of the pair.
	Second Y
}

// Zip zips o1 with o2.
// If o1 and o2 are both non-empty, it returns an Optional Pair holding the values of o1 & o2.
//
// It returns an empty Optional otherwise.
func Zip[X, Y any](o1 *Optional[X], o2 *Optional[Y]) *Optional[*Pair[X, Y]] {
	if o1.IsPresent() && o2.IsPresent() {
		return Of(&Pair[X, Y]{First: o1.Unwrap(), Second: o2.Unwrap()})
	}

	return Empty[*Pair[X, Y]]()
}

// Unzip unzips the given optional pair.
// If o is not empty, it returns the unwrapped pair, or two empty Optionals otherwise.
func Unzip[X, Y any](o *Optional[*Pair[*Optional[X], *Optional[Y]]]) (*Optional[X], *Optional[Y]) {
	if o.IsPresent() {
		pair := o.Unwrap()
		return pair.First, pair.Second
	}

	return Empty[X](), Empty[Y]()
}

// ZipWith zips o1 with o2.
// If o1 and o2 are both non-empty, it returns an Optional with a value
// that results from the application of the given mapper to the values of o1 & o2.
// It returns an empty Optional otherwise.
//
// It o1 & o2 are both non-empty and mapper is nil, it returns an empty Optional of the target type.
func ZipWith[X, Y, Z any](o1 *Optional[X], o2 *Optional[Y], mapper func(X, Y) Z) *Optional[Z] {
	if o1.IsPresent() && o2.IsPresent() {
		if mapper == nil {
			return Empty[Z]()
		}

		return Of(mapper(o1.Unwrap(), o2.Unwrap()))
	}

	return Empty[Z]()
}

// Flatten flattens the given Optional.
func Flatten[T any](o *Optional[*Optional[T]]) *Optional[T] {
	if o.IsPresent() {
		return o.Unwrap()
	}

	return Empty[T]()
}

// Is checks if the value of this instance satisfies the given predicate.
// If this instance is empty, it returns false.
//
// If this instance is not empty and predicate is nil, it returns false.
func (o *Optional[T]) Is(predicate func(T) bool) bool {
	if o.IsEmpty() || predicate == nil {
		return false
	}

	return predicate(o.Unwrap())
}

// Val returns the value held by this instance, if any. It returns ErrNoValue otherwise.
func (o *Optional[T]) Val() (T, error) {
	if o.IsPresent() {
		return o.Unwrap(), nil
	}

	return getZeroOfType[T](), ErrNoValue
}

// ValOr returns the value held by this instance, if any. It returns the given error otherwise.
// On the other hand, if this instance is empty and err is nil, it returns ErrNoValue.
func (o *Optional[T]) ValOr(err error) (T, error) {
	if o.IsPresent() {
		return o.Unwrap(), nil
	}

	if err == nil {
		return getZeroOfType[T](), ErrNoValue
	}

	return getZeroOfType[T](), err
}

// ValOrElse returns the value held by this instance, if any.
// It returns the error provided by the given supplier otherwise.
//
// If this instance is empty and supplier is either nil or returns a nil err, it returns ErrNoValue.
func (o *Optional[T]) ValOrElse(supplier func() error) (T, error) {
	if o.IsPresent() {
		return o.Unwrap(), nil
	}

	if supplier == nil {
		return getZeroOfType[T](), ErrNoValue
	}

	if err := supplier(); err != nil {
		return getZeroOfType[T](), err
	}
	return getZeroOfType[T](), ErrNoValue
}

// setValue changes the state of this instance to non-empty
// iff the given value is valid, or to empty otherwise.
func (o *Optional[T]) setValue(value T) {
	if o == nil {
		// should never happen, by design.
		return
	}

	if isValueValid(value) {
		o.value = value
		o.isValueValid = true
	} else {
		o.unsetValue()
	}
}

// unsetValue changes the state of this instance to empty.
func (o *Optional[T]) unsetValue() {
	if o == nil {
		// should never happen, by design.
		return
	}

	o.value = getZeroOfType[T]()
	o.isValueValid = false
}

// isValueValid returns false if the given value is either nil or invalid.
func isValueValid[T any](value T) bool {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return false
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		if v.IsNil() {
			return false
		}
		fallthrough
	default:
		return true
	}
}

func getZeroOfType[T any]() T {
	var zero T
	return zero
}
