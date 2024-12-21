-- name: FindMealLogsForUserAndDate :many
SELECT sqlc.embed(meallogs), sqlc.embed(meals), sqlc.embed(items), sqlc.embed(portions)
FROM meallogs
JOIN meals ON meallogs.meal_id = meals.id
JOIN items ON meallogs.item_id = items.id
JOIN portions ON meallogs.portion_id = portions.id
WHERE meals.owner_id = ?
AND strftime('%d', meallogs.datetime) = strftime('%d', date('now'));

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
