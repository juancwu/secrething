package permissions

import (
	"strings"
)

const (
	// Vault-level permissions
	PermissionVaultOwner  PermissionBits = 1 << 0 // 0x0001 - Can transfer ownership
	PermissionVaultAdmin  PermissionBits = 1 << 1 // 0x0002 - Can manage permissions
	PermissionVaultEditor PermissionBits = 1 << 2 // 0x0004 - Can edit vault and secrets
	PermissionVaultShare  PermissionBits = 1 << 3 // 0x0008 - Can share with others
	PermissionVaultRead   PermissionBits = 1 << 4 // 0x0010 - Can view vault metadata

	// Secret-level permissions
	PermissionSecretCreate PermissionBits = 1 << 5 // 0x0020 - Can create secrets
	PermissionSecretRead   PermissionBits = 1 << 6 // 0x0040 - Can read secrets
	PermissionSecretUpdate PermissionBits = 1 << 7 // 0x0080 - Can update secrets
	PermissionSecretDelete PermissionBits = 1 << 8 // 0x0100 - Can delete secrets

	// History and audit permissions
	PermissionViewHistory   PermissionBits = 1 << 9  // 0x0200 - Can view access history
	PermissionViewAuditLogs PermissionBits = 1 << 10 // 0x0400 - Can view audit logs

	// Team-level permissions
	PermissionTeamManageVaults PermissionBits = 1 << 11 // 0x0800 - Can create vaults owned by the team

	// Common permission sets
	PermissionsOwner PermissionBits = PermissionVaultOwner | PermissionVaultAdmin | PermissionVaultShare |
		PermissionVaultRead | PermissionSecretCreate | PermissionSecretRead |
		PermissionSecretUpdate | PermissionSecretDelete | PermissionViewHistory |
		PermissionViewAuditLogs

	PermissionsAdmin PermissionBits = PermissionVaultAdmin | PermissionVaultShare | PermissionVaultRead |
		PermissionSecretCreate | PermissionSecretRead | PermissionSecretUpdate |
		PermissionSecretDelete | PermissionViewHistory | PermissionViewAuditLogs

	PermissionsEditor PermissionBits = PermissionVaultEditor | PermissionVaultRead |
		PermissionSecretCreate | PermissionSecretRead | PermissionSecretUpdate | PermissionSecretDelete

	PermissionsReadOnly PermissionBits = PermissionVaultRead | PermissionSecretRead
)

func (p PermissionBits) String() string {
	if p == 0 {
		return "None"
	}

	// Pre-allocate slice with capacity based on maximum possible permissions
	perms := make([]string, 0, 12)

	// Map of bit positions to permission names
	permNames := map[PermissionBits]string{
		PermissionVaultOwner:       "VaultOwner",
		PermissionVaultAdmin:       "VaultAdmin",
		PermissionVaultEditor:      "VaultEditor",
		PermissionVaultShare:       "VaultShare",
		PermissionVaultRead:        "VaultRead",
		PermissionSecretCreate:     "SecretCreate",
		PermissionSecretRead:       "SecretRead",
		PermissionSecretUpdate:     "SecretUpdate",
		PermissionSecretDelete:     "SecretDelete",
		PermissionViewHistory:      "ViewHistory",
		PermissionViewAuditLogs:    "ViewAuditLogs",
		PermissionTeamManageVaults: "TeamManageVaults",
	}

	// Check named permission sets first for more readable output
	if p == PermissionsOwner {
		return "Owner"
	} else if p == PermissionsAdmin {
		return "Admin"
	} else if p == PermissionsEditor {
		return "Editor"
	} else if p == PermissionsReadOnly {
		return "ReadOnly"
	}

	// Check individual bits
	for bit, name := range permNames {
		if p&bit != 0 {
			perms = append(perms, name)
		}
	}

	// Use strings.Builder for efficient string concatenation
	var b strings.Builder
	for i, perm := range perms {
		if i > 0 {
			b.WriteString("|")
		}
		b.WriteString(perm)
	}

	return b.String()
}

// All returns true if the permission has all of the specified permission bits
func (p PermissionBits) All(perms ...PermissionBits) bool {
	for _, perm := range perms {
		if (p & perm) != perm {
			return false
		}
	}
	return true
}

// Any returns true if the permission has any of the specified permission bits
func (p PermissionBits) Any(perms ...PermissionBits) bool {
	for _, perm := range perms {
		if (p & perm) != 0 {
			return true
		}
	}
	return false
}

// Matches returns true if the given permission bit matches of the specified permission bit.
func (p PermissionBits) Matches(perm PermissionBits) bool {
	return (p & perm) != perm
}
