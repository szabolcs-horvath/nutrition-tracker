CREATE TABLE IF NOT EXISTS nutritions (
    id                          INTEGER PRIMARY KEY AUTOINCREMENT,
    calories_per_100g           REAL NOT NULL DEFAULT 0
        CHECK ( calories_per_100g >= 0 ),
    fats_per_100g               REAL NOT NULL DEFAULT 0
        CHECK ( fats_per_100g >= 0 ),
    fats_saturated_per_100g     REAL
        CHECK ( fats_saturated_per_100g >= 0 AND fats_saturated_per_100g <= fats_per_100g ),
    carbs_per_100g              REAL NOT NULL DEFAULT 0
        CHECK ( carbs_per_100g >= 0 ),
    carbs_sugar_per_100g        REAL
        CHECK ( carbs_sugar_per_100g >= 0 AND carbs_sugar_per_100g <= carbs_per_100g ),
    carbs_slow_release_per_100g REAL
        CHECK ( carbs_slow_release_per_100g >= 0 AND carbs_slow_release_per_100g <= carbs_per_100g ),
    carbs_fast_release_per_100g REAL
        CHECK ( carbs_fast_release_per_100g >= 0 AND carbs_fast_release_per_100g <= carbs_per_100g )
        CHECK ( carbs_fast_release_per_100g + carbs_slow_release_per_100g <= carbs_per_100g ),
    proteins_per_100g           REAL NOT NULL DEFAULT 0
        CHECK ( proteins_per_100g >= 0 ),
    salt_per_100g               REAL
        CHECK ( salt_per_100g >= 0 ),
    owner_id                    INTEGER REFERENCES users,
    user_created                BOOLEAN NOT NULL DEFAULT FALSE
        CHECK ( user_created IS FALSE OR owner_id IS NOT NULL )
);

INSERT INTO nutritions (id,
                        calories_per_100g,
                        fats_per_100g,
                        fats_saturated_per_100g,
                        carbs_per_100g,
                        carbs_sugar_per_100g,
                        carbs_slow_release_per_100g,
                        carbs_fast_release_per_100g,
                        proteins_per_100g,
                        salt_per_100g,
                        owner_id,
                        user_created)
SELECT id,
       calories_per_100g,
       fats_per_100g,
       fats_saturated_per_100g,
       carbs_per_100g,
       carbs_sugar_per_100g,
       carbs_slow_release_per_100g,
       carbs_fast_release_per_100g,
       proteins_per_100g,
       salt_per_100g,
       owner_id,
       user_created
FROM items;

ALTER TABLE items ADD COLUMN nutrition_id INTEGER REFERENCES nutritions NOT NULL DEFAULT 1;
ALTER TABLE items ADD COLUMN original_nutrition_id INTEGER REFERENCES nutritions;

UPDATE items
SET nutrition_id = nutritions.id,
    original_nutrition_id = nutritions.id
FROM nutritions
WHERE items.id = nutritions.id;

ALTER TABLE items DROP COLUMN calories_per_100g;
ALTER TABLE items DROP COLUMN fats_saturated_per_100g;
ALTER TABLE items DROP COLUMN fats_per_100g;
ALTER TABLE items DROP COLUMN carbs_fast_release_per_100g;
ALTER TABLE items DROP COLUMN carbs_slow_release_per_100g;
ALTER TABLE items DROP COLUMN carbs_sugar_per_100g;
ALTER TABLE items DROP COLUMN carbs_per_100g;
ALTER TABLE items DROP COLUMN proteins_per_100g;
ALTER TABLE items DROP COLUMN salt_per_100g;
