package goptional

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	opt := Empty[interface{}]()
	require.Nil(t, opt.wrappedValue)
}

func TestOf_Interface(t *testing.T) {
	var iFace interface{} = []string{"goptional"}
	optIFace := Of(iFace)
	require.NotNil(t, optIFace.wrappedValue)
	require.EqualValues(t, optIFace.wrappedValue.value, iFace)
}

func TestOf_Channel(t *testing.T) {
	ch := make(chan int)
	optCh := Of(ch)
	require.NotNil(t, optCh.wrappedValue)
	require.EqualValues(t, optCh.wrappedValue.value, ch)
}

func TestOf_Function(t *testing.T) {
	optFn := Of(func(x bool) string {
		return "goptional"
	})
	require.NotNil(t, optFn.wrappedValue)
}

func TestOf_Pointer(t *testing.T) {
	str := "goptional"
	optPtr := Of(&str)
	require.NotNil(t, optPtr.wrappedValue)
	require.EqualValues(t, optPtr.wrappedValue.value, &str)
}

func TestOf_Slice(t *testing.T) {
	sl := []interface{}{"a", true, map[string]interface{}{
		"k1": []string{"a", "b", "c"},
		"k2": false,
	}}
	optSl := Of(sl)
	require.NotNil(t, optSl.wrappedValue)
	require.EqualValues(t, optSl.wrappedValue.value, sl)
}

func TestOf_Map(t *testing.T) {
	m := map[string]interface{}{
		"k1": []string{"a", "b", "c"},
		"k2": false,
	}
	optMap := Of(m)
	require.NotNil(t, optMap.wrappedValue)
	require.EqualValues(t, optMap.wrappedValue.value, m)
}

func TestOf_Array(t *testing.T) {
	arr := [3]interface{}{1, false, "goptional"}
	optArr := Of(arr)
	require.NotNil(t, optArr.wrappedValue)
	require.EqualValues(t, optArr.wrappedValue.value, arr)
}

func TestOf_Struct(t *testing.T) {
	st := struct {
		a string
		b uint
		c bool
		d []interface{}
	}{a: "a", b: 2, c: false}
	optStruct := Of(st)
	require.NotNil(t, optStruct.wrappedValue)
	require.EqualValues(t, optStruct.wrappedValue.value, st)
}

func TestOf_Char(t *testing.T) {
	optChar := Of('a')
	require.NotNil(t, optChar.wrappedValue)
	require.EqualValues(t, optChar.wrappedValue.value, 'a')
}

func TestOf_Float(t *testing.T) {
	optFloat := Of(1.234)
	require.NotNil(t, optFloat.wrappedValue)
	require.EqualValues(t, optFloat.wrappedValue.value, 1.234)
}

func TestOf_Bool(t *testing.T) {
	optBool := Of(true)
	require.NotNil(t, optBool.wrappedValue)
	require.EqualValues(t, optBool.wrappedValue.value, true)
}

func TestOf_Int(t *testing.T) {
	optInt := Of(1)
	require.NotNil(t, optInt.wrappedValue)
	require.EqualValues(t, optInt.wrappedValue.value, 1)
}

func TestOf_String(t *testing.T) {
	optStr := Of("goptional")
	require.NotNil(t, optStr.wrappedValue)
	require.EqualValues(t, optStr.wrappedValue.value, "goptional")
}

func TestOf_ZeroSlice(t *testing.T) {
	optSl := Of([]interface{}{})
	require.NotNil(t, optSl.wrappedValue)
	require.EqualValues(t, optSl.wrappedValue.value, []interface{}{})
}

func TestOf_ZeroMap(t *testing.T) {
	optMap := Of(map[string]interface{}{})
	require.NotNil(t, optMap.wrappedValue)
	require.EqualValues(t, optMap.wrappedValue.value, map[string]interface{}{})
}

func TestOf_ZeroArray(t *testing.T) {
	var arr [3]interface{}
	optArr := Of(arr)
	require.NotNil(t, optArr.wrappedValue)
	require.EqualValues(t, optArr.wrappedValue.value, arr)
}

func TestOf_ZeroStruct(t *testing.T) {
	st := struct {
		a string
		b uint
		c bool
		d []interface{}
	}{}
	optStruct := Of(st)
	require.NotNil(t, optStruct.wrappedValue)
	require.EqualValues(t, optStruct.wrappedValue.value, st)
}

func TestOf_ZeroChar(t *testing.T) {
	optChar := Of(' ')
	require.NotNil(t, optChar.wrappedValue)
	require.EqualValues(t, optChar.wrappedValue.value, ' ')
}

func TestOf_ZeroFloat(t *testing.T) {
	optFloat := Of(0.)
	require.NotNil(t, optFloat.wrappedValue)
	require.EqualValues(t, optFloat.wrappedValue.value, 0.)
}

func TestOf_ZeroBool(t *testing.T) {
	optBool := Of(false)
	require.NotNil(t, optBool.wrappedValue)
	require.EqualValues(t, optBool.wrappedValue.value, false)
}

