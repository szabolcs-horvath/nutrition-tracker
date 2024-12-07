ALTER TABLE items
    RENAME COLUMN owner TO owner_id;

ALTER TABLE nutritions
    RENAME COLUMN owner TO owner_id;

ALTER TABLE notifications
    RENAME COLUMN owner TO owner_id;
