package permissions

// PermissionBits represents the underlying type for all permissions
type PermissionBits uint64

// Permission is a common interface for all permission types
type Permission interface {
	Bits() PermissionBits
	String() string
	All(perms ...PermissionBits) bool
	Any(perms ...PermissionBits) bool
}
