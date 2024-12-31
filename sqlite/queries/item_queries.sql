-- name: ListItems :many
SELECT DISTINCT sqlc.embed(items), sqlc.embed(items_users_view), sqlc.embed(languages), sqlc.embed(portions)
FROM items
LEFT JOIN items_users_view ON items.owner_id = items_users_view.id
JOIN languages ON items.language_id = languages.id
JOIN portions ON items.default_portion_id = portions.id;

-- name: FindItemById :one
SELECT DISTINCT sqlc.embed(items), sqlc.embed(items_users_view), sqlc.embed(languages), sqlc.embed(portions)
FROM items
LEFT JOIN items_users_view ON items.owner_id = items_users_view.id
JOIN languages ON items.language_id = languages.id
JOIN portions ON items.default_portion_id = portions.id
WHERE items.id = ?
LIMIT 1;

-- name: SearchItemsByNameAndUser :many
SELECT DISTINCT sqlc.embed(items), sqlc.embed(items_users_view), sqlc.embed(languages), sqlc.embed(portions)
FROM items
LEFT JOIN items_users_view ON items.owner_id = items_users_view.id
JOIN languages ON items.language_id = languages.id
JOIN portions ON items.default_portion_id = portions.id
WHERE (items.owner_id IS NULL OR items.owner_id IS ?)
AND items.name LIKE '%' || sqlc.arg(query) || '%' COLLATE NOCASE;

-- name: CreateItem :one
INSERT INTO items(name,
                  owner_id,
                  language_id,
                  liquid,
                  default_portion_id,
                  calories_per_100,
                  fats_per_100,
                  fats_saturated_per_100,
                  carbs_per_100,
                  carbs_sugar_per_100,
                  carbs_slow_release_per_100,
                  carbs_fast_release_per_100,
                  proteins_per_100,
                  salt_per_100)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: CreateItemsPortionsJoiningTableRecord :exec
INSERT INTO items_portions_joining_table(item_id, portion_id)
VALUES (?, ?);

-- name: UpdateItem :one
UPDATE items
SET name = ?,
    owner_id = ?,
    language_id = ?,
    liquid = ?,
    default_portion_id = ?,
    calories_per_100 = ?,
    fats_per_100 = ?,
    fats_saturated_per_100 = ?,
    carbs_per_100 = ?,
    carbs_sugar_per_100 = ?,
    carbs_slow_release_per_100 = ?,
    carbs_fast_release_per_100 = ?,
    proteins_per_100 = ?,
    salt_per_100 = ?
WHERE id = ?
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = ?;

-- name: GetOwnerIdByItemId :one
SELECT owner_id
FROM items
WHERE id = ?
LIMIT 1;
