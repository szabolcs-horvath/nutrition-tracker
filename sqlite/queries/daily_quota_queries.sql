-- name: FindDailyQuotaById :one
SELECT sqlc.embed(daily_quotas), sqlc.embed(users)
FROM daily_quotas
JOIN users ON daily_quotas.owner_id = users.id
WHERE daily_quotas.id = ?
LIMIT 1;

-- name: ListDailyQuotasForUser :many
SELECT sqlc.embed(daily_quotas)
FROM daily_quotas
WHERE owner_id = ?;

-- name: FindDailyQuotaByOwnerAndDate :one
SELECT sqlc.embed(daily_quotas)
FROM daily_quotas
WHERE owner_id = ?
AND (archived_date_time IS NULL OR archived_date_time > sqlc.arg(date))
ORDER BY archived_date_time NULLS LAST
LIMIT 1;

-- name: CreateDailyQuota :one
INSERT INTO daily_quotas(owner_id,
                         calories,
                         fats,
                         fats_saturated,
                         carbs,
                         carbs_sugar,
                         carbs_slow_release,
                         carbs_fast_release,
                         proteins,
                         salt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateDailyQuota :one
UPDATE daily_quotas
SET owner_id = ?,
    calories = ?,
    fats = ?,
    fats_saturated = ?,
    carbs = ?,
    carbs_sugar = ?,
    carbs_slow_release = ?,
    carbs_fast_release = ?,
    proteins = ?,
    salt = ?
WHERE id = ?
RETURNING *;

-- name: ArchiveDailyQuota :exec
UPDATE daily_quotas
SET archived_date_time = datetime('now')
WHERE id = ?;