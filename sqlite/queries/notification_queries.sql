-- name: FindNotificationById :one
SELECT *
FROM notifications
WHERE id = ?
LIMIT 1;

-- name: FindNotificationByIdWithRelations :one
SELECT sqlc.embed(notifications), sqlc.embed(users)
FROM notifications
JOIN users ON notifications.owner_id = users.id
WHERE notifications.id = ?
LIMIT 1;

-- name: FindNotificationsByUserId :many
SELECT sqlc.embed(notifications), sqlc.embed(users)
FROM notifications
JOIN users ON notifications.owner_id = users.id
WHERE notifications.owner_id = ?;

-- name: FindNotificationByUserIdAndTime :one
SELECT sqlc.embed(notifications), sqlc.embed(users)
FROM notifications
JOIN users ON notifications.owner_id = users.id
WHERE notifications.owner_id = ?
    AND notifications.time = ?
LIMIT 1;

-- name: CreateNotification :one
INSERT INTO notifications(owner_id,
                          time,
                          delay,
                          delay_date,
                          name)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateNotification :one
UPDATE notifications
SET owner_id = ?,
    time = ?,
    delay = ?,
    delay_date = ?,
    name = ?
WHERE id = ?
RETURNING *;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = ?;
