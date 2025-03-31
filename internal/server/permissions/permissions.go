package permissions

type VaultPermission uint64
type TeamPermission uint64

const (
	// Vault permissions
	PermissionVaultRead   VaultPermission = 0x0000_0000_0000_0001
	PermissionVaultUpdate VaultPermission = 0x0000_0000_0000_0002
	PermissionVaultCreate VaultPermission = 0x0000_0000_0000_0004
	PermissionVaultDelete VaultPermission = 0x0000_0000_0000_0008

	// Team permissions
	PermissionTeamMembersRead   TeamPermission = 0x0000_0000_0000_0010
	PermissionTeamMembersUpdate TeamPermission = 0x0000_0000_0000_0020
	PermissionTeamMembersCreate TeamPermission = 0x0000_0000_0000_0040
	PermissionTeamMembersDelete TeamPermission = 0x0000_0000_0000_0080
	PermissionTeamVaultRead     TeamPermission = 0x0000_0000_0000_0100
	PermissionTeamVaultUpdate   TeamPermission = 0x0000_0000_0000_0200
	PermissionTeamVaultCreate   TeamPermission = 0x0000_0000_0000_0400
	PermissionTeamVaultDelete   TeamPermission = 0x0000_0000_0000_0800
)
