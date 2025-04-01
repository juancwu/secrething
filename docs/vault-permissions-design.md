# Vault and Permissions System Design

This document outlines a simplified permissions design for the vault system.

## Database Schema

```sql
-- Vaults (core resource)
CREATE TABLE vaults (
    vault_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_by_user_id TEXT NOT NULL,  -- Who created it (always a user)
    owner_type TEXT NOT NULL,          -- "user" or "team"
    owner_id TEXT NOT NULL,            -- user_id or team_id
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (created_by_user_id) REFERENCES users(user_id)
)

-- Secrets within vaults
CREATE TABLE secrets (
    secret_id TEXT PRIMARY KEY,
    vault_id TEXT NOT NULL,
    name TEXT NOT NULL,
    value BLOB NOT NULL,  -- Encrypted value
    created_by TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(user_id),
    UNIQUE (vault_id, name)
)

-- Permissions for vault access
CREATE TABLE permissions (
    permission_id TEXT PRIMARY KEY,
    vault_id TEXT NOT NULL,
    grantee_type TEXT NOT NULL,        -- "user" or "team"
    grantee_id TEXT NOT NULL,          -- user_id or team_id
    permission_bits BIGINT NOT NULL,   -- Bitmask for granular permissions
    granted_by TEXT NOT NULL,          -- user_id who granted access
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(user_id),
    UNIQUE (vault_id, grantee_type, grantee_id)
)
```

## Permission Structure

```go
const (
    // Vault-level permissions
    PermissionVaultOwner         = 1 << 0  // 0x0001 - Can transfer ownership
    PermissionVaultAdmin         = 1 << 1  // 0x0002 - Can manage permissions
    PermissionVaultShare         = 1 << 2  // 0x0004 - Can share with others
    PermissionVaultRead          = 1 << 3  // 0x0008 - Can view vault metadata
    
    // Secret-level permissions
    PermissionSecretCreate       = 1 << 4  // 0x0010 - Can create secrets
    PermissionSecretRead         = 1 << 5  // 0x0020 - Can read secrets
    PermissionSecretUpdate       = 1 << 6  // 0x0040 - Can update secrets
    PermissionSecretDelete       = 1 << 7  // 0x0080 - Can delete secrets
    
    // History and audit permissions
    PermissionViewHistory        = 1 << 8  // 0x0100 - Can view access history
    PermissionViewAuditLogs      = 1 << 9  // 0x0200 - Can view audit logs
    
    // Common permission sets
    PermissionsOwner = PermissionVaultOwner | PermissionVaultAdmin | PermissionVaultShare |
                      PermissionVaultRead | PermissionSecretCreate | PermissionSecretRead |
                      PermissionSecretUpdate | PermissionSecretDelete | PermissionViewHistory |
                      PermissionViewAuditLogs
                      
    PermissionsAdmin = PermissionVaultAdmin | PermissionVaultShare | PermissionVaultRead |
                      PermissionSecretCreate | PermissionSecretRead | PermissionSecretUpdate |
                      PermissionSecretDelete | PermissionViewHistory | PermissionViewAuditLogs
                      
    PermissionsReadOnly = PermissionVaultRead | PermissionSecretRead
)
```

## Core Business Rules

1. **Vault Creation:**
   - Any user can create a vault
   - Upon creation, they become the owner (with PermissionVaultOwner)
   - Team members with PermissionTeamManageVaults can create vaults owned by the team

2. **Ownership:**
   - A vault has exactly one owner (user or team)
   - Only owners can transfer ownership
   - When a team owns a vault, the team's permission structure determines access

3. **Permission Inheritance:**
   - Team permissions cascade to all team members
   - If a user has both direct and team-based access, permissions are combined
   - A user can have multiple permission paths to the same vault

4. **Access Control:**
   - Only users with PermissionVaultShare can share vaults
   - Users can't grant permissions they don't have themselves
   - Vault owners always retain ownership permissions