package goptional

import "reflect"

type Optional[T any] struct {
	wrappedValue *valueWrapper[T]
}

type valueWrapper[T any] struct {
	value T
}

const noSuchElementErrMsg = "no value present"

func Empty[T any]() *Optional[T] {
	return &Optional[T]{}
}

func Of[T any](value T) *Optional[T] {
	if reflect.ValueOf(&value).Elem().IsZero() {
		return Empty[T]()
	}
	return &Optional[T]{
		wrappedValue: &valueWrapper[T]{value: value},
	}
}

func (o *Optional[T]) IsPresent() bool {
	return o.wrappedValue != nil
}

func (o *Optional[T]) IsEmpty() bool {
	return o.wrappedValue == nil
}

func (o *Optional[T]) Get() T {
	if o.IsEmpty() {
		panic(noSuchElementErrMsg)
	}
	return o.wrappedValue.value
}

func (o *Optional[T]) IfPresent(action func(T)) {
	if o.IsPresent() {
		action(o.Get())
	}
}

func (o *Optional[T]) IfPresentOrElse(action func(T), emptyAction func()) {
	if o.IsPresent() {
		action(o.Get())
	} else {
		emptyAction()
	}
}

func (o *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	if o.IsEmpty() {
		return o
	}
	if predicate(o.Get()) {
		return o
	}
	return Empty[T]()
}

func Map[X, Y any](input *Optional[X], mapper func(X) Y) *Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return Of(mapper(input.Get()))
}

func FlatMap[X, Y any](input *Optional[X], mapper func(X) *Optional[Y]) *Optional[Y] {
	if input.IsEmpty() {
		return Empty[Y]()
	}
	return mapper(input.Get())
}

func (o *Optional[T]) Or(supplier func() *Optional[T]) *Optional[T] {
	if o.IsPresent() {
		return o
	}
	return supplier()
}

func (o *Optional[T]) OrElse(other T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return other
}

func (o *Optional[T]) OrElseGet(supplier func() T) T {
	if o.IsPresent() {
		return o.Get()
	}
	return supplier()
}

func (o *Optional[T]) OrElsePanic() T {
	if o.IsEmpty() {
		panic(noSuchElementErrMsg)
	}
	return o.Get()
}

func (o *Optional[T]) OrElsePanicWithErr(errSupplier func() error) T {
	if o.IsEmpty() {
		panic(errSupplier().Error())
	}
	return o.Get()
}
