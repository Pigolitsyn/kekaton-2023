
-- +migrate Up

CREATE TABLE IF NOT EXISTS users
(
    id       UUID        NOT NULL UNIQUE PRIMARY KEY DEFAULT gen_random_uuid(),
    email    VARCHAR(32) NOT NULL UNIQUE,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(64) NOT NULL,
    salt     VARCHAR(8)  NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_users_email ON users (email);

-- +migrate Down

DROP INDEX IF EXISTS uidx_users_email ON users;

DROP TABLE IF EXISTS users;
