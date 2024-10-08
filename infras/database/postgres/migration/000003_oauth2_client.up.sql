CREATE TABLE oauth2_clients (
    id BIGINT PRIMARY KEY,
    user_id BIGINT REFERENCES users (id),
    name VARCHAR,
    is_confidential BOOLEAN,
    hashed_secret VARCHAR,
    updated_at timestamp
);

ALTER TABLE users DROP COLUMN created_at;
ALTER TABLE users ADD CONSTRAINT username_uniq UNIQUE (username);
ALTER TABLE refresh_tokens DROP COLUMN created_at;
