DROP TABLE IF EXISTS languages;

ALTER TABLE users
    DROP COLUMN language_id;

DROP TABLE IF EXISTS portions;

ALTER TABLE items
    DROP COLUMN language_id;

ALTER TABLE items
    DROP COLUMN liquid;

ALTER TABLE items
    DROP COLUMN default_portion_id;

ALTER TABLE notifications
    DROP COLUMN name;
