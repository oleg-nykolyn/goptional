package goptional

import (
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
