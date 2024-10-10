DROP TABLE oauth2_clients;

ALTER TABLE users ADD COLUMN created_at timestamp;
ALTER TABLE users DROP CONSTRAINT username_uniq;
ALTER TABLE refresh_tokens ADD COLUMN created_at timestamp;
