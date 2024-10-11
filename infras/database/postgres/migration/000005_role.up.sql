ALTER TABLE users ADD COLUMN role VARCHAR;

UPDATE users
SET role='admin'
WHERE id IN (SELECT id FROM users ORDER BY id ASC LIMIT 1);

UPDATE users
SET role='user'
WHERE id NOT IN (SELECT id FROM users ORDER BY id ASC LIMIT 1);
