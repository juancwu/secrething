package permission

import (
	"testing"
)

func TestPermissionOperations(t *testing.T) {
	tests := []struct {
		perm     Permission
		expected uint64
	}{
		{NoOp, 0b0},
		{Read, 0b1},
		{WriteName, 0b10},
		{WriteIngredientName, 0b100},
		{WriteIngredientValue, 0b1000},
		{DeleteIngredient, 0b1_0000},
		{Delete, 0b10_0000},
		{Admin, 0b100_0000},
		{AddGroup, 0b1000_0000},
		{Owner, 0b1000_0000_0000_0000_0000_0000_0000_0000},
		{GroupOwner, 0b0100_0000_0000_0000_0000_0000_0000_0000},
	}

	for _, tt := range tests {
		if uint64(tt.perm) != tt.expected {
			t.Errorf("Permission %064b does not match expected %064b", tt.perm, tt.expected)
		}
	}
}
