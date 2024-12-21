-- name: ListNotificationsByUserId :many
SELECT sqlc.embed(notifications), sqlc.embed(users)
FROM notifications
JOIN users ON notifications.owner_id = users.id
WHERE notifications.owner_id = ?;

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
