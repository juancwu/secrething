CREATE TRIGGER update_users_passwords_resets_updated_at
BEFORE UPDATE ON users_passwords_resets
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
