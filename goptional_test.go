package goptional

import (
	"errors"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	opt := Empty[interface{}]()
	require.Nil(t, opt)
}

func TestOf_ValidValue(t *testing.T) {
	opt := Of(123)
	require.NotEmpty(t, opt)
	require.EqualValues(t, opt[0], 123)
}

func TestOf_InvalidValue(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
		require.EqualValues(t, err, "value is invalid")
	}()
	Of[interface{}](nil)
}

func TestOf_NilValue(t *testing.T) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
		require.EqualValues(t, err, "value is nil")
	}()
	Of[*string](nil)
}

func TestOfNillable_NilValue(t *testing.T) {
	require.Empty(t, OfNillable[string](nil))
}

func TestOfNillable_NotNilValue(t *testing.T) {
	v := []int{1, 2, 3}
	opt := OfNillable(&v)
	require.NotEmpty(t, opt)
	require.EqualValues(t, opt[0], &v)
}

func TestIsPresent_Empty(t *testing.T) {
	opt := Empty[string]()
	require.False(t, opt.IsPresent())
}

func TestIsPresent_NilValue(t *testing.T) {
	opt := OfNillable[map[string]interface{}](nil)
	require.False(t, opt.IsPresent())
}

func TestIsPresent_NotEmpty(t *testing.T) {
	opt := Of("goptional")
	require.True(t, opt.IsPresent())
}

func TestIsPresent_ZeroValue(t *testing.T) {
	opt := Of("")
	require.True(t, opt.IsPresent())
}

func TestIsEmpty_Empty(t *testing.T) {
	opt := Empty[string]()
	require.True(t, opt.IsEmpty())
}

func TestIsEmpty_NilValue(t *testing.T) {
	opt := OfNillable[map[string]interface{}](nil)
	require.True(t, opt.IsEmpty())
}

func TestIsEmpty_NotEmpty(t *testing.T) {
	opt := Of("goptional")
	require.False(t, opt.IsEmpty())
}

func TestIsEmpty_ZeroValue(t *testing.T) {
	opt := Of("")
	require.False(t, opt.IsEmpty())
}

func TestGet_NotEmpty(t *testing.T) {
	s := "goptional"
	opt := Of(s)
	require.EqualValues(t, opt.Get(), s)
}

func TestGet_Empty(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.ErrorIs(t, r.(error), ErrNoValue)
	}()
	opt := Empty[string]()
	_ = opt.Get()
}

func TestGet_NilValue(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.ErrorIs(t, r.(error), ErrNoValue)
	}()
	opt := OfNillable[string](nil)
	_ = opt.Get()
}

func TestIfPresent_NotEmpty(t *testing.T) {
	optVal := 0
	Of(123).IfPresent(func(x int) { optVal = x })
	require.EqualValues(t, optVal, 123)
}

func TestIfPresent_Empty(t *testing.T) {
	called := false
	Empty[int]().IfPresent(func(_ int) { called = true })
	require.False(t, called)
}

func TestIfPresent_NilValue(t *testing.T) {
	called := false
	OfNillable[[]string](nil).IfPresent(func(_ *[]string) { called = true })
	require.False(t, called)
}

func TestIfPresent_NilActionOnEmpty(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	OfNillable[[]string](nil).IfPresent(nil)
}

func TestIfPresent_NilActionOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Of([]string{"a", "b", "c"}).IfPresent(nil)
}

func TestIfPresentOrElse_Empty(t *testing.T) {
	var actionCalled, emptyActionCalled bool
	Empty[string]().IfPresentOrElse(func(_ string) { actionCalled = true }, func() { emptyActionCalled = true })
	require.False(t, actionCalled)
	require.True(t, emptyActionCalled)
}

func TestIfPresentOrElse_NilValue(t *testing.T) {
	var actionCalled, emptyActionCalled bool
	OfNillable[string](nil).IfPresentOrElse(func(_ *string) { actionCalled = true }, func() { emptyActionCalled = true })
	require.False(t, actionCalled)
	require.True(t, emptyActionCalled)
}

func TestIfPresentOrElse_NotEmpty(t *testing.T) {
	var actionCalled, emptyActionCalled bool
	Of(123).IfPresentOrElse(func(_ int) { actionCalled = true }, func() { emptyActionCalled = true })
	require.True(t, actionCalled)
	require.False(t, emptyActionCalled)
}

