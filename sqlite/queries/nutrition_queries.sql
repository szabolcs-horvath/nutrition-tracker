-- name: ListNutritions :many
SELECT *
FROM nutritions;

-- name: FindNutritionById :one
SELECT *
FROM nutritions
WHERE id = ?
LIMIT 1;

-- name: CreateNutrition :one
INSERT INTO nutritions(calories_per_100g,
                       fats_per_100g,
                       fats_saturated_per_100g,
                       carbs_per_100g,
                       carbs_sugar_per_100g,
                       proteins_per_100g,
                       salt_per_100g)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;
