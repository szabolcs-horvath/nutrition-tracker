CREATE VIEW IF NOT EXISTS items_users_view AS SELECT users.*
FROM items
LEFT JOIN users ON items.owner_id = users.id;

CREATE VIEW IF NOT EXISTS portions_users_view AS SELECT users.*
FROM portions
LEFT JOIN users ON portions.owner_id = users.id;

CREATE VIEW IF NOT EXISTS portions_languages_view AS SELECT languages.*
FROM portions
LEFT JOIN languages ON portions.language_id = languages.id;