package goptional

import (
	"errors"
	"fmt"
	"strings"
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
	require.True(t, Of[interface{}](nil).IsEmpty())
}

func TestOf_NilValue(t *testing.T) {
	require.True(t, Of[*string](nil).IsEmpty())
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

func TestUnwrap_NotEmpty(t *testing.T) {
	s := "goptional"
	opt := Of(s)
	require.EqualValues(t, opt.Unwrap(), s)
}

func TestUnwrap_Empty(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.ErrorIs(t, r.(error), ErrNoValue)
	}()
	opt := Empty[string]()
	_ = opt.Unwrap()
}

func TestUnwrap_NilValue(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		require.ErrorIs(t, r.(error), ErrNoValue)
	}()
	opt := Of[*string](nil)
	_ = opt.Unwrap()
}

func TestIfPresent_NotEmpty(t *testing.T) {
	optVal := 0
	Of(123).IfPresent(func(x *int) { optVal = *x })
	require.EqualValues(t, optVal, 123)
}

func TestIfPresent_Empty(t *testing.T) {
	called := false
	Empty[int]().IfPresent(func(_ *int) { called = true })
	require.False(t, called)
}

func TestIfPresent_NilValue(t *testing.T) {
	called := false
	Of[[]string](nil).IfPresent(func(_ *[]string) { called = true })
	require.False(t, called)
}

func TestIfPresent_NilActionOnEmpty(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Of[[]string](nil).IfPresent(nil)
}

func TestIfPresent_NilActionOnNotEmpty(t *testing.T) {
	Of([]string{"a", "b", "c"}).IfPresent(nil)
}

func TestIfPresentOrElse_Empty(t *testing.T) {
	var actionCalled, emptyActionCalled bool
	Empty[string]().IfPresentOrElse(func(_ *string) { actionCalled = true }, func() { emptyActionCalled = true })
	require.False(t, actionCalled)
	require.True(t, emptyActionCalled)
}

func TestIfPresentOrElse_NilValue(t *testing.T) {
	var actionCalled, emptyActionCalled bool
	Of[*string](nil).IfPresentOrElse(func(_ **string) { actionCalled = true }, func() { emptyActionCalled = true })
	require.False(t, actionCalled)
	require.True(t, emptyActionCalled)
}

func TestIfPresentOrElse_NotEmpty(t *testing.T) {
	var actionCalled, emptyActionCalled bool
	Of(123).IfPresentOrElse(func(_ *int) { actionCalled = true }, func() { emptyActionCalled = true })
	require.True(t, actionCalled)
	require.False(t, emptyActionCalled)
}

func TestIfPresentOrElse_NilActionOnNotEmpty(t *testing.T) {
	Of(123).IfPresentOrElse(nil, func() {})
}

func TestIfPresentOrElse_NilEmptyActionOnEmpty(t *testing.T) {
	Empty[string]().IfPresentOrElse(func(_ *string) {}, nil)
}

func TestIfPresentOrElse_NilEmptyActionOnNilValue(t *testing.T) {
	Of[*string](nil).IfPresentOrElse(func(_ **string) {}, nil)
}

func TestFilter_Empty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Filter(func(_ *string) bool { return true })
	require.True(t, opt.IsEmpty())
}

func TestFilter_NilValue(t *testing.T) {
	opt := Of[*[]string](nil)
	opt = opt.Filter(func(_ **[]string) bool { return true })
	require.True(t, opt.IsEmpty())
}

func TestFilter_NotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Filter(func(_ *int) bool { return true })
	require.True(t, opt.IsPresent())
}

func TestFilter_NilPredicateOnEmpty(t *testing.T) {
	require.True(t, Empty[string]().Filter(nil).IsEmpty())
}

func TestFilter_NilPredicateOnNotEmpty(t *testing.T) {
	opt := Of(123).Filter(nil)
	require.True(t, opt.IsEmpty())
}

func TestFilter_PredicateNotOkOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Filter(func(_ *string) bool { return false })
	require.True(t, opt.IsEmpty())
}

func TestFilter_PredicateNotOkOnNilValue(t *testing.T) {
	opt := Of[*string](nil)
	opt = opt.Filter(func(_ **string) bool { return false })
	require.True(t, opt.IsEmpty())
}

func TestFilter_PredicateNotOkOnNotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Filter(func(_ *int) bool { return false })
	require.True(t, opt.IsEmpty())
}

func TestMap_Empty(t *testing.T) {
	opt := Map(Empty[string](), func(s *string) string { return *s })
	require.True(t, opt.IsEmpty())
}

