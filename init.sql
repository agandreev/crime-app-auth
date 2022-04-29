CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    login    VARCHAR(16),
    password VARCHAR(60),
    UNIQUE(login)
);