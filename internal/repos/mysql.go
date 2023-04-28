package repos

import (
	"database/sql"
	"fmt"

	"graphlabsts.core/internal/types"
	"graphlabsts.core/pkg/rand"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/sha3"
)

type MySQLRepo struct {
	db *sql.DB
}

// TODO Обернуть ошибку
func (r *MySQLRepo) Connect(dsn string) error {
	var err error
	r.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = r.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// TODO Вынести в конф. файл длину соли
// TODO Обернуть ошибку
// TODO Валидировать данные
// TODO Проверить логин и TelegramID на существование в БД
func (r *MySQLRepo) AddUser(user *types.User) error {
	var nullLogin, nullPassword, nullSalt, nullTelegramID sql.NullString
	var err error
	if user.Login != "" {
		nullLogin.String = user.Login
		nullLogin.Valid = true
	} else {
		nullLogin.Valid = false
	}
	if user.Password != "" {
		nullPassword.String, nullSalt.String, err = r.getPasswordString(user.Password)
		if err != nil {
			return err
		}
		nullPassword.Valid = true
		nullSalt.Valid = true
	} else {
		nullPassword.Valid = false
		nullSalt.Valid = false
	}
	if user.TelegramID != "" {
		nullTelegramID.String = user.TelegramID
		nullTelegramID.Valid = true
	} else {
		nullTelegramID.Valid = false
	}
	_, err = r.db.Exec(
		"INSERT INTO users (name, surname, last_name, login, password, salt, role, tg_id)"+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?)", user.Name, user.Surname, user.LastName,
		nullLogin, nullPassword, nullSalt, user.RoleCode, nullTelegramID,
	)
	if err != nil {
		return err
	}
	return nil
}

// TODO Проверить название предмета на существование в БД
// TODO Проверить creatorId на существование
// TODO Обернуть ошибку
func (r *MySQLRepo) AddSubject(title string, creatorId uint64) error {
	_, err := r.db.Exec(
		"INSERT INTO subjects (title, creator_id) VALUES (?, ?)", title, creatorId,
	)
	if err != nil {
		return err
	}
	return nil
}

// TODO Проверить название группы на существование в БД
// TODO Проверить creatorId на существование
// TODO Обернуть ошибку
func (r *MySQLRepo) AddGroup(name string, creatorId uint64) error {
	_, err := r.db.Exec(
		"INSERT INTO student_groups (name, creator_id) VALUES (?, ?)", name, creatorId,
	)
	if err != nil {
		return err
	}
	return nil
}

// TODO Обернуть ошибку
func (r *MySQLRepo) isSaltExists(salt string) (bool, error) {
	var result string
	row := r.db.QueryRow("SELECT id FROM users WHERE salt = ?;", salt)
	err := row.Scan(&result)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return true, err
	}
	return true, nil
}

func (r *MySQLRepo) getPasswordString(password string) (string, string, error) {
	var salt string
	for {
		salt = rand.String(32)
		f, err := r.isSaltExists(salt)
		if err != nil {
			return "", "", err
		}
		if !f {
			break
		}
	}
	bytes := sha3.Sum256([]byte(password + salt))
	result := fmt.Sprintf("%x", bytes)
	return result, salt, nil
}
