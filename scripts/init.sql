-- DROP TABLE IF EXISTS ans_opt_to_questions;
DROP TABLE IF EXISTS terms_to_questions;
DROP TABLE IF EXISTS tests_to_subjects;
DROP TABLE IF EXISTS teachers_to_subjects;
DROP TABLE IF EXISTS teachers_to_groups;
DROP TABLE IF EXISTS students_to_groups;
DROP TABLE IF EXISTS groups_to_subjects;
DROP TABLE IF EXISTS terms_to_answers;
DROP TABLE IF EXISTS terms_to_results;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS test_results;
-- DROP TABLE IF EXISTS answer_options;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS terms;
DROP TABLE IF EXISTS tests;
DROP TABLE IF EXISTS student_groups;
DROP TABLE IF EXISTS subjects;
DROP TABLE IF EXISTS users;

-- 2 - student
-- 1 - teacher
-- 0 - admin
CREATE TABLE users (
    id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    surname VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) DEFAULT NULL,
    login VARCHAR(20) UNIQUE,
    password VARCHAR(64),
    salt VARCHAR(32) UNIQUE,
    role TINYINT NOT NULL DEFAULT 2, 
    tg_id VARCHAR(13) UNIQUE DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE student_groups (
    id TINYINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(15) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    creator_id SMALLINT UNSIGNED,
    PRIMARY KEY (id),
    FOREIGN KEY (creator_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE students_to_groups (
    student_id SMALLINT UNSIGNED UNIQUE,
    group_id TINYINT UNSIGNED,
    PRIMARY KEY (student_id, group_id),
    FOREIGN KEY (student_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (group_id) REFERENCES student_groups (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE teachers_to_groups (
    teacher_id SMALLINT UNSIGNED UNIQUE,
    group_id TINYINT UNSIGNED,
    PRIMARY KEY (teacher_id, group_id),
    FOREIGN KEY (teacher_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (group_id) REFERENCES student_groups (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE subjects (
    id TINYINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    creator_id SMALLINT UNSIGNED,
    PRIMARY KEY (id),
    FOREIGN KEY (creator_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT 
);

CREATE TABLE teachers_to_subjects (
    teacher_id SMALLINT UNSIGNED,
    subject_id TINYINT UNSIGNED,
    PRIMARY KEY (teacher_id, subject_id),
    FOREIGN KEY (teacher_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (subject_id) REFERENCES subjects (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE groups_to_subjects (
    group_id TINYINT UNSIGNED,
    subject_id TINYINT UNSIGNED,
    PRIMARY KEY (group_id, subject_id),
    FOREIGN KEY (group_id) REFERENCES student_groups (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (subject_id) REFERENCES subjects (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE tests (
    id SMALLINT UNSIGNED NOT NULL AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    opened_at TIMESTAMP NOT NULL,
    closed_at TIMESTAMP,
    questions_count TINYINT UNSIGNED,
    is_multiple BOOLEAN NOT NULL DEFAULT FALSE,
    is_adaptive BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    creator_id SMALLINT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (creator_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE tests_to_subjects (
    test_id SMALLINT UNSIGNED,
    subject_id TINYINT UNSIGNED,
    PRIMARY KEY (test_id, subject_id),
    FOREIGN KEY (test_id) REFERENCES tests (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (subject_id) REFERENCES subjects (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE terms (
    id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
    value VARCHAR (100) NOT NULL,
    test_id SMALLINT UNSIGNED,
    PRIMARY KEY (id),
    FOREIGN KEY (test_id) REFERENCES tests (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE questions (
    id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
    test_id SMALLINT UNSIGNED,
    type_id TINYINT UNSIGNED NOT NULL,
    content TEXT NOT NULL,
    answer_options TEXT NOT NULL,
    answer TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (test_id) REFERENCES tests (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

-- CREATE TABLE answer_options (
--     id MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT,
--     question_id MEDIUMINT UNSIGNED,
--     value TEXT NOT NULL,
--     PRIMARY KEY (id),
--     FOREIGN KEY (question_id) REFERENCES questions (id)
--         ON DELETE RESTRICT ON UPDATE RESTRICT
-- );

-- CREATE TABLE ans_opt_to_questions (
--     question_id MEDIUMINT UNSIGNED,
--     ansopt_id MEDIUMINT UNSIGNED,
--     PRIMARY KEY (question_id, ansopt_id),
--     FOREIGN KEY (question_id) REFERENCES questions (id)
--         ON DELETE RESTRICT ON UPDATE RESTRICT,
--     FOREIGN KEY (ansopt_id) REFERENCES answer_options (id)
--         ON DELETE RESTRICT ON UPDATE RESTRICT
-- );

CREATE TABLE terms_to_questions (
    question_id MEDIUMINT UNSIGNED,
    term_id MEDIUMINT UNSIGNED,
    PRIMARY KEY (question_id, term_id),
    FOREIGN KEY (question_id) REFERENCES questions (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (term_id) REFERENCES terms (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE test_results (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    test_id SMALLINT UNSIGNED,
    student_id SMALLINT UNSIGNED,
    result TINYINT UNSIGNED NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (test_id) REFERENCES tests (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (student_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE terms_to_results (
    result_id INT UNSIGNED,
    term_id MEDIUMINT UNSIGNED,
    count SMALLINT UNSIGNED NOT NULL,
    PRIMARY KEY (result_id, term_id),
    FOREIGN KEY (result_id) REFERENCES test_results (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (term_id) REFERENCES terms (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE answers (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    question_id MEDIUMINT UNSIGNED,
    student_id SMALLINT UNSIGNED,
    value TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    FOREIGN KEY (question_id) REFERENCES questions (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (student_id) REFERENCES users (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);

CREATE TABLE terms_to_answers (
    answer_id BIGINT UNSIGNED,
    term_id MEDIUMINT UNSIGNED,
    is_correct BOOLEAN NOT NULL,
    PRIMARY KEY (answer_id, term_id),
    FOREIGN KEY (answer_id) REFERENCES answers (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT,
    FOREIGN KEY (term_id) REFERENCES terms (id)
        ON DELETE RESTRICT ON UPDATE RESTRICT
);