CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    display_name VARCHAR,
    username VARCHAR,
    hashed_pass VARCHAR,
    created_at timestamp,
    updated_at timestamp
);
