-- name: FindNotificationById :one
SELECT sqlc.embed(notifications), sqlc.embed(users)
FROM notifications
JOIN users ON notifications.owner_id = users.id
WHERE notifications.id = ?
LIMIT 1;

-- name: ListNotificationsByUserId :many
SELECT sqlc.embed(notifications), sqlc.embed(users)
FROM notifications
JOIN users ON notifications.owner_id = users.id
WHERE notifications.owner_id = ?;

-- name: CreateNotification :one
INSERT INTO notifications(owner_id,
                          time,
                          delay_seconds,
                          delay_date)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: UpdateNotification :one
UPDATE notifications
SET owner_id = ?,
    time = ?,
    delay_seconds = ?,
    delay_date = ?
WHERE id = ?
RETURNING *;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = ?;
