-- name: ListUsers :many
SELECT DISTINCT sqlc.embed(users), sqlc.embed(languages), sqlc.embed(users_daily_quotas_view)
FROM users
JOIN languages ON users.language_id = languages.id
LEFT JOIN users_daily_quotas_view ON users.daily_quota_id = users_daily_quotas_view.id;

-- name: FindUserById :one
SELECT DISTINCT sqlc.embed(users), sqlc.embed(languages), sqlc.embed(users_daily_quotas_view)
FROM users
JOIN languages ON users.language_id = languages.id
LEFT JOIN users_daily_quotas_view ON users.daily_quota_id = users_daily_quotas_view.id
WHERE users.id = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users(language_id, daily_quota_id)
VALUES (?, ?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET language_id = ?,
    daily_quota_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
