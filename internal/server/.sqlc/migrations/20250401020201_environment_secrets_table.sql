-- +goose Up
-- +goose StatementBegin
-- Modify the vault_secrets table to remove UNIQUE constraint
DROP INDEX IF EXISTS unique_vault_secret_name;
ALTER TABLE vault_secrets ADD COLUMN environment_id TEXT REFERENCES environments(environment_id) ON DELETE CASCADE;
-- Re-add the unique constraint including environment_id
CREATE UNIQUE INDEX unique_vault_env_secret_name ON vault_secrets(vault_id, environment_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the new unique constraint
DROP INDEX IF EXISTS unique_vault_env_secret_name;
-- Remove the environment_id column
CREATE TABLE vault_secrets_temp (
    secret_id TEXT NOT NULL PRIMARY KEY,
    vault_id TEXT NOT NULL CHECK (vault_id != ''),
    name TEXT NOT NULL CHECK (name != ''),
    value BLOB NOT NULL,
    created_at TEXT NOT NULL CHECK (created_at != ''),
    updated_at TEXT NOT NULL CHECK (updated_at != '')
);
INSERT INTO vault_secrets_temp SELECT secret_id, vault_id, name, value, created_at, updated_at FROM vault_secrets;
DROP TABLE vault_secrets;
ALTER TABLE vault_secrets_temp RENAME TO vault_secrets;
-- Re-add the original unique constraint
CREATE UNIQUE INDEX unique_vault_secret_name ON vault_secrets(vault_id, name);
-- Re-add foreign key constraint
ALTER TABLE vault_secrets ADD CONSTRAINT fk_vault_id FOREIGN KEY (vault_id) REFERENCES vaults(vault_id) ON DELETE CASCADE;
-- +goose StatementEnd
