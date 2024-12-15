CREATE TABLE IF NOT EXISTS items_portions_joining_table (
    item_id    INTEGER REFERENCES items NOT NULL,
    portion_id INTEGER REFERENCES portions NOT NULL,
    PRIMARY KEY (item_id, portion_id)
);
