SET NAMES utf8;
INSERT INTO users (name, surname, last_name, login, password, salt, role, tg_id, created_at, updated_at, deleted_at) VALUES ("Админист", "Администраторов", "Администраторович", "admin01", "fc93e7998b2d184c101006e800c8a9c6fe74b5d99a428d0364ad41eaed79eece", "@vO5Yb538i(lTESZS%Rs3oq88C)A,hrU", 0, NULL, "2023-03-01 12:00:00", NULL, NULL);

INSERT INTO student_groups (name, created_at, updated_at, deleted_at, creator_id) VALUES
("Б20-504", "2023-03-01 12:20:00", NULL, NULL, 1),
("Б20-505", "2023-03-01 12:21:00", NULL, NULL, 1),
("Б20-501", "2023-03-01 12:22:00", NULL, NULL, 1);

INSERT INTO subjects (title, created_at, updated_at, deleted_at, creator_id) VALUES
("Теория графов", "2023-03-01 13:00:00", NULL, NULL, 1),
("Теория функции комплексного переменного", "2023-03-01 13:00:01", NULL, NULL, 1),
("Математический анализ", "2023-03-01 13:00:02", NULL, NULL, 1);

INSERT INTO teachers_to_subjects (teacher_id, subject_id) VALUES
(1, 1),
(1, 2),
(1, 3);

INSERT INTO teachers_to_groups (teacher_id, group_id) VALUES
(1, 1),
(1, 2),
(1, 3);