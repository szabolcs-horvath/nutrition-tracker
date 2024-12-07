ALTER TABLE items
    ADD COLUMN owner INTEGER REFERENCES users;

ALTER TABLE items
    ADD COLUMN user_created BOOLEAN NOT NULL DEFAULT FALSE
        CHECK ( user_created = 0 OR owner IS NOT NULL );

ALTER TABLE items
    ADD COLUMN original_nutrition_id INTEGER REFERENCES nutritions
        CHECK ( user_created = 1 OR nutrition_id = original_nutrition_id );

UPDATE items
    SET original_nutrition_id = nutrition_id
    WHERE user_created = 0;