func TestOf_ZeroInt(t *testing.T) {
	optInt := Of(0)
	require.NotNil(t, optInt.wrappedValue)
	require.EqualValues(t, optInt.wrappedValue.value, 0)
}

func TestOf_ZeroString(t *testing.T) {
	optStr := Of("")
	require.NotNil(t, optStr.wrappedValue)
	require.EqualValues(t, optStr.wrappedValue.value, "")
}

func TestOf_NilFunction(t *testing.T) {
	optFn := Of[func(int) bool](nil)
	require.Nil(t, optFn.wrappedValue)
}

func TestOf_NilChannel(t *testing.T) {
	optCh := Of[chan int](nil)
	require.Nil(t, optCh.wrappedValue)
}

func TestOf_NilMap(t *testing.T) {
	optMap := Of[map[string]interface{}](nil)
	require.Nil(t, optMap.wrappedValue)
}

func TestOf_NilSlice(t *testing.T) {
	optSlice := Of[[]string](nil)
	require.Nil(t, optSlice.wrappedValue)
}

func TestOf_NilInterface(t *testing.T) {
	optIFace := Of[interface{}](nil)
	require.Nil(t, optIFace.wrappedValue)
}

func TestOf_NilPointer(t *testing.T) {
	optPtr := Of[*string](nil)
	require.Nil(t, optPtr.wrappedValue)
}

func TestIsPresent_Empty(t *testing.T) {
	opt := Empty[string]()
	require.False(t, opt.IsPresent())
}

func TestIsPresent_NilValue(t *testing.T) {
	opt := Of[map[string]interface{}](nil)
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
	opt := Of[map[string]interface{}](nil)
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
		require.NotNil(t, recover())
	}()
	opt := Empty[string]()
	_ = opt.Get()
}

func TestGet_NilValue(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	opt := Of[*string](nil)
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
	Of[[]string](nil).IfPresent(func(_ []string) { called = true })
	require.False(t, called)
}

func TestIfPresent_NilActionOnEmpty(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Of[[]string](nil).IfPresent(nil)
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
	Of[*string](nil).IfPresentOrElse(func(_ *string) { actionCalled = true }, func() { emptyActionCalled = true })
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
	Of[*string](nil).IfPresentOrElse(func(_ *string) {}, nil)
}

func TestFilter_Empty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Filter(func(_ string) bool { return true })
	require.True(t, opt.IsEmpty())
}

func TestFilter_NilValue(t *testing.T) {
	opt := Of[[]string](nil)
	opt = opt.Filter(func(_ []string) bool { return true })
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
	opt := Of[*string](nil)
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
		require.NotNil(t, recover())
	}()
	Map(nil, func(i int) string { return "goptional" })
}

func TestMap_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
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
		require.NotNil(t, recover())
	}()
	MapOr(nil, func(i int) string { return "goptional" }, "default")
}

func TestMapOr_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
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
		require.NotNil(t, recover())
	}()
	MapOrElse(nil, func(i int) string { return "goptional" }, func() string { return "default" })
}

func TestMapOrElse_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
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
	opt := FlatMap(Empty[string](), func(x string) *Optional[int] { return Of(123) })
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_NilMapperOnEmpty(t *testing.T) {
	opt := FlatMap[string, interface{}](Empty[string](), nil)
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_MapToNotEmptyOnNotEmpty(t *testing.T) {
	opt := FlatMap(Of(123), func(x int) *Optional[string] { return Of(fmt.Sprintf("%v", x)) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestFlatMap_MapToEmptyOnNotEmpty(t *testing.T) {
	opt := FlatMap(Of(123), func(x int) *Optional[string] { return Empty[string]() })
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
		require.NotNil(t, recover())
	}()
	FlatMap(nil, func(x int) *Optional[string] { return Of("123") })
}

func TestFlatMap_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	FlatMap[bool, bool](nil, nil)
}

func TestAnd_Empty(t *testing.T) {
	require.True(t, Empty[string]().And(func() *Optional[string] { return Of("123") }).IsEmpty())
}

func TestAnd_NilSupplierOnEmpty(t *testing.T) {
	require.True(t, Empty[string]().And(nil).IsEmpty())
}

func TestAnd_SuppliedEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.And(func() *Optional[int] { return Empty[int]() })
	require.True(t, opt.IsEmpty())
}

func TestAnd_SuppliedNotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.And(func() *Optional[int] { return Of(321) })
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
	opt = opt.Or(func() *Optional[int] { return Of(321) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), 123)
}

func TestOr_SuppliedNotEmptyOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Or(func() *Optional[string] { return Of("123") })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Get(), "123")
}

func TestOr_SuppliedEmptyOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Or(func() *Optional[string] { return Empty[string]() })
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
		require.NotNil(t, recover())
	}()
	Empty[string]().Xor(nil)
}

func TestXor_NilOptOnNotEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
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
