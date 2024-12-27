BEGIN TRANSACTION;

-- languages
INSERT INTO languages(name, native_name)
VALUES ('en', 'English'),
       ('hu', 'magyar');

-- portions
INSERT INTO portions(name, owner_id, language_id, liquid, weigth_in_grams, volume_in_ml)
VALUES ('g', NULL, NULL, FALSE, 1, NULL),
       ('ml', NULL, NULL, TRUE, NULL, 1),
       ('dkg', NULL, NULL, FALSE, 10, NULL),
       ('cl', NULL, NULL, TRUE, NULL, 10),
       ('dl', NULL, NULL, TRUE, NULL, 100);

-- items
-- TODO

COMMIT TRANSACTION;
