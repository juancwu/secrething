CREATE TRIGGER update_email_verifications_updated_at
BEFORE UPDATE ON email_verifications
FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
