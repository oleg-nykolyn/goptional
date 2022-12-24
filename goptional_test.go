package goptional

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDoNothing_AlwaysValid(t *testing.T) {
	DoNothing()
	require.True(t, true)
}
