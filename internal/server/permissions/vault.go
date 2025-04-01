package permissions

// VaultPermission is for vault-level operations
type VaultPermission PermissionBit

// EnvironmentPermission is for environment-level operations
type EnvironmentPermission PermissionBit

// Ensure permission types implement the Permission interface
var _ Permission = VaultPermission(0)
var _ Permission = EnvironmentPermission(0)

const (
	// Vault permissions
	PermissionVaultRead   VaultPermission = 0x0000_0000_0000_0001
	PermissionVaultUpdate VaultPermission = 0x0000_0000_0000_0002
	PermissionVaultDelete VaultPermission = 0x0000_0000_0000_0004
	PermissionVaultCreate VaultPermission = 0x0000_0000_0000_0008

	// Environment permissions
	PermissionEnvironmentRead   EnvironmentPermission = 0x0000_0000_0100_0000
	PermissionEnvironmentCreate EnvironmentPermission = 0x0000_0000_0200_0000
	PermissionEnvironmentUpdate EnvironmentPermission = 0x0000_0000_0400_0000
	PermissionEnvironmentDelete EnvironmentPermission = 0x0000_0000_0800_0000
)

// Convert to VaultPermission for database storage compatibility
func (p EnvironmentPermission) ToVaultPermission() VaultPermission {
	return VaultPermission(p)
}

// CanReadEnvironment checks if the user has permission to read environments
// Users with vault read permissions can read environments even if they don't have explicit environment read permissions
func (p VaultPermission) CanReadEnvironment() bool {
	return p&VaultPermission(PermissionEnvironmentRead) == VaultPermission(PermissionEnvironmentRead) ||
		p&PermissionVaultRead == PermissionVaultRead
}

// EnvironmentPermission implementation
func (p EnvironmentPermission) CanReadEnvironment() bool {
	return p&PermissionEnvironmentRead == PermissionEnvironmentRead
}

// String returns a human-readable permission string
func (p VaultPermission) String() string {
	switch p {
	case PermissionVaultRead:
		return "vault.read"
	case PermissionVaultUpdate:
		return "vault.update"
	case PermissionVaultDelete:
		return "vault.delete"
	case PermissionVaultCreate:
		return "vault.create"
	}
	return ""
}

// String returns a human-readable permission string
func (p EnvironmentPermission) String() string {
	switch p {
	case PermissionEnvironmentRead:
		return "environment.read"
	case PermissionEnvironmentCreate:
		return "environment.create"
	case PermissionEnvironmentUpdate:
		return "environment.update"
	case PermissionEnvironmentDelete:
		return "environment.delete"
	}
	return ""
}

// AllVaultPermissions checks the given user permission has all the required vault permissions.
func AllVaultPermissions(userPerm VaultPermission, perms ...VaultPermission) bool {
	for _, perm := range perms {
		if userPerm&perm == 0 {
			return false
		}
	}

	return true
}

// AnyVaultPermission checks the given user permission has any of the vault permissions.
func AnyVaultPermission(userPerm VaultPermission, perms ...VaultPermission) bool {
	for _, perm := range perms {
		if userPerm&perm == perm {
			return true
		}
	}

	return false
}

// AllEnvironmentPermissions checks the given user permission has all the required environment permissions.
func AllEnvironmentPermissions(userPerm VaultPermission, perms ...EnvironmentPermission) bool {
	for _, perm := range perms {
		if userPerm&VaultPermission(perm) == 0 {
			return false
		}
	}

	return true
}

// AnyEnvironmentPermission checks the given user permission has any of the environment permissions.
func AnyEnvironmentPermission(userPerm VaultPermission, perms ...EnvironmentPermission) bool {
	for _, perm := range perms {
		if userPerm&VaultPermission(perm) == VaultPermission(perm) {
			return true
		}
	}

	return false
}