func TestMap_NilMapperOnEmpty(t *testing.T) {
	opt := Map[string, interface{}](Empty[string](), nil)
	require.True(t, opt.IsEmpty())
}

func TestMap_NotEmpty(t *testing.T) {
	opt := Map(Of(123), func(x *int) string { return fmt.Sprintf("%v", *x) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "123")
}

func TestMap_NilMapperOnNotEmpty(t *testing.T) {
	require.True(t, Map[int, string](Of(123), nil).IsEmpty())
}

func TestMap_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Map(nil, func(_ *int) string { return "goptional" })
}

func TestMap_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	Map[bool, bool](nil, nil)
}

func TestMapOr_Empty(t *testing.T) {
	opt := MapOr(Empty[string](), func(s *string) string { return *s }, "default")
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "default")
}

func TestMapOr_NilMapperOnEmpty(t *testing.T) {
	opt := MapOr[string, interface{}](Empty[string](), nil, "default")
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "default")
}

func TestMapOr_NotEmpty(t *testing.T) {
	opt := MapOr(Of(123), func(x *int) string { return fmt.Sprintf("%v", *x) }, "default")
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "123")
}

func TestMapOr_NilMapperOnNotEmpty(t *testing.T) {
	require.True(t, MapOr(Of(123), nil, "default").IsEmpty())
}

func TestMapOr_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOr(nil, func(_ *int) string { return "goptional" }, "default")
}

func TestMapOr_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOr[bool](nil, nil, "default")
}

func TestMapOrElse_Empty(t *testing.T) {
	opt := MapOrElse(Empty[string](), func(s *string) string { return *s }, func() string { return "default" })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "default")
}

func TestMapOrElse_NilMapperOnEmpty(t *testing.T) {
	opt := MapOrElse(Empty[string](), nil, func() string { return "default" })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "default")
}

func TestMapOrElse_NotEmpty(t *testing.T) {
	opt := MapOrElse(Of(123), func(x *int) string { return fmt.Sprintf("%v", *x) }, func() string { return "default" })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "123")
}

func TestMapOrElse_NilMapperOnNotEmpty(t *testing.T) {
	require.True(t, MapOrElse(Of(123), nil, func() string { return "default" }).IsEmpty())
}

func TestMapOrElse_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOrElse(nil, func(_ *int) string { return "goptional" }, func() string { return "default" })
}

func TestMapOrElse_NilMapperOnNilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	MapOrElse[bool](nil, nil, func() string { return "default" })
}

func TestMapOrElse_NilSupplierOnEmpty(t *testing.T) {
	require.True(t, MapOrElse(Empty[string](), func(_ *string) int { return 0 }, nil).IsEmpty())
}

func TestFlatMap_Empty(t *testing.T) {
	opt := FlatMap(Empty[string](), func(_ *string) Optional[int] { return Of(123) })
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_NilMapperOnEmpty(t *testing.T) {
	opt := FlatMap[string, interface{}](Empty[string](), nil)
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_MapToNotEmptyOnNotEmpty(t *testing.T) {
	opt := FlatMap(Of(123), func(x *int) Optional[string] { return Of(fmt.Sprintf("%v", *x)) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "123")
}

func TestFlatMap_MapToEmptyOnNotEmpty(t *testing.T) {
	opt := FlatMap(Of(123), func(_ *int) Optional[string] { return Empty[string]() })
	require.True(t, opt.IsEmpty())
}

func TestFlatMap_NilMapperOnNotEmpty(t *testing.T) {
	require.True(t, FlatMap[int, string](Of(123), nil).IsEmpty())
}

func TestFlatMap_NilInput(t *testing.T) {
	defer func() {
		require.Nil(t, recover())
	}()
	FlatMap(nil, func(_ *int) Optional[string] { return Of("123") })
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
	require.EqualValues(t, opt.Unwrap(), 321)
}

func TestAnd_NilSupplierOnNotEmpty(t *testing.T) {
	require.True(t, Of(123).And(nil).IsEmpty())
}

func TestOr_NilSupplierOnNotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Or(nil)
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), 123)
}

func TestOr_NotEmpty(t *testing.T) {
	opt := Of(123)
	opt = opt.Or(func() Optional[int] { return Of(321) })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), 123)
}

func TestOr_SuppliedNotEmptyOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Or(func() Optional[string] { return Of("123") })
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), "123")
}

func TestOr_SuppliedEmptyOnEmpty(t *testing.T) {
	opt := Empty[string]()
	opt = opt.Or(func() Optional[string] { return Empty[string]() })
	require.True(t, opt.IsEmpty())
}

