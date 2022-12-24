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

func (o *Optional[T]) Get() T {
	if !o.IsPresent() {
		panic(noSuchElementErrMsg)
	}
	return o.wrappedValue.value
}

func (o *Optional[T]) IsPresent() bool {
	return o.wrappedValue != nil
}

func (o *Optional[T]) IsEmpty() bool {
	return o.wrappedValue == nil
}

func (o *Optional[T]) IfPresent(action func(T)) {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) IfPresentOrElse(action func(T), runnable func()) {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) Filter(predicate func(T) bool) *Optional[T] {
	//TODO implement me
	panic("implement me")
}

func Map[X, Y any](input *Optional[X], mapper func(X) Y) Optional[Y] {
	//TODO implement me
	panic("implement me")
}

func FlatMap[X, Y any](input *Optional[X], mapper func(X) *Optional[Y]) *Optional[Y] {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) Or(supplier func() Optional[T]) *Optional[T] {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) OrElse(other T) T {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) OrElseGet(supplier func() T) T {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) OrElsePanic() T {
	//TODO implement me
	panic("implement me")
}

func (o *Optional[T]) OrElsePanicWithErr(f func() error) T {
	//TODO implement me
	panic("implement me")
}
