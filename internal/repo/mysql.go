package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/sha3"

	"graphlabsts.core/internal/models"
)

type MySQLRepo struct {
	DB *sql.DB
}

var (
	ErrConnectingDB  = errors.New("error connecting db")
	ErrNoUser        = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
)

// TODO В будущем перенести всё в Docker и считывать переменную окружения
var pepper string = ""

func (r *MySQLRepo) Connect(dsn string) error {
	var err error
	r.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return ErrConnectingDB
	}

	err = r.DB.Ping()
	if err != nil {
		return ErrConnectingDB
	}

	return nil
}

func (r *MySQLRepo) Authorize(login string, password string) (*models.UserAuthData, error) {
	uad := &models.UserAuthData{}
	var salt, userPassword string

	row := r.DB.QueryRow("SELECT id, role, salt, password FROM users WHERE login = ?;", login)
	err := row.Scan(&uad.Id, &uad.RoleCode, &salt, &userPassword)
	if err == sql.ErrNoRows {
		return nil, ErrNoUser
	}

	if userPassword != getPasswordHash(password, salt) {
		fmt.Println(getPasswordHash(password, salt))
		return nil, ErrWrongPassword
	}

	return uad, nil
}

func getPasswordHash(password string, salt string) string {
	return fmt.Sprintf("%x", sha3.Sum256([]byte(password+salt+pepper)))
}
