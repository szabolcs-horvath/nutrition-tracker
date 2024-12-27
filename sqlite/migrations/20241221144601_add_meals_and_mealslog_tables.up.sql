CREATE TABLE IF NOT EXISTS meals
(
    id                       INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id                 INTEGER REFERENCES users NOT NULL,
    notification_id          INTEGER REFERENCES notifications,
    name                     TEXT NOT NULL,
    time                     TIME NOT NULL,
    calories_quota           REAL,
    fats_quota               REAL,
    fats_saturated_quota     REAL,
    carbs_quota              REAL,
    carbs_sugar_quota        REAL,
    carbs_slow_release_quota REAL,
    carbs_fast_release_quota REAL,
    proteins_quota           REAL,
    salt_quota               REAL,
    archived                 BOOLEAN NOT NULL DEFAULT FALSE,
    CHECK (
        (
            CAST(strftime('%H', time) AS INTEGER) * 3600 +
            CAST(strftime('%M', time) AS INTEGER) * 60 +
            CAST(strftime('%f', time) AS REAL)
        ) < 86400 -- 24 hours
    ),
    CHECK ( calories_quota >= 0 ),
    CHECK ( fats_quota >= 0 ),
    CHECK ( fats_saturated_quota >= 0 ),
    CHECK ( carbs_quota >= 0 ),
    CHECK ( carbs_sugar_quota >= 0 ),
    CHECK ( carbs_slow_release_quota >= 0 ),
    CHECK ( carbs_fast_release_quota >= 0 ),
    CHECK ( proteins_quota >= 0 ),
    CHECK ( salt_quota >= 0 ),
    CHECK ( fats_saturated_quota <= fats_quota ),
    CHECK ( carbs_sugar_quota <= carbs_quota ),
    CHECK ( carbs_slow_release_quota + carbs_fast_release_quota <= carbs_quota )
);

CREATE VIEW IF NOT EXISTS meals_notifications_view AS SELECT notifications.*
FROM meals
LEFT JOIN notifications ON meals.notification_id = notifications.id;

CREATE TABLE IF NOT EXISTS meallogs
(
    id                 INTEGER PRIMARY KEY AUTOINCREMENT,
    meal_id            INTEGER REFERENCES meals NOT NULL,
    item_id            INTEGER REFERENCES items NOT NULL,
    portion_id         INTEGER REFERENCES portions NOT NULL,
    portion_multiplier REAL NOT NULL,
    datetime           DATETIME NOT NULL,
    CHECK (
        (
            CAST(strftime('%H', datetime) AS INTEGER) * 3600 +
            CAST(strftime('%M', datetime) AS INTEGER) * 60 +
            CAST(strftime('%f', datetime) AS REAL)
            ) < 86400 -- 24 hours
        )
)
