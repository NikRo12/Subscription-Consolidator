CREATE TABLE subs (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    url        VARCHAR(500) NOT NULL,
    price      NUMERIC(10, 2) NOT NULL
);