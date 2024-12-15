-- name: ListUsers :many
SELECT sqlc.embed(users), sqlc.embed(languages)
FROM users
JOIN languages ON users.language_id = languages.id;

-- name: FindUserById :one
SELECT sqlc.embed(users), sqlc.embed(languages)
FROM users
JOIN languages ON users.language_id = languages.id
WHERE users.id = ?
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users(language_id)
VALUES (?)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET language_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
