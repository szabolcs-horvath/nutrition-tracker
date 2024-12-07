CREATE TABLE IF NOT EXISTS nutritions
(
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    calories_per_100g       REAL NOT NULL CHECK (0 <= calories_per_100g),
    fats_per_100g           REAL NOT NULL CHECK (0 <= fats_per_100g),
    fats_saturated_per_100g REAL CHECK (0 <= fats_saturated_per_100g AND fats_saturated_per_100g <= fats_per_100g),
    carbs_per_100g          REAL NOT NULL CHECK (0 <= carbs_per_100g),
    carbs_sugar_per_100g    REAL CHECK (0 <= carbs_sugar_per_100g AND carbs_sugar_per_100g <= carbs_per_100g),
    proteins_per_100g       REAL NOT NULL CHECK (0 <= proteins_per_100g),
    salt_per_100g           REAL CHECK (0 <= salt_per_100g)
);

CREATE TABLE IF NOT EXISTS items
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT NOT NULL,
    nutrition REFERENCES nutritions NOT NULL,
    icon      BLOB
);
