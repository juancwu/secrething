package permission

import (
	"encoding/binary"
	"errors"
)

const (
	// NoOp just means no bit set for any permission. Useful for clearing permission.
	NoOp uint64 = 0
	// Read allows reading all contents of a bento
	Read uint64 = 1
	// WriteName allows changing a bento's name
	WriteName uint64 = 1 << 1
	// WriteIngredientName allows changing the name of an ingredient
	WriteIngredientName uint64 = 1 << 2
	// WriteIngredientName allows changing the value of an ingredient
	WriteIngredientValue uint64 = 1 << 3
	// DeleteIngredient allows deleting an ingredient
	DeleteIngredient uint64 = 1 << 4
	// Delete allows deletion of an entire bento + ingredients
	Delete uint64 = 1 << 5
	// Admin has access to all permissions + granting others access to the bento.
	// An admin cannot assign others as admins, and the owner is the only one
	// who can assign new admins and remove existing admins.
	Admin uint64 = 1 << 6
	// AddGroup is a special permission that allows the user to add the bento
	// to a group. By adding a bento to a group, it essentially grants access
	// to the bento to all users in the group. One required condition is that
	// the user needs to have this permission, be bento admin/owner and
	// be group admin/owner to be able to add this bento to a group.
	// This is to prevent non-admin/owner to be able to just add a bento to a group
	// just by having this permission enabled. This check doesn't apply to bento owner.
	AddGroup uint64 = 1 << 7
	// Owner is the ultimate permission level.
	Owner uint64 = 1 << 63

	// Group specific permissions
	// Only group owner and admins can add bentos to the group
	// Bentos owned by owner/admin can be added
	// Bentos that owner/admin has admin privilege and has AddGroup permission

	// AddUserToGroup allows the addition of new users to the group
	AddUserToGroup uint64 = 1 << 8
	// DeleteUserFromGroup allows the removal of users from group
	DeleteUserFromGroup uint64 = 1 << 9
	// GroupAdmin has all access and the privilege to add bentos to the group
	// given that thet condidtions are satisified. The conditions are explained
	// in AddGroup permission block.
	GroupAdmin uint64 = 1 << 10
	// GroupOwner is the ultimate permission level for a group.
	GroupOwner uint64 = 1 << 62
)

// GetBentoOwnerPermissions gets the combination of bits that an owner of a bento would have.
func GetBentoOwnerPermissions() uint64 {
	return Owner | WriteName | WriteIngredientName | WriteIngredientValue | Read | Delete | DeleteIngredient | AddGroup | Admin
}

// ToBytes transform a uint64 into bytes
func ToBytes(permission uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, permission)
	return b
}

// FromBytes transform bytes into uint64
func FromBytes(permission []byte) (uint64, error) {
	if permission == nil {
		return 0, errors.New("Cannot scan Nil bytes")
	}
	if len(permission) != 8 {
		return 0, errors.New("Invalid permission bytes length")
	}

	p := binary.BigEndian.Uint64(permission)

	return p, nil
}
