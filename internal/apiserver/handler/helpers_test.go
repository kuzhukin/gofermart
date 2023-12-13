package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLunaValidation(t *testing.T) {
	require.False(t, validateOrderID("4561261212345464"))
	require.True(t, validateOrderID("4561261212345467"))
}
