ALTER TABLE items ADD COLUMN user_created BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE portions ADD COLUMN user_created BOOLEAN NOT NULL DEFAULT FALSE;