func TestOr_NilSupplierOnEmpty(t *testing.T) {
	require.True(t, Empty[string]().Or(nil).IsEmpty())
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
	require.Empty(t, Empty[string]().OrElseGet(nil), "")
}

func TestOrPanicWith_NotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrPanicWith(func() error { return errors.New("woops") }), 123)
}

func TestOrPanicWith_NilSupplierOnNotEmpty(t *testing.T) {
	require.EqualValues(t, Of(123).OrPanicWith(nil), 123)
}

func TestOrPanicWith_Empty(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		require.Error(t, err)
		require.EqualError(t, err, "woops")
	}()
	Empty[string]().OrPanicWith(func() error { return errors.New("woops") })
}

func TestOrPanicWith_SuppliedNilOnEmpty(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		err, ok := r.(error)
		require.True(t, ok)
		require.ErrorIs(t, err, ErrNoValue)
	}()
	Empty[string]().OrPanicWith(func() error { return nil })
}

func TestOrPanicWith_NilSupplierOnEmpty(t *testing.T) {
	defer func() {
		require.NotNil(t, recover())
	}()
	Empty[string]().OrPanicWith(nil)
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
	require.EqualValues(t, opt.Unwrap(), 321)
}

func TestXor_SecondEmpty(t *testing.T) {
	opt := Of(123).Xor(Empty[int]())
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), 123)
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

func TestOrDefault_Empty(t *testing.T) {
	assert.EqualValues(t, Empty[string]().OrDefault(), "")
	assert.False(t, Empty[bool]().OrDefault())
	assert.Nil(t, Empty[*string]().OrDefault())
	assert.EqualValues(t, Empty[int]().OrDefault(), 0)
	assert.Nil(t, Empty[[]string]().OrDefault())
}

func TestOrDefault_NotEmpty(t *testing.T) {
	assert.EqualValues(t, Of("abc").OrDefault(), "abc")
	assert.True(t, Of(true).OrDefault())
	s := "abc"
	assert.EqualValues(t, Of(&s).OrDefault(), &s)
	assert.EqualValues(t, Of(123).OrDefault(), 123)
	v := []string{"a", "b", "c"}
	assert.EqualValues(t, Of(v).OrDefault(), v)
}

func TestTake_Empty(t *testing.T) {
	var opt Optional[int]
	opt2 := opt.Take()
	require.True(t, opt.IsEmpty())
	require.True(t, opt2.IsEmpty())
}

func TestTake_Nil(t *testing.T) {
	opt := Of[*string](nil)
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
	require.EqualValues(t, opt2.Unwrap(), 123)
}

func TestTake_Ptr(t *testing.T) {
	v := []interface{}{"a", 123, 321, false, []string{}, nil}
	opt := Of(&v)
	opt2 := opt.Take()

	require.Nil(t, opt)
	require.True(t, opt.IsEmpty())

	require.True(t, opt2.IsPresent())
	require.EqualValues(t, opt2.Unwrap(), &v)
}

func TestReplace_Empty(t *testing.T) {
	opt := Empty[int]()
	opt2 := opt.Replace(321)

	require.EqualValues(t, opt.Unwrap(), 321)
	require.True(t, opt2.IsEmpty())
}

func TestReplace_NotEmpty(t *testing.T) {
	meh := 123
	lfg := 69_420

	opt := Of(meh)
	opt2 := opt.Replace(lfg)

	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), lfg)

	require.True(t, opt2.IsPresent())
	require.EqualValues(t, opt2.Unwrap(), meh)
}

type sampleStruct struct {
	X string   `json:"x"`
	Y bool     `json:"y"`
	Z []string `json:"z"`
}

var sampleStructInst = &sampleStruct{
	X: "gmgn",
	Y: true,
	Z: []string{"a", "b", "c"},
}

const sampleJSON = `{"x":"gmgn","y":true,"z":["a","b","c"]}`

func TestMarshalJSON_Empty(t *testing.T) {
	jsonBytes, err := Empty[int]().MarshalJSON()
	require.NoError(t, err)
	require.EqualValues(t, jsonBytes, nilAsJSON)
}

func TestMarshalJSON_NotEmpty(t *testing.T) {
	jsonBytes, err := Of("gmgn").MarshalJSON()
	require.NoError(t, err)
	require.EqualValues(t, jsonBytes, []byte("\"gmgn\""))

	jsonBytes, err = Of(true).MarshalJSON()
	require.NoError(t, err)
	require.EqualValues(t, jsonBytes, []byte("true"))

	jsonBytes, err = Of(sampleStructInst).MarshalJSON()
	require.NoError(t, err)
	require.EqualValues(t, jsonBytes, []byte(sampleJSON))
}

