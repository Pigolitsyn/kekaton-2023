
-- +migrate-up

CREATE TABLE IF NOT EXISTS comments
(
    comment_id   SERIAL       NOT NULL PRIMARY KEY,
    user_id      UUID         NOT NULL,
    point_id     INTEGER      NOT NULL,
    comment_text VARCHAR(600) NOT NULL,
    rating       INT2         NOT NULL
);

-- +migrate-down

DROP TABLE IF EXISTS comments;



