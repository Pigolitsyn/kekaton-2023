
-- +migrate Up

INSERT INTO tag_types (name) VALUES ('чистый'), ('сломан'), ('стульчак'), ('туалетная бумага');

-- +migrate Down

DELETE FROM tag_types WHERE name in ('чистый', 'сломан', 'стульчак', 'туалетная бумага');
