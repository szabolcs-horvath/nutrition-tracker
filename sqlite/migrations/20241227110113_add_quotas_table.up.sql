CREATE TABLE IF NOT EXISTS daily_quotas
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id           INTEGER REFERENCES users NOT NULL,
    archived_date_time DATETIME,
    calories           REAL CHECK ( calories >= 0 ),
    fats               REAL CHECK ( fats >= 0 ),
    fats_saturated     REAL CHECK ( fats_saturated >= 0 AND fats_saturated <= fats ),
    carbs              REAL CHECK ( carbs >= 0 ),
    carbs_sugar        REAL CHECK ( carbs_sugar >= 0 AND carbs_sugar <= carbs ),
    carbs_slow_release REAL CHECK ( carbs_slow_release >= 0 AND carbs_slow_release <= carbs ),
    carbs_fast_release REAL CHECK ( carbs_fast_release >= 0 AND carbs_fast_release <= carbs ),
    proteins           REAL CHECK ( proteins >= 0 ),
    salt               REAL CHECK ( salt >= 0 )
);

ALTER TABLE users
ADD COLUMN daily_quota_id INTEGER REFERENCES daily_quotas;

DROP VIEW IF EXISTS items_users_view;

CREATE VIEW IF NOT EXISTS items_users_view AS SELECT users.*
FROM items
LEFT JOIN users ON items.owner_id = users.id;

DROP VIEW IF EXISTS portions_users_view;

CREATE VIEW IF NOT EXISTS portions_users_view AS SELECT users.*
FROM portions
LEFT JOIN users ON portions.owner_id = users.id;

CREATE VIEW IF NOT EXISTS users_daily_quotas_view AS
SELECT daily_quotas.*
FROM users
LEFT JOIN daily_quotas ON users.daily_quota_id = daily_quotas.id;
