package test

import (
	"konbini/server/utils"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPasswordHash(t *testing.T) {
	password := "password"
	hashed, err := utils.GeneratePasswordHash(password)
	require.NoError(t, err)
	require.Regexp(t, regexp.MustCompile(`^\$argon2id\$v=(\d+)\$m=(\d+),t=(\d+),p=(\d+)\$([A-Za-z0-9+/]+)\$([A-Za-z0-9+/]+)$`), hashed)

	matches, err := utils.ComparePasswordAndHash(password, hashed)
	require.NoError(t, err)
	require.True(t, matches)

	matches, err = utils.ComparePasswordAndHash("not the same", hashed)
	require.NoError(t, err)
	require.False(t, matches)
}
