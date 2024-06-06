CREATE TRIGGER update_personal_bentos_updated_at
BEFORE UPDATE ON personal_bentos
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
