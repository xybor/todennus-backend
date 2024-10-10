ALTER TABLE users ADD COLUMN allowed_scope VARCHAR;
ALTER TABLE oauth2_clients ADD COLUMN allowed_scope VARCHAR;

UPDATE users
SET allowed_scope='*'
WHERE allowed_scope IS NULL;

UPDATE oauth2_clients
SET allowed_scope='*'
WHERE allowed_scope IS NULL AND is_confidential IS TRUE;

UPDATE oauth2_clients
SET allowed_scope='read'
WHERE allowed_scope IS NULL AND is_confidential IS FALSE;
