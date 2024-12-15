CREATE TABLE IF NOT EXISTS languages (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT NOT NULL,
    native_name TEXT NOT NULL
);

--Default languages
INSERT INTO languages(id, name, native_name)
VALUES (1, 'en', 'English'),
       (2, 'hu', 'magyar');

ALTER TABLE users
    ADD COLUMN language_id INTEGER REFERENCES languages NOT NULL DEFAULT 1;

CREATE TABLE IF NOT EXISTS portions (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    weigth_in_grams REAL CHECK ( weigth_in_grams >= 0 ),
    volume_in_ml    REAL CHECK ( volume_in_ml >= 0 )
        CHECK ( weigth_in_grams IS NOT NULL OR volume_in_ml IS NOT NULL ),
    owner_id        INTEGER REFERENCES users,
    user_created    BOOLEAN NOT NULL DEFAULT FALSE
        CHECK ( user_created = 0 OR owner_id IS NOT NULL ),
    language_id     INTEGER REFERENCES languages
        CHECK ( user_created = 0 OR language_id IS NOT NULL )
);

--Default portions
INSERT INTO portions(id, name, weigth_in_grams, volume_in_ml, owner_id, user_created, language_id)
VALUES (1, 'g', 1, NULL, NULL, FALSE, NULL),
       (2, 'ml', NULL, 1, NULL, FALSE, NULL);

ALTER TABLE items
    ADD COLUMN language_id INTEGER REFERENCES languages NOT NULL DEFAULT 1;

ALTER TABLE items
    ADD COLUMN liquid BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE items
    ADD COLUMN default_portion_id INTEGER REFERENCES portions NOT NULL DEFAULT 1;

UPDATE items
    SET default_portion_id = 2
    WHERE liquid = TRUE;

ALTER TABLE notifications
    ADD COLUMN name TEXT NOT NULL DEFAULT '';
