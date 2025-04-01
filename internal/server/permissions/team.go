package permissions

type TeamPermission PermissionBit

const (
	// Team permissions
	PermissionTeamMembersRead   TeamPermission = 0x0000_0000_0000_0001
	PermissionTeamMembersUpdate TeamPermission = 0x0000_0000_0000_0002
	PermissionTeamMembersCreate TeamPermission = 0x0000_0000_0000_0004
	PermissionTeamMembersDelete TeamPermission = 0x0000_0000_0000_0008
	PermissionTeamVaultRead     TeamPermission = 0x0000_0000_0000_0010
	PermissionTeamVaultUpdate   TeamPermission = 0x0000_0000_0000_0020
	PermissionTeamVaultCreate   TeamPermission = 0x0000_0000_0000_0040
	PermissionTeamVaultDelete   TeamPermission = 0x0000_0000_0000_0080
)

// String returns a human-readable permission string
func (p TeamPermission) String() string {
	switch p {
	case PermissionTeamMembersRead:
		return "team.members.read"
	case PermissionTeamMembersUpdate:
		return "team.members.update"
	case PermissionTeamMembersCreate:
		return "team.members.create"
	case PermissionTeamMembersDelete:
		return "team.members.delete"
	case PermissionTeamVaultRead:
		return "team.vault.read"
	case PermissionTeamVaultUpdate:
		return "team.vault.update"
	case PermissionTeamVaultCreate:
		return "team.vault.create"
	case PermissionTeamVaultDelete:
		return "team.vault.delete"
	}
	return ""
}

// AllTeamPermissions checks the given user permission has all the required team permissions.
func AllTeamPermissions(userPerm TeamPermission, perms ...TeamPermission) bool {
	for _, perm := range perms {
		if userPerm&perm == 0 {
			return false
		}
	}

	return true
}

// AnyTeamPermission checks the given user permission has any of the team permissions.
func AnyTeamPermission(userPerm TeamPermission, perms ...TeamPermission) bool {
	for _, perm := range perms {
		if userPerm&perm == perm {
			return true
		}
	}

	return false
}
