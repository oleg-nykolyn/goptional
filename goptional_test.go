package goptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	opt := Empty[interface{}]()
	require.NotNil(t, opt)
	require.Nil(t, opt.wrappedValue)
	require.True(t, opt.IsEmpty())
	require.EqualValues(t, opt, Empty[interface{}]())
}
