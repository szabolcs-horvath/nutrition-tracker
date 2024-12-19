CREATE TABLE IF NOT EXISTS portions_backup
(
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT                  NOT NULL,
    owner_id        INTEGER REFERENCES users,
    language_id     INTEGER REFERENCES languages,
    weigth_in_grams REAL,
    volume_in_ml    REAL,
    CHECK (volume_in_ml >= 0),
    CHECK (weigth_in_grams >= 0),
    CHECK (weigth_in_grams IS NOT NULL OR volume_in_ml IS NOT NULL)
);

INSERT INTO portions_backup
SELECT id, name, owner_id, language_id, weigth_in_grams, volume_in_ml FROM portions;

DROP TABLE IF EXISTS portions;

ALTER TABLE portions_backup RENAME TO portions;