func TestIfPresentOrElse_NilActionOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Of(123).IfPresentOrElse(nil, func() {})
}

func TestIfPresentOrElse_NilEmptyActionOnEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Empty[string]().IfPresentOrElse(func(_ string) {}, nil)
}

func TestIfPresentOrElse_NilEmptyActionOnNilValue(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	OfNillable[string](nil).IfPresentOrElse(func(_ *string) {}, nil)
}

func TestFilter_Empty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Filter(func(_ string) bool { return true })
	require.True(t, opt.IsEmpty())
}

func TestFilter_NilValue(t *testing.T) {
	opt := OfNillable[[]string](nil)
	opt = opt.Filter(func(_ *[]string) bool { return true })
	require.True(t, opt.IsEmpty())
}

func TestFilter_NotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Filter(func(_ int) bool { return true })
	require.True(t, opt.IsPresent())
}

func TestFilter_NilPredicateOnEmpty(t *testing.T) {
	require.True(t, Empty[string]().Filter(nil).IsEmpty())
}

func TestFilter_NilPredicateOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Of(123).Filter(nil)
}

func TestFilter_PredicateNotOkOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Filter(func(_ string) bool { return false })
	require.True(t, opt.IsEmpty())
}

func TestFilter_PredicateNotOkOnNilValue(t *testing.T) {
	opt := OfNillable[string](nil)
	opt = opt.Filter(func(_ *string) bool { return false })
	require.True(t, opt.IsEmpty())
}

func TestFilter_PredicateNotOkOnNotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Filter(func(_ int) bool { return false })
	require.True(t, opt.IsEmpty())
}

func TestMap_Empty(t *testing.T) {
	opt := Map(Empty[string](), func(s string) string { return s })
	require.True(t, opt.IsEmpty())
}

func TestMap_NilMapperOnEmpty(t *testing.T) {
	opt := Map[string, interface{}](Empty[string](), nil)
	require.True(t, opt.IsEmpty())
}

func TestMap_NotEmpty(t *testing.T) {
	opt := Map(Of(123), func(x int) string { return fmt.Sprintf("%v", x) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestMap_NilMapperOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Map[int, string](Of(123), nil)
}

func TestMap_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Map(nil, func(i int) string { return "goptional" })
}

func TestMap_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Map[bool, bool](nil, nil)
}

func TestMapOr_Empty(t *testing.T) {
	opt := MapOr(Empty[string](), func(s string) string { return s }, "default")
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "default")
}

func TestMapOr_NilMapperOnEmpty(t *testing.T) {
	opt := MapOr[string, interface{}](Empty[string](), nil, "default")
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "default")
}

func TestMapOr_NotEmpty(t *testing.T) {
	opt := MapOr(Of(123), func(x int) string { return fmt.Sprintf("%v", x) }, "default")
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestMapOr_NilMapperOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	MapOr(Of(123), nil, "default")
}

func TestMapOr_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOr(nil, func(i int) string { return "goptional" }, "default")
}

func TestMapOr_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOr[bool](nil, nil, "default")
}

func TestMapOrElse_Empty(t *testing.T) {
	opt := MapOrElse(Empty[string](), func(s string) string { return s }, func() string { return "default" })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "default")
}

func TestMapOrElse_NilMapperOnEmpty(t *testing.T) {
	opt := MapOrElse(Empty[string](), nil, func() string { return "default" })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "default")
}

func TestMapOrElse_NotEmpty(t *testing.T) {
	opt := MapOrElse(Of(123), func(x int) string { return fmt.Sprintf("%v", x) }, func() string { return "default" })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestMapOrElse_NilMapperOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	MapOrElse(Of(123), nil, func() string { return "default" })
}

func TestMapOrElse_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOrElse(nil, func(i int) string { return "goptional" }, func() string { return "default" })
}

func TestMapOrElse_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOrElse[bool](nil, nil, func() string { return "default" })
}

func TestMapOrElse_NilSupplierOnEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	MapOrElse(Empty[string](), func(x string) int { return 0 }, nil)
}

