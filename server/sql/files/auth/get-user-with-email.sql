SELECT 
    id, 
    first_name, 
    last_name, 
    email, 
    email_verified, 
    created_at, 
    updated_at, 
    pgp_sym_decrypt(pub_key, $2) 
FROM 
    users 
WHERE 
    email = $1;
