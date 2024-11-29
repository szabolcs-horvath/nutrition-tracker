-- name: ListItemsWithNutritions :many
SELECT sqlc.embed(items), sqlc.embed(nutritions)
FROM items JOIN nutritions ON items.nutrition_id = nutritions.id;

-- name: FindItemById :one
SELECT *
FROM items
WHERE id = ?
LIMIT 1;

-- name: FindItemByIdWithNutrition :one
SELECT sqlc.embed(items), sqlc.embed(nutritions)
FROM items JOIN nutritions ON items.nutrition_id = nutritions.id
where items.id = ?
LIMIT 1;

-- name: CreateItem :one
INSERT INTO items(name,
                  nutrition_id,
                  icon)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateItem :one
UPDATE items
SET name = ?,
    nutrition_id = ?,
    icon = ?
WHERE id = ?
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = ?;
