package permissions

// PermissionBit represents the underlying type for all permissions
type PermissionBit uint64

// Permission is a common interface for all permission types
type Permission interface {
	String() string
}
