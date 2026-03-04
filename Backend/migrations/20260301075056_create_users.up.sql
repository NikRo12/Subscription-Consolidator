CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL,
    refresh_token VARCHAR(255) UNIQUE NOT NULL
);