func TestFlatMap_Empty(t *testing.T) {
	opt := FlatMap(Empty[string](), func(x string) Optional[int] { return Of(123) })
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_NilMapperOnEmpty(t *testing.T) {
	opt := FlatMap[string, interface{}](Empty[string](), nil)
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_MapToNotEmptyOnNotEmpty(t *testing.T) {
	opt := FlatMap(Of(123), func(x int) Optional[string] { return Of(fmt.Sprintf("%v", x)) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestFlatMap_MapToEmptyOnNotEmpty(t *testing.T) {
	opt := FlatMap(Of(123), func(x int) Optional[string] { return Empty[string]() })
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_NilMapperOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	FlatMap[int, string](Of(123), nil)
}

func TestFlatMap_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	FlatMap(nil, func(x int) Optional[string] { return Of("123") })
}

func TestFlatMap_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	FlatMap[bool, bool](nil, nil)
}

func TestAnd_Empty(t *testing.T) {
	require.True(t, Empty[string]().And(func() Optional[string] { return Of("123") }).IsEmpty())
}

func TestAnd_NilSupplierOnEmpty(t *testing.T) {
	require.True(t, Empty[string]().And(nil).IsEmpty())
}

func TestAnd_SuppliedEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.And(func() Optional[int] { return Empty[int]() })
	require.True(t, opt.IsEmpty())
}

func TestAnd_SuppliedNotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.And(func() Optional[int] { return Of(321) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), 321)
}

func TestAnd_NilSupplierOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Of(123).And(nil)
}

func TestOr_NilSupplierOnNotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Or(nil)
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), 123)
}

func TestOr_NotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Or(func() Optional[int] { return Of(321) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), 123)
}

func TestOr_SuppliedNotEmptyOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Or(func() Optional[string] { return Of("123") })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestOr_SuppliedEmptyOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Or(func() Optional[string] { return Empty[string]() })
	require.True(t, opt.IsEmpty())
}

func TestOr_NilSupplierOnEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Empty[string]().Or(nil)
}

func TestOrElse_NotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrElse(321), 123)
}

func TestOrElse_Empty(t *testing.T) {
	require.EqualValues(t, Empty[int]().OrElse(123), 123)
}

func TestOrElseGet_NotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrElseGet(func() int { return 321 }), 123)
}

func TestOrElseGet_NilSupplierOnNotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrElseGet(nil), 123)
}

func TestOrElseGet_Empty(t *testing.T) {
	require.EqualValues(t, Empty[int]().OrElseGet(func() int { return 321 }), 321)
}

func TestOrElseGet_NilSupplierOnEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Empty[string]().OrElseGet(nil)
}

func TestOrElsePanicWithErr_NotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrElsePanicWithErr(func() error { return errors.New("woops") }), 123)
}

func TestOrElsePanicWithErr_NilSupplierOnNotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrElsePanicWithErr(nil), 123)
}

func TestOrElsePanicWithErr_Empty(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		require.Error(t, err)
		require.EqualError(t, err, "woops")
	}()
	Empty[string]().OrElsePanicWithErr(func() error { return errors.New("woops") })
}

func TestOrElsePanicWithErr_SuppliedNilOnEmpty(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		require.ErrorIs(t, err, ErrNoValue)
	}()
	Empty[string]().OrElsePanicWithErr(func() error { return nil })
}

func TestOrElsePanicWithErr_NilSupplierOnEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Empty[string]().OrElsePanicWithErr(nil)
}

func TestXor_NilOptOnEmpty(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Empty[string]().Xor(nil)
}

func TestXor_NilOptOnNotEmpty(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Of(123).Xor(nil)
}

func TestXor_BothEmpty(t *testing.T) {
	require.True(t, Empty[int]().Xor(Empty[int]()).IsEmpty())
}

func TestXor_BothNotEmpty(t *testing.T) {
	require.True(t, Of(123).Xor(Of(321)).IsEmpty())
}

func TestXor_FirstEmpty(t *testing.T) {
	opt := Empty[int]().Xor(Of(321))
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), 321)
}

func TestXor_SecondEmpty(t *testing.T) {
	opt := Of(123).Xor(Empty[int]())
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), 123)
}

func TestString_Empty(t *testing.T) {
	require.EqualValues(t, Empty[int]().String(), "Optional.empty")
}

func TestString_NotEmptySimple(t *testing.T) {
	v := 123
	require.EqualValues(t, Of(v).String(), spew.Sprintf("Optional[%#+v]", v))
}

func TestString_NotEmptyComposite(t *testing.T) {
	v := struct {
		X string
		Y int
		Z []interface{}
	}{
		X: "abc",
		Y: 123,
		Z: []interface{}{"abc", 123, []string{}, nil},
	}
	require.EqualValues(t, Of(v).String(), spew.Sprintf("Optional[%#+v]", v))
}

