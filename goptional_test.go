package goptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	opt := Empty[interface{}]()
	require.Nil(t, opt.wrappedValue)
}

func TestOf_NonZeroValues(t *testing.T) {
	// string
	optStr := Of("goptional")
	require.NotNil(t, optStr.wrappedValue)
	require.EqualValues(t, optStr.wrappedValue.value, "goptional")

	// int
	optInt := Of(1)
	require.NotNil(t, optInt.wrappedValue)
	require.EqualValues(t, optInt.wrappedValue.value, 1)

	// bool
	optBool := Of(true)
	require.NotNil(t, optBool.wrappedValue)
	require.EqualValues(t, optBool.wrappedValue.value, true)

	// float
	optFloat := Of(1.234)
	require.NotNil(t, optFloat.wrappedValue)
	require.EqualValues(t, optFloat.wrappedValue.value, 1.234)

	// char
	optChar := Of('a')
	require.NotNil(t, optChar.wrappedValue)
	require.EqualValues(t, optChar.wrappedValue.value, 'a')

	// struct
	st := struct {
		a string
		b uint
		c bool
		d []interface{}
	}{a: "a", b: 2, c: false}
	optStruct := Of(st)
	require.NotNil(t, optStruct.wrappedValue)
	require.EqualValues(t, optStruct.wrappedValue.value, st)

	// array
	arr := [3]interface{}{1, false, "goptional"}
	optArr := Of(arr)
	require.NotNil(t, optArr.wrappedValue)
	require.EqualValues(t, optArr.wrappedValue.value, arr)

	// map
	m := map[string]interface{}{
		"k1": []string{"a", "b", "c"},
		"k2": false,
	}
	optMap := Of(m)
	require.NotNil(t, optMap.wrappedValue)
	require.EqualValues(t, optMap.wrappedValue.value, m)

	// slice
	sl := []interface{}{"a", true, m}
	optSl := Of(sl)
	require.NotNil(t, optSl.wrappedValue)
	require.EqualValues(t, optSl.wrappedValue.value, sl)

	// pointer
	str := "goptional"
	optPtr := Of(&str)
	require.NotNil(t, optPtr.wrappedValue)
	require.EqualValues(t, optPtr.wrappedValue.value, &str)

	// function
	optFn := Of(func(x bool) string {
		return str
	})
	require.NotNil(t, optFn.wrappedValue)

	// channel
	ch := make(chan int)
	optCh := Of(ch)
	require.NotNil(t, optCh.wrappedValue)
	require.EqualValues(t, optCh.wrappedValue.value, ch)

	// interface

	var iFace interface{} = []string{str}
	optIFace := Of(iFace)
	require.NotNil(t, optIFace.wrappedValue)
	require.EqualValues(t, optIFace.wrappedValue.value, iFace)
}

func TestOf_ZeroValues(t *testing.T) {
	// string
	optStr := Of("")
	require.NotNil(t, optStr.wrappedValue)
	require.EqualValues(t, optStr.wrappedValue.value, "")

	// int
	optInt := Of(0)
	require.NotNil(t, optInt.wrappedValue)
	require.EqualValues(t, optInt.wrappedValue.value, 0)

	// bool
	optBool := Of(false)
	require.NotNil(t, optBool.wrappedValue)
	require.EqualValues(t, optBool.wrappedValue.value, false)

	// float
	optFloat := Of(0.)
	require.NotNil(t, optFloat.wrappedValue)
	require.EqualValues(t, optFloat.wrappedValue.value, 0.)

	// char
	optChar := Of(' ')
	require.NotNil(t, optChar.wrappedValue)
	require.EqualValues(t, optChar.wrappedValue.value, ' ')

	// struct
	st := struct {
		a string
		b uint
		c bool
		d []interface{}
	}{}
	optStruct := Of(st)
	require.NotNil(t, optStruct.wrappedValue)
	require.EqualValues(t, optStruct.wrappedValue.value, st)

	// array
	var arr [3]interface{}
	optArr := Of(arr)
	require.NotNil(t, optArr.wrappedValue)
	require.EqualValues(t, optArr.wrappedValue.value, arr)

	// map
	optMap := Of(map[string]interface{}{})
	require.NotNil(t, optMap.wrappedValue)
	require.EqualValues(t, optMap.wrappedValue.value, map[string]interface{}{})

	// slice
	optSl := Of([]interface{}{})
	require.NotNil(t, optSl.wrappedValue)
	require.EqualValues(t, optSl.wrappedValue.value, []interface{}{})
}

func TestOf_NilValues(t *testing.T) {
	optPtr := Of[*string](nil)
	require.Nil(t, optPtr.wrappedValue)

	optIFace := Of[interface{}](nil)
	require.Nil(t, optIFace.wrappedValue)

	optSlice := Of[[]string](nil)
	require.Nil(t, optSlice.wrappedValue)

	optMap := Of[map[string]interface{}](nil)
	require.Nil(t, optMap.wrappedValue)

	optCh := Of[chan int](nil)
	require.Nil(t, optCh.wrappedValue)

	optFn := Of[func(int) bool](nil)
	require.Nil(t, optFn.wrappedValue)
}
