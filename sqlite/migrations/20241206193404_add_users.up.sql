CREATE TABLE IF NOT EXISTS users
(
  id INTEGER PRIMARY KEY AUTOINCREMENT
);

ALTER TABLE nutritions
    ADD COLUMN owner INTEGER REFERENCES users;

ALTER TABLE nutritions
    ADD COLUMN user_created BOOLEAN NOT NULL DEFAULT FALSE
        CHECK ( user_created = 0 OR owner IS NOT NULL );
