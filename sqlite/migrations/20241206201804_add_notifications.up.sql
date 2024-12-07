CREATE TABLE IF NOT EXISTS notifications
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    owner      REFERENCES users NOT NULL,
    time       TIME NOT NULL,
    delay      TIME,
    delay_date DATE,
    UNIQUE (owner, time),
    CHECK (
        (
            delay IS NULL
            AND
            delay_date IS NULL
        )
        OR
        (
            delay IS NOT NULL
            AND
            delay_date IS NOT NULL
            AND
            (
                CAST(strftime('%H', time) AS INTEGER) * 3600 +
                CAST(strftime('%M', time) AS INTEGER) * 60 +
                CAST(strftime('%f', time) AS REAL)
            ) +
            (
                CAST(strftime('%H', delay) AS INTEGER) * 3600 +
                CAST(strftime('%M', delay) AS INTEGER) * 60 +
                CAST(strftime('%f', delay) AS REAL)
            ) < 86400 -- 24 hours
        )
    )
);