ALTER TABLE nutritions
    ADD COLUMN carbs_slow_release_per_100g REAL
        CHECK (carbs_slow_release_per_100g >= 0 AND carbs_slow_release_per_100g <= carbs_per_100g);

ALTER TABLE nutritions
    ADD COLUMN carbs_fast_release_per_100g REAL
        CHECK (carbs_fast_release_per_100g >= 0 AND carbs_fast_release_per_100g <= carbs_per_100g)
        CHECK (carbs_slow_release_per_100g + carbs_fast_release_per_100g <= carbs_per_100g);
