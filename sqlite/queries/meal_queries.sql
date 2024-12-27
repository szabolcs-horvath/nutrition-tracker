-- name: FindMealById :one
SELECT DISTINCT sqlc.embed(meals), sqlc.embed(users), sqlc.embed(meals_notifications_view)
FROM meals
JOIN users ON meals.owner_id = users.id
LEFT JOIN meals_notifications_view ON meals.notification_id = meals_notifications_view.id
WHERE meals.id = ?
LIMIT 1;

-- name: FindMealsForUser :many
SELECT DISTINCT sqlc.embed(meals), sqlc.embed(users), sqlc.embed(meals_notifications_view)
FROM meals
JOIN users ON meals.owner_id = users.id
LEFT JOIN meals_notifications_view ON meals.notification_id = meals_notifications_view.id
WHERE meals.owner_id = ?
AND meals.archived = ?;

-- name: CreateMeal :one
INSERT INTO meals(owner_id,
                  notification_id,
                  name,
                  time,
                  calories_quota,
                  fats_quota,
                  fats_saturated_quota,
                  carbs_quota,
                  carbs_sugar_quota,
                  carbs_slow_release_quota,
                  carbs_fast_release_quota,
                  proteins_quota,
                  salt_quota)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateMeal :one
UPDATE meals
SET notification_id = ?,
    name = ?,
    time = ?,
    calories_quota = ?,
    fats_quota = ?,
    fats_saturated_quota = ?,
    carbs_quota = ?,
    carbs_sugar_quota = ?,
    carbs_slow_release_quota = ?,
    carbs_fast_release_quota = ?,
    proteins_quota = ?,
    salt_quota = ?
WHERE id = ?
RETURNING *;

-- name: ArchiveMeal :exec
UPDATE meals
SET archived = TRUE
WHERE id = ?;
