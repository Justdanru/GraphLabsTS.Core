package repo

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/sha3"

	"graphlabsts.core/internal/models"
)

type MySQLRepo struct {
	DB *sql.DB
}

// TODO В будущем перенести всё в Docker и считывать переменную окружения
var pepper string = ""

func (r *MySQLRepo) Connect(dsn string) error {
	var err error
	r.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return ErrConnectingDB
	}

	// Если раскоментить, то не будет работать в docker-compose! Непонятно! Разобраться!
	//err = r.DB.Ping()
	//if err != nil {
	//	fmt.Println("Ping")
	//	fmt.Println(err.Error())
	//	return ErrConnectingDB
	//}

	return nil
}

func (r *MySQLRepo) Authenticate(userCredentials *models.UserCredentials) (*models.UserAuthData, error) {
	uad := &models.UserAuthData{}
	var salt, userPassword string

	row := r.DB.QueryRow("SELECT id, role, salt, password FROM users WHERE login = ?;", userCredentials.Login)
	err := row.Scan(&uad.Id, &uad.RoleCode, &salt, &userPassword)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}

	if userPassword != getPasswordHash(userCredentials.Password, salt) {
		return nil, ErrWrongPassword
	}

	return uad, nil
}

func (r *MySQLRepo) AddRefreshSession(rs *models.RefreshSession) error {
	_, err := r.DB.Exec(
		"INSERT INTO refresh_sessions (refresh_token, fingerprint, user_id) VALUES (?, ?, ?);",
		rs.RefreshToken, rs.Fingerprint, rs.UserId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepo) GetRefreshSessionsCountByUserId(userId int64) (int, error) {
	var resultStr string
	row := r.DB.QueryRow("SELECT COUNT(id) FROM refresh_sessions WHERE user_id = ?;", userId)
	err := row.Scan(&resultStr)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(resultStr)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *MySQLRepo) GetRefreshSessionByToken(token string) (*models.RefreshSession, error) {
	rs := &models.RefreshSession{}

	row := r.DB.QueryRow("SELECT refresh_token, fingerprint, user_id FROM refresh_sessions WHERE refresh_token = ?", token)
	err := row.Scan(&rs.RefreshToken, &rs.Fingerprint, &rs.UserId)
	if err == sql.ErrNoRows {
		return nil, ErrNoRefreshSessions
	}

	return rs, nil
}

func (r *MySQLRepo) UpdateRefreshSession(rs *models.RefreshSession, oldRefreshToken string) error {
	_, err := r.DB.Exec(
		"UPDATE refresh_sessions SET refresh_token = ?, fingerprint = ?, created_at = NOW() WHERE refresh_token = ?;",
		rs.RefreshToken, rs.Fingerprint, oldRefreshToken,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepo) DeleteAllRefreshSessionsByUserId(userId int64) error {
	_, err := r.DB.Exec(
		"DELETE FROM refresh_sessions WHERE user_id = ?;", userId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MySQLRepo) GetUser(userId int64) (*models.User, error) {
	user := &models.User{}

	var telegramId, lastName, createdAt, updatedAt sql.NullString
	row := r.DB.QueryRow("SELECT name, surname, last_name, tg_id, created_at, updated_at FROM users WHERE id = ?", userId)
	err := row.Scan(&user.Name, &user.Surname, &lastName, &telegramId, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	if lastName.Valid {
		user.LastName = lastName.String
	} else {
		user.LastName = ""
	}
	if telegramId.Valid {
		user.TelegramId = telegramId.String
	} else {
		user.TelegramId = ""
	}
	if createdAt.Valid {
		user.CreatedAt = createdAt.String
	} else {
		user.CreatedAt = ""
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.String
	} else {
		user.UpdatedAt = ""
	}

	subjectRows, err := r.DB.Query(`SELECT subjects.title FROM subjects
	INNER JOIN teachers_to_subjects ON teachers_to_subjects.subject_id = subjects.id
	INNER JOIN users ON teachers_to_subjects.teacher_id = users.id
	WHERE users.id = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer subjectRows.Close()

	for subjectRows.Next() {
		var subject string
		err = subjectRows.Scan(&subject)
		if err != nil {
			return nil, err
		}
		user.Subjects = append(user.Subjects, subject)
	}

	groupsRows, err := r.DB.Query(`SELECT student_groups.name FROM student_groups
	INNER JOIN teachers_to_groups ON teachers_to_groups.group_id = student_groups.id
	INNER JOIN users ON teachers_to_groups.teacher_id = users.id
	WHERE users.id = ?`, userId)
	if err != nil {
		return nil, err
	}
	defer groupsRows.Close()

	for groupsRows.Next() {
		var group string
		err = groupsRows.Scan(&group)
		if err != nil {
			return nil, err
		}
		user.Groups = append(user.Groups, group)
	}

	return user, nil
}

func getPasswordHash(password string, salt string) string {
	return fmt.Sprintf("%x", sha3.Sum256([]byte(password+salt+pepper)))
}