func TestUnmarshalJSON_NoDataOnEmpty(t *testing.T) {
	opt := Empty[int]()
	err := opt.UnmarshalJSON(nil)
	require.NoError(t, err)
	require.True(t, opt.IsEmpty())
}

func TestUnmarshalJSON_NullDataOnEmpty(t *testing.T) {
	opt := Empty[int]()
	err := opt.UnmarshalJSON(nilAsJSON)
	require.NoError(t, err)
	require.True(t, opt.IsEmpty())
}

func TestUnmarshalJSON_NoDataOnNotEmpty(t *testing.T) {
	opt := Of(123)
	err := opt.UnmarshalJSON(nil)
	require.NoError(t, err)
	require.True(t, opt.IsEmpty())
}

func TestUnmarshalJSON_NullDataOnNotEmpty(t *testing.T) {
	opt := Of(123)
	err := opt.UnmarshalJSON(nilAsJSON)
	require.NoError(t, err)
	require.True(t, opt.IsEmpty())
}

func TestUnmarshalJSON_InvalidData(t *testing.T) {
	opt := Empty[*sampleStruct]()
	err := opt.UnmarshalJSON([]byte(sampleJSON)[1:])
	require.Error(t, err)
	require.True(t, opt.IsEmpty())
}

func TestUnmarshalJSON_ValidDataOnEmpty(t *testing.T) {
	opt := Empty[*sampleStruct]()
	err := opt.UnmarshalJSON([]byte(sampleJSON))
	require.NoError(t, err)
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), sampleStructInst)
}

func TestUnmarshalJSON_ValidDataOnNotEmpty(t *testing.T) {
	s := sampleStruct{
		X: "diff",
		Z: nil,
	}
	opt := Of(s)
	err := opt.UnmarshalJSON([]byte(sampleJSON))
	require.NoError(t, err)
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), *sampleStructInst)
}

func TestFlatten_Empty(t *testing.T) {
	require.True(t, Flatten(Empty[Optional[int]]()).IsEmpty())
}

func TestFlatten_NotEmpty(t *testing.T) {
	opt := Flatten(Of(Of(123)))
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), 123)
}

func TestZip_SomeEmpty(t *testing.T) {
	require.True(t, Zip(Empty[int](), Empty[string]()).IsEmpty())
	require.True(t, Zip(Of(123), Empty[string]()).IsEmpty())
	require.True(t, Zip(Empty[int](), Of("gm")).IsEmpty())
}

func TestZip_BothNotEmpty(t *testing.T) {
	opt := Zip(Of(123), Of("gm"))
	require.True(t, opt.IsPresent())

	v := opt.Unwrap()
	require.EqualValues(t, v.First, 123)
	require.EqualValues(t, v.Second, "gm")
}

func TestUnzip_Empty(t *testing.T) {
	o1, o2 := Unzip(Empty[*Pair[Optional[int], Optional[string]]]())
	require.True(t, o1.IsEmpty())
	require.True(t, o2.IsEmpty())
}

func TestUnzip_BothNotEmpty(t *testing.T) {
	pair := &Pair[Optional[int], Optional[string]]{First: Of(123), Second: Of("gm")}
	o1, o2 := Unzip(Of(pair))

	require.True(t, o1.IsPresent())
	require.EqualValues(t, o1, pair.First)

	require.True(t, o2.IsPresent())
	require.EqualValues(t, o2, pair.Second)

}

func TestUnzip_LeftEmpty(t *testing.T) {
	pair := &Pair[Optional[int], Optional[string]]{First: Empty[int](), Second: Of("gm")}
	o1, o2 := Unzip(Of(pair))

	require.True(t, o1.IsEmpty())
	require.EqualValues(t, o1, pair.First)

	require.True(t, o2.IsPresent())
	require.EqualValues(t, o2, pair.Second)
}

func TestUnzip_RightEmpty(t *testing.T) {
	pair := &Pair[Optional[int], Optional[string]]{First: Of(123), Second: Empty[string]()}
	o1, o2 := Unzip(Of(pair))

	require.True(t, o1.IsPresent())
	require.EqualValues(t, o1, pair.First)

	require.True(t, o2.IsEmpty())
	require.EqualValues(t, o2, pair.Second)
}

