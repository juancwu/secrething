package permissions

type VaultPermission uint64

const (
	// Vault permissions
	PermissionVaultRead   VaultPermission = 0x0000_0000_0000_0001
	PermissionVaultUpdate VaultPermission = 0x0000_0000_0000_0002
	PermissionVaultCreate VaultPermission = 0x0000_0000_0000_0004
	PermissionVaultDelete VaultPermission = 0x0000_0000_0000_0008
)
