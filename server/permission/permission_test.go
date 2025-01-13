package permission

import (
	"testing"
)

func TestPermissionOperations(t *testing.T) {
	tests := []struct {
		perm     uint64
		expected uint64
	}{
		{NoOp, 0},
		{Read, 1},
		{WriteName, 1 << 1},
		{WriteIngredientName, 1 << 2},
		{WriteIngredientValue, 1 << 3},
		{DeleteIngredient, 1 << 4},
		{Delete, 1 << 5},
		{Admin, 1 << 6},
		{AddGroup, 1 << 7},
		{Owner, 1 << 63},
		{GroupOwner, 1 << 62},
	}

	for _, tt := range tests {
		if uint64(tt.perm) != tt.expected {
			t.Errorf("Permission %064b does not match expected %064b", tt.perm, tt.expected)
		}
	}
}
