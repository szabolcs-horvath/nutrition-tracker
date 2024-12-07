ALTER TABLE items
    RENAME COLUMN owner_id TO owner;

ALTER TABLE nutritions
    RENAME COLUMN owner_id TO owner;

ALTER TABLE notifications
    RENAME COLUMN owner_id TO owner;
