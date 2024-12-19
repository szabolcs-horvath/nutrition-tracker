ALTER TABLE items DROP COLUMN user_created;

CREATE TABLE IF NOT EXISTS portions_backup (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    owner_id        INTEGER REFERENCES users,
    language_id     INTEGER REFERENCES languages,
    liquid          BOOLEAN DEFAULT FALSE NOT NULL,
    weigth_in_grams REAL,
    volume_in_ml    REAL,
    CHECK (liquid IS FALSE OR volume_in_ml IS NOT NULL),
    CHECK (liquid IS TRUE OR weigth_in_grams IS NOT NULL),
    CHECK (volume_in_ml >= 0),
    CHECK (weigth_in_grams >= 0),
    CHECK (weigth_in_grams IS NOT NULL OR volume_in_ml IS NOT NULL)
);

INSERT INTO portions_backup
SELECT id, name, owner_id, language_id, liquid, weigth_in_grams, volume_in_ml FROM portions;

DROP TABLE IF EXISTS portions;

ALTER TABLE portions_backup RENAME TO portions;
