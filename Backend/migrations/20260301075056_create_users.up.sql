CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    refresh_token VARCHAR(255) UNIQUE NOT NULL
);