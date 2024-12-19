ALTER TABLE items ADD COLUMN calories_per_100g REAL NOT NULL DEFAULT 0
    CHECK ( calories_per_100g >= 0 );
ALTER TABLE items ADD COLUMN fats_per_100g REAL NOT NULL DEFAULT 0
    CHECK ( fats_per_100g >= 0 );
ALTER TABLE items ADD COLUMN fats_saturated_per_100g REAL
    CHECK ( fats_saturated_per_100g >= 0 AND fats_saturated_per_100g <= fats_per_100g );
ALTER TABLE items ADD COLUMN carbs_per_100g REAL NOT NULL DEFAULT 0
    CHECK ( carbs_per_100g >= 0 );
ALTER TABLE items ADD COLUMN carbs_sugar_per_100g REAL
    CHECK ( carbs_sugar_per_100g >= 0 AND carbs_sugar_per_100g <= carbs_per_100g );
ALTER TABLE items ADD COLUMN carbs_slow_release_per_100g REAL
    CHECK ( carbs_slow_release_per_100g >= 0 AND carbs_slow_release_per_100g <= carbs_per_100g );
ALTER TABLE items ADD COLUMN carbs_fast_release_per_100g REAL
    CHECK ( carbs_fast_release_per_100g >= 0 AND carbs_fast_release_per_100g <= carbs_per_100g )
    CHECK ( carbs_fast_release_per_100g + carbs_slow_release_per_100g <= carbs_per_100g );
ALTER TABLE items ADD COLUMN proteins_per_100g REAL NOT NULL DEFAULT 0
    CHECK ( proteins_per_100g >= 0 );
ALTER TABLE items ADD COLUMN salt_per_100g REAL
    CHECK ( salt_per_100g >= 0 );

UPDATE items
SET calories_per_100g = (SELECT calories_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    fats_per_100g = (SELECT fats_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    fats_saturated_per_100g = (SELECT fats_saturated_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    carbs_per_100g = (SELECT carbs_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    carbs_sugar_per_100g = (SELECT carbs_sugar_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    carbs_slow_release_per_100g = (SELECT carbs_slow_release_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    carbs_fast_release_per_100g = (SELECT carbs_fast_release_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    proteins_per_100g = (SELECT proteins_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id),
    salt_per_100g = (SELECT salt_per_100g FROM nutritions WHERE items.nutrition_id = nutritions.id)
WHERE EXISTS (SELECT 1 FROM nutritions WHERE items.nutrition_id = nutritions.id);


ALTER TABLE items DROP COLUMN original_nutrition_id;
ALTER TABLE items DROP COLUMN nutrition_id;

DROP TABLE IF EXISTS nutritions;

