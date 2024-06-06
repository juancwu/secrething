CREATE TRIGGER update_personal_bento_access_updated_at
BEFORE UPDATE ON personal_bento_access
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column(); 
