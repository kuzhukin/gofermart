package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLunaValidation(t *testing.T) {
	require.False(t, validateOrderId("4561261212345464"))
	require.True(t, validateOrderId("4561261212345467"))
}
