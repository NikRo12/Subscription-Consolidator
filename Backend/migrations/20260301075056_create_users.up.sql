CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    refresh_token varchar not null 
)