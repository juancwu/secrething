package store

import (
	"strconv"
	"strings"
	"testing"
)

// formats the integer into a binary string with format 0bxxxx_xxxx.
// Longest it will be is 32bits and it will pad start with 0 to complete 4 bits.
func formatBinaryWithUnderscores(n int) string {
	binStr := strconv.FormatInt(int64(n), 2)
	length := len(binStr)
	if length%4 != 0 {
		padding := 4 - length%4
		binStr = strings.Repeat("0", padding) + binStr
		length += padding
	}

	var builder strings.Builder
	// prefix
	builder.WriteString("0b")
	for i, digit := range binStr {
		if i > 0 && i%4 == 0 {
			builder.WriteRune('_')
		}
		builder.WriteRune(digit)
	}
	return builder.String()
}

// Test that all the constants are the right value
func TestBentoPermissionBinaries(t *testing.T) {
	expectedPerms := []int{
		0b0000_0000_0000_0000,
		0b0000_0000_0000_0001,
		0b0000_0000_0000_0010,
		0b0000_0000_0000_0100,
		0b0000_0000_0000_1000,
		0b0000_0000_0001_0000,
		0b0000_0000_0010_0000,
		0b0000_0000_0100_0000,
		0b0000_0000_1000_0000,
	}
	testPerms := []int{
		O_NO_PERMS,
		O_WRITE,
		O_SHARE,
		O_GRANT_SHARE,
		O_DELETE,
		O_WRITE_INGRIDIENT,
		O_DELETE_INGRIDIENT,
		O_RENAME_INGRIDIENT,
		O_RENAME_BENTO,
	}

	for i := 0; i < len(testPerms); i++ {
		if testPerms[i] != expectedPerms[i] {
			t.Logf("Expected O_NO_PERMS to be %s but received %s", formatBinaryWithUnderscores(expectedPerms[i]), formatBinaryWithUnderscores(testPerms[i]))
			t.Fail()
		}
	}
}
