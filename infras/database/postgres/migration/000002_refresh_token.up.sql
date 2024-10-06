CREATE TABLE refresh_tokens (
    refresh_token_id BIGINT PRIMARY KEY,
    access_token_id BIGINT,
    seq INT,
    created_at timestamp,
    updated_at timestamp
);
