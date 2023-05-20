package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupportedCurrency(t *testing.T) {
	curr := CAD

	require.True(t, IsSupportedCurrency(curr))

	invalidCurr := "hi"
	require.False(t, IsSupportedCurrency(invalidCurr))
}
