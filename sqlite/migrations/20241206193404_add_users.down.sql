DROP TABLE IF EXISTS users;

ALTER TABLE nutritions
    DROP COLUMN user_created;

ALTER TABLE nutritions
    DROP COLUMN owner;
