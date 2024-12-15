PRAGMA ignore_check_constraints = ON;

ALTER TABLE portions ADD COLUMN liquid BOOLEAN DEFAULT FALSE
    CHECK ( liquid IS FALSE OR volume_in_ml IS NOT NULL )
    CHECK ( liquid IS TRUE OR weigth_in_grams IS NOT NULL );

UPDATE portions
SET liquid = TRUE
WHERE volume_in_ml IS NOT NULL;

PRAGMA ignore_check_constraints = OFF;
