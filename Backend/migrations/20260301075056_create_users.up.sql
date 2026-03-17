CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    google_id     VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) NOT NULL,
    access_token  VARCHAR(255) NOT NULL
);