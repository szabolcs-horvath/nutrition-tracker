-- name: FindPortionsForItemAndUser :many
SELECT sqlc.embed(portions)
FROM portions
JOIN items_portions_joining_table on portions.id = items_portions_joining_table.portion_id
JOIN items on items_portions_joining_table.item_id = items.id
WHERE items.id = sqlc.arg(item_id)
AND portions.owner_id = sqlc.narg(owner_id)
AND portions.owner_id IS NOT NULL
    UNION ALL
SELECT sqlc.embed(portions)
FROM portions
JOIN items ON portions.liquid = items.liquid
WHERE items.id = sqlc.arg(item_id)
AND portions.owner_id IS NULL;

-- name: FindPortionById :one
SELECT sqlc.embed(portions), sqlc.embed(portions_users_view), sqlc.embed(portions_languages_view)
FROM portions
LEFT JOIN portions_users_view ON portions.owner_id = portions_users_view.id
LEFT JOIN portions_languages_view ON portions.language_id = portions_languages_view.id
WHERE portions.id = ?
LIMIT 1;

-- name: CreatePortion :one
INSERT INTO portions(name,
                     owner_id,
                     language_id,
                     liquid,
                     weigth_in_grams,
                     volume_in_ml)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdatePortion :one
UPDATE portions
SET name = ?,
    owner_id = ?,
    language_id = ?,
    liquid = ?,
    weigth_in_grams = ?,
    volume_in_ml = ?
WHERE id = ?
RETURNING *;

-- name: DeletePortion :exec
DELETE FROM portions
WHERE id = ?;