func TestEquals_BothEmpty(t *testing.T) {
	opt1 := Empty[string]()
	opt2 := Empty[string]()
	require.True(t, opt1.Equals(opt2))
}

func TestEquals_FirstEmpty(t *testing.T) {
	opt1 := Empty[string]()
	opt2 := Of("abc")
	require.False(t, opt1.Equals(opt2))
}

func TestEquals_SecondEmpty(t *testing.T) {
	opt1 := Of("abc")
	opt2 := Empty[string]()
	require.False(t, opt1.Equals(opt2))
}

func TestEquals_EqualValues(t *testing.T) {
	opt1 := Of("abc")
	opt2 := Of("abc")
	require.True(t, opt1.Equals(opt2))
}

func TestEquals_NotEqualValues(t *testing.T) {
	opt1 := Of("abc")
	opt2 := Of("abcd")
	require.False(t, opt1.Equals(opt2))
}

func TestEquals_EqualValuesComposite(t *testing.T) {
	v := struct {
		X string
		Y int
		Z []interface{}
	}{
		X: "abc",
		Y: 123,
		Z: []interface{}{"abc", 123, []string{}, nil},
	}
	opt1 := Of(v)
	opt2 := Of(v)
	require.True(t, opt1.Equals(opt2))
}

func TestEquals_NotEqualValuesComposite(t *testing.T) {
	v := struct {
		X string
		Y int
		Z []interface{}
	}{
		X: "abc",
		Y: 123,
		Z: []interface{}{"abc", 123, []string{}, nil},
	}
	v2 := struct {
		X string
		Y int
		Z []interface{}
	}{
		X: "abcd",
		Y: 1234,
		Z: []interface{}{"abc", 321, []string{"a", "b"}, nil},
	}
	opt1 := Of(v)
	opt2 := Of(v2)
	require.False(t, opt1.Equals(opt2))
}

func TestOrZero_Empty(t *testing.T) {
	assert.EqualValues(t, Empty[string]().OrZero(), "")
	assert.False(t, Empty[bool]().OrZero())
	assert.Nil(t, Empty[*string]().OrZero())
	assert.EqualValues(t, Empty[int]().OrZero(), 0)
	assert.Nil(t, Empty[[]string]().OrZero())
}

func TestOrZero_NotEmpty(t *testing.T) {
	assert.EqualValues(t, Of("abc").OrZero(), "abc")
	assert.True(t, Of(true).OrZero())
	s := "abc"
	assert.EqualValues(t, Of(&s).OrZero(), &s)
	assert.EqualValues(t, Of(123).OrZero(), 123)
	v := []string{"a", "b", "c"}
	assert.EqualValues(t, Of(v).OrZero(), v)
}

func TestTake_Empty(t *testing.T) {
	var opt Optional[int]
	opt2 := opt.Take()
	require.True(t, opt.IsEmpty())
	require.True(t, opt2.IsEmpty())
}

func TestTake_Nil(t *testing.T) {
	opt := OfNillable[string](nil)
	opt2 := opt.Take()
	require.True(t, opt.IsEmpty())
	require.True(t, opt2.IsEmpty())
}

func TestTake_NotEmpty(t *testing.T) {
	opt := Of(123)
	opt2 := opt.Take()

	require.Nil(t, opt)
	require.True(t, opt.IsEmpty())

	require.True(t, opt2.IsPresent())
	require.EqualValues(t, opt2.Get(), 123)
}

func TestTake_Nillable(t *testing.T) {
	v := []interface{}{"a", 123, 321, false, []string{}, nil}
	opt := OfNillable(&v)
	opt2 := opt.Take()

	require.Nil(t, opt)
	require.True(t, opt.IsEmpty())

	require.True(t, opt2.IsPresent())
	require.EqualValues(t, opt2.Get(), &v)
}

func TestReplace_Empty(t *testing.T) {
	opt := Empty[int]()
	opt2 := opt.Replace(321)

	require.EqualValues(t, opt.Get(), 321)
	require.True(t, opt2.IsEmpty())
}

func TestReplace_NotEmpty(t *testing.T) {
	meh := 123
	lfg := 69_420

	opt := Of(meh)
	opt2 := opt.Replace(lfg)

	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), lfg)

	require.True(t, opt2.IsPresent())
	require.EqualValues(t, opt2.Get(), meh)
}
