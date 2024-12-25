CREATE TABLE IF NOT EXISTS notifications_backup
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id      INTEGER REFERENCES users NOT NULL,
    time          TIME NOT NULL,
    delay_seconds INTEGER,
    delay_date    DATE,
    UNIQUE (owner_id, time),
    CHECK (
        (
            delay_seconds IS NULL
            AND
            delay_date IS NULL
        )
        OR
        (
            delay_seconds IS NOT NULL
            AND
            delay_date IS NOT NULL
            AND
            (
                CAST(strftime('%H', time) AS INTEGER) * 3600 +
                CAST(strftime('%M', time) AS INTEGER) * 60 +
                CAST(strftime('%f', time) AS REAL)
            ) + delay_seconds < 86400 -- 24 hours
        )
    )
);

INSERT INTO notifications_backup
SELECT id, owner_id, time, CAST(strftime('%H', delay) AS INTEGER) * 3600 + CAST(strftime('%M', delay) AS INTEGER) * 60 + CAST(strftime('%f', delay) AS REAL), delay_date FROM notifications;

DROP VIEW IF EXISTS meals_notifications_view;

DROP TABLE IF EXISTS notifications;

ALTER TABLE notifications_backup RENAME TO notifications;

CREATE VIEW IF NOT EXISTS meals_notifications_view AS SELECT notifications.*
FROM meals
LEFT JOIN notifications ON meals.notification_id = notifications.id;
