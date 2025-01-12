package permission

type Permission uint64

const (
	// NoOp just means no bit set for any permission. Useful for clearing permission.
	NoOp Permission = 0b0000_0000_0000_0000_0000_0000_0000_0000
	// Read allows reading all contents of a bento
	Read Permission = 0b0000_0000_0000_0000_0000_0000_0000_0001
	// WriteName allows changing a bento's name
	WriteName Permission = 0b0000_0000_0000_0000_0000_0000_0000_0010
	// WriteIngredientName allows changing the name of an ingredient
	WriteIngredientName Permission = 0b0000_0000_0000_0000_0000_0000_0000_0100
	// WriteIngredientName allows changing the value of an ingredient
	WriteIngredientValue Permission = 0b0000_0000_0000_0000_0000_0000_0000_1000
	// DeleteIngredient allows deleting an ingredient
	DeleteIngredient Permission = 0b0000_0000_0000_0000_0000_0000_0001_0000
	// Delete allows deletion of an entire bento + ingredients
	Delete Permission = 0b0000_0000_0000_0000_0000_0000_0010_0000
	// Admin has access to all permissions + granting others access to the bento.
	// An admin cannot assign others as admins, and the owner is the only one
	// who can assign new admins and remove existing admins.
	Admin Permission = 0b0000_0000_0000_0000_0000_0000_0100_0000
	// AddGroup is a special permission that allows the user to add the bento
	// to a group. By adding a bento to a group, it essentially grants access
	// to the bento to all users in the group. One required condition is that
	// the user needs to have this permission, be bento admin/owner and
	// be group admin/owner to be able to add this bento to a group.
	// This is to prevent non-admin/owner to be able to just add a bento to a group
	// just by having this permission enabled. This check doesn't apply to bento owner.
	AddGroup Permission = 0b0000_0000_0000_0000_0000_0000_1000_0000
	// Owner is the ultimate permission level.
	Owner Permission = 0b1000_0000_0000_0000_0000_0000_0000_0000

	// Group specific permissions
	// Only group owner and admins can add bentos to the group
	// Bentos owned by owner/admin can be added
	// Bentos that owner/admin has admin privilege and has AddGroup permission

	// AddUserToGroup allows the addition of new users to the group
	AddUserToGroup Permission = 0b0000_0000_0000_0000_0000_0001_0000_0000
	// DeleteUserFromGroup allows the removal of users from group
	DeleteUserFromGroup Permission = 0b0000_0000_0000_0000_0000_0010_0000_0000
	// GroupAdmin has all access and the privilege to add bentos to the group
	// given that thet condidtions are satisified. The conditions are explained
	// in AddGroup permission block.
	GroupAdmin Permission = 0b0000_0000_0000_0000_0000_0100_0000_0000
	// GroupOwner is the ultimate permission level for a group.
	GroupOwner Permission = 0b0100_0000_0000_0000_0000_0000_0000_0000
)
