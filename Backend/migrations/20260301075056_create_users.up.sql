CREATE TABLE users (
    id            SERIAL PRIMARY KEY,
    google_id     VARCHAR(1025) UNIQUE NOT NULL,
    refresh_token VARCHAR(1025) NOT NULL,
    access_token  VARCHAR(1025) NOT NULL
);