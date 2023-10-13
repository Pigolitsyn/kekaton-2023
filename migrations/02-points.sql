
-- +migrate Up

CREATE TABLE IF NOT EXISTS points
(
    id          SERIAL       NOT NULL PRIMARY KEY,
    coordinates POINT        NOT NULL,
    address     VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    open_time   TIME,
    close_time  TIME,
    created_by  UUID
);

CREATE TABLE IF NOT EXISTS tag_types
(
    id   SERIAL      NOT NULL PRIMARY KEY,
    name VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS tags
(
    id       SERIAL  NOT NULL PRIMARY KEY,
    point_id INTEGER,
    type_id  INTEGER,

    FOREIGN KEY (point_id) REFERENCES points    (id),
    FOREIGN KEY (type_id)  REFERENCES tag_types (id)
);

-- +migrate Down

DROP TABLE IF EXISTS tags;

DROP TABLE IF EXISTS tag_types;

DROP TABLE IF EXISTS points;
