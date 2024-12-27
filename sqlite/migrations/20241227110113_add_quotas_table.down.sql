DROP VIEW IF EXISTS users_daily_quotas_view;

DROP VIEW IF EXISTS items_users_view;

CREATE VIEW IF NOT EXISTS items_users_view AS SELECT users.*
FROM items
LEFT JOIN users ON items.owner_id = users.id;

DROP VIEW IF EXISTS portions_users_view;

CREATE VIEW IF NOT EXISTS portions_users_view AS SELECT users.*
FROM portions
LEFT JOIN users ON portions.owner_id = users.id;

ALTER TABLE users
DROP COLUMN daily_quota_id;

DROP TABLE IF EXISTS daily_quotas;
