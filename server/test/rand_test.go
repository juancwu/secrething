package test

import (
	"github.com/juancwu/konbini/server/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomDigits(t *testing.T) {
	t.Run("sequence of 6", func(t *testing.T) {
		seq, err := utils.RandomDigits(6)
		require.NoError(t, err)
		require.Equal(t, 6, len(seq))
	})
}
