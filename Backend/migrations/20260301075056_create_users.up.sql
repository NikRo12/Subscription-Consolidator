CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    google_id     TEXT UNIQUE NOT NULL,
    refresh_token TEXT NOT NULL,
    access_token  TEXT NOT NULL
);