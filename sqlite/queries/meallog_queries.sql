-- name: FindMealLogById :one
SELECT sqlc.embed(meallogs), sqlc.embed(meals), sqlc.embed(items), sqlc.embed(portions)
From meallogs
JOIN meals ON meallogs.meal_id = meals.id
JOIN items ON meallogs.item_id = items.id
JOIN portions ON meallogs.portion_id = portions.id
WHERE meallogs.id = ?
LIMIT 1;

-- name: FindMealLogsForUserAndDate :many
SELECT sqlc.embed(meallogs), sqlc.embed(meals), sqlc.embed(items), sqlc.embed(portions)
FROM meallogs
JOIN meals ON meallogs.meal_id = meals.id
JOIN items ON meallogs.item_id = items.id
JOIN portions ON meallogs.portion_id = portions.id
WHERE meals.owner_id = ?
AND date(meallogs.datetime) IS date(sqlc.arg(date));

-- name: CreateMealLog :one
INSERT INTO meallogs(meal_id,
                     item_id,
                     portion_id,
                     portion_multiplier,
                     datetime)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateMealLog :one
UPDATE meallogs
SET meal_id = ?,
    item_id = ?,
    portion_id = ?,
    portion_multiplier = ?,
    datetime = ?
WHERE id = ?
RETURNING *;

-- name: DeleteMealLog :exec
DELETE FROM meallogs
WHERE id = ?;
