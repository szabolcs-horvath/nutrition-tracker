BEGIN TRANSACTION;

-- daily_quotas (including archived quotas)
INSERT INTO daily_quotas(id, owner_id, archived_date_time, calories, fats, fats_saturated, carbs, carbs_sugar, carbs_slow_release, carbs_fast_release, proteins,
                         salt)
VALUES (1, 1, '2021-12-27 12:54:21', 2021, 70, 20, 260, 90, 150, 110, 50, 6),
       (2, 2, '2021-12-27 12:54:21', 2021, 70, 20, 260, 90, 150, 110, 50, 6),
       (3, 1, '2022-01-01 00:00:00', 2022, 70, 20, 260, 90, 150, 110, 50, 6),
       (4, 2, '2022-01-01 00:00:00', 2022, 70, 20, 260, 90, 150, 110, 50, 6),
       (5, 1, NULL, 2023, 70, 20, 260, 90, 150, 110, 50, 6),
       (6, 2, NULL, 2023, 70, 20, 260, 90, 150, 110, 50, 6);

-- users
INSERT INTO users(id, language_id, daily_quota_id)
VALUES (1, 1, 5),
       (2, 2, 6);

-- meals
INSERT INTO meals(id, owner_id, notification_id, name, time, calories_quota, fats_quota, fats_saturated_quota, carbs_quota, carbs_sugar_quota,
                  carbs_slow_release_quota, carbs_fast_release_quota, proteins_quota, salt_quota, archived)
VALUES (1, 1, NULL, 'Breakfast', '08:00:00', 500, 20, 10, 50, 10, 20, 20, 20, 5, FALSE),
       (2, 1, NULL, 'Lunch', '12:00:00', 800, 30, 15, 80, 20, 40, 40, 30, 5, FALSE),
       (3, 1, NULL, 'Dinner', '18:00:00', 700, 25, 12, 70, 15, 30, 30, 25, 5, FALSE),
       (4, 2, NULL, 'Reggeli', '08:00:00', 500, 20, 10, 50, 10, 20, 20, 20, 5, FALSE),
       (5, 2, NULL, 'Ebéd', '12:00:00', 800, 30, 15, 80, 20, 40, 40, 30, 5, FALSE),
       (6, 2, NULL, 'Vacsora', '18:00:00', 700, 25, 12, 70, 15, 30, 30, 25, 5, FALSE);

-- languages
INSERT INTO languages(id, name, native_name)
VALUES (1, 'en', 'English'),
       (2, 'hu', 'magyar');

-- portions
INSERT INTO portions(id, name, owner_id, language_id, liquid, weigth_in_grams, volume_in_ml)
VALUES (1, 'g', NULL, NULL, FALSE, 1, NULL),
       (2, 'ml', NULL, NULL, TRUE, NULL, 1),
       (3, 'dkg', NULL, NULL, FALSE, 10, NULL),
       (4, 'cl', NULL, NULL, TRUE, NULL, 10),
       (5, 'dl', NULL, NULL, TRUE, NULL, 100),
       (6, 'piece', 1, 1, FALSE, 50, NULL),
       (7, 'darab', 2, 2, FALSE, 50, NULL),
       (8, 'cup', 1, 1, TRUE, NULL, 250),
       (9, 'csésze', 2, 2, TRUE, NULL, 250);

-- items
INSERT INTO items(id, name, owner_id, language_id, liquid, default_portion_id, calories_per_100, fats_per_100, fats_saturated_per_100, carbs_per_100,
                  carbs_sugar_per_100, carbs_slow_release_per_100, carbs_fast_release_per_100, proteins_per_100, salt_per_100)
VALUES (1, 'Apple', NULL, 1, FALSE, 6, 52, 0.2, 0.1, 14, 10, 10, 4, 0.3, 0.002),
       (2, 'Alma', NULL, 2, FALSE, 7, 52, 0.2, 0.1, 14, 10, 10, 4, 0.3, 0.002),
       (3, 'Banana', NULL, 1, FALSE, 6, 89, 0.3, 0.1, 23, 12, 10, 13, 1.1, 0.001),
       (4, 'Banán', NULL, 2, FALSE, 7, 89, 0.3, 0.1, 23, 12, 10, 13, 1.1, 0.001),
       (5, 'Bread', 1, 1, FALSE, 1, 265, 1.1, 0.2, 49, 2, 40, 9, 8.5, 0.5),
       (6, 'Kenyér', 2, 2, FALSE, 1, 265, 1.1, 0.2, 49, 2, 40, 9, 8.5, 0.5),
       (7, 'Butter', 1, 1, FALSE, 1, 717, 81.1, 51.4, 0.1, 0.1, 0, 0, 0.9, 1.8),
       (8, 'Vaj', 2, 2, FALSE, 1, 717, 81.1, 51.4, 0.1, 0.1, 0, 0, 0.9, 1.8),
       (9, 'Cheese', 1, 1, FALSE, 1, 402, 33.1, 21.1, 1.3, 0.5, 0, 0, 25, 1.6),
       (10, 'Sajt', 2, 2, FALSE, 1, 402, 33.1, 21.1, 1.3, 0.5, 0, 0, 25, 1.6),
       (11, 'Chicken', 1, 1, FALSE, 1, 165, 3.6, 1, 0, 0, 0, 0, 31, 0.1),
       (12, 'Csirke', 2, 2, FALSE, 1, 165, 3.6, 1, 0, 0, 0, 0, 31, 0.1),
       (13, 'Chocolate', 1, 1, FALSE, 1, 546, 31.3, 18.8, 59.4, 56.3, 0, 0, 5.4, 0.1),
       (14, 'Csokoládé', 2, 2, FALSE, 1, 546, 31.3, 18.8, 59.4, 56.3, 0, 0, 5.4, 0.1),
       (15, 'Coca Cola', 1, 1, TRUE, 4, 42, 0, 0, 10.6, 10.6, 0, 0, 0, 0.01),
       (16, 'Coca Cola', 2, 2, TRUE, 4, 42, 0, 0, 10.6, 10.6, 0, 0, 0, 0.01),
       (17, 'Coffee', 1, 1, TRUE, 4, 2, 0.1, 0, 0.3, 0, 0.3, 0, 0.1, 0.01),
       (18, 'Kávé', 2, 2, TRUE, 4, 2, 0.1, 0, 0.3, 0, 0.3, 0, 0.1, 0.01),
       (19, 'Cucumber', 1, 1, FALSE, 1, 15, 0.1, 0, 3.6, 1.7, 0, 0, 0.6, 0.002),
       (20, 'Uborka', 2, 2, FALSE, 1, 15, 0.1, 0, 3.6, 1.7, 0, 0, 0.6, 0.002),
       (21, 'Egg', 1, 1, FALSE, 6, 155, 11, 3.3, 1.1, 0.2, 0, 0, 13, 0.5),
       (22, 'Tojás', 2, 2, FALSE, 7, 155, 11, 3.3, 1.1, 0.2, 0, 0, 13, 0.5),
       (23, 'Fish', 1, 1, FALSE, 1, 206, 13.4, 3.3, 0, 0, 0, 0, 20, 0.1),
       (24, 'Hal', 2, 2, FALSE, 1, 206, 13.4, 3.3, 0, 0, 0, 0, 20, 0.1);

-- items_portions_joining_table (only portions which have an owner id should have entries in the joining table; the languages should be matching on the item and the portion)
INSERT INTO items_portions_joining_table(item_id, portion_id)
VALUES (1, 6),
       (2, 7),
       (3, 6),
       (4, 7),
       (21, 6),
       (22, 7);

COMMIT TRANSACTION;