func TestUnzip_BothEmpty(t *testing.T) {
	pair := &Pair[Optional[int], Optional[string]]{First: Empty[int](), Second: Empty[string]()}
	o1, o2 := Unzip(Of(pair))

	require.True(t, o1.IsEmpty())
	require.EqualValues(t, o1, pair.First)

	require.True(t, o2.IsEmpty())
	require.EqualValues(t, o2, pair.Second)
}

func TestZipWith_SomeEmpty(t *testing.T) {
	require.True(t, ZipWith[string, int, interface{}](Empty[string](), Empty[int](), nil).IsEmpty())
	require.True(t, ZipWith[string, int, interface{}](Of("gm"), Empty[int](), nil).IsEmpty())
	require.True(t, ZipWith[string, int, interface{}](Empty[string](), Of(123), nil).IsEmpty())
}

func TestZipWith_NilMapperOnNotEmpty(t *testing.T) {
	require.True(t, ZipWith[string, int, interface{}](Of("gm"), Of(123), nil).IsEmpty())
}

func TestZipWith_BothNotEmpty(t *testing.T) {
	opt := ZipWith(Of("gm"), Of([]int{1, 2, 3, 4}), func(x *string, y *[]int) []interface{} {
		return []interface{}{*x, *y}
	})
	require.True(t, opt.IsPresent())
	require.EqualValues(t, opt.Unwrap(), []interface{}{"gm", []int{1, 2, 3, 4}})
}

func TestZipWith_BothNotEmptyWithNilReturn(t *testing.T) {
	opt := ZipWith(Of("gm"), Of([]int{1, 2, 3, 4}), func(x *string, y *[]int) []interface{} {
		return nil
	})
	require.True(t, opt.IsEmpty())
}

func TestIs_Empty(t *testing.T) {
	require.False(t, Empty[int]().Is(nil))
	require.False(t, Empty[int]().Is(func(_ *int) bool { return true }))
}

func TestIs_NilPredicateOnNotEmpty(t *testing.T) {
	require.False(t, Of(123).Is(nil))
}

func TestIs_NotEmpty(t *testing.T) {
	require.True(t, Of(123).Is(func(x *int) bool { return *x%2 != 0 }))
	require.True(t, Of(1234).Is(func(x *int) bool { return *x > 100 }))
	require.True(t, Of([]string{"gm", "Gn"}).Is(func(x *[]string) bool { return strings.ToLower((*x)[1]) == "gn" }))
}

func TestVal_NotEmpty(t *testing.T) {
	v, err := Of(123).Val()
	require.EqualValues(t, v, 123)
	require.NoError(t, err)
}

func TestVal_Empty(t *testing.T) {
	v, err := Empty[string]().Val()
	require.EqualValues(t, v, "")
	require.ErrorIs(t, err, ErrNoValue)
}

func TestValOr_NotEmpty(t *testing.T) {
	v, err := Of(123).ValOr(nil)
	require.EqualValues(t, v, 123)
	require.NoError(t, err)

	v2, err := Of(321).ValOr(errors.New("woops"))
	require.EqualValues(t, v2, 321)
	require.NoError(t, err)
}

func TestValOr_Empty(t *testing.T) {
	inErr := errors.New("woops")
	v, err := Empty[string]().ValOr(inErr)
	require.EqualValues(t, v, "")
	require.ErrorIs(t, err, inErr)
}

func TestValOr_NilErrOnEmpty(t *testing.T) {
	_, err := Empty[string]().ValOr(nil)
	require.ErrorIs(t, err, ErrNoValue)
}

func TestValOrElse_NotEmpty(t *testing.T) {
	v, err := Of(123).ValOrElse(nil)
	require.EqualValues(t, v, 123)
	require.NoError(t, err)

	v2, err := Of(321).ValOrElse(func() error { return errors.New("woops") })
	require.EqualValues(t, v2, 321)
	require.NoError(t, err)
}

func TestValOrElse_Empty(t *testing.T) {
	inErr := errors.New("woops")
	v, err := Empty[string]().ValOrElse(func() error { return inErr })
	require.EqualValues(t, v, "")
	require.ErrorIs(t, err, inErr)
}

func TestValOrElse_NilSupplierOnEmpty(t *testing.T) {
	v, err := Empty[string]().ValOrElse(nil)
	require.EqualValues(t, v, "")
	require.ErrorIs(t, err, ErrNoValue)
}

func TestValOrElse_SuppliedNilOnEmpty(t *testing.T) {
	v, err := Empty[string]().ValOrElse(func() error { return nil })
	require.EqualValues(t, v, "")
	require.ErrorIs(t, err, ErrNoValue)
}
