SELECT id, ref_id, status, user_id, created_at, updated_at FROM email_verifications WHERE ref_id = $1;
