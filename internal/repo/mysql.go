package repo

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
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

func getPasswordHash(password string, salt string) string {
	return fmt.Sprintf("%x", sha3.Sum256([]byte(password+salt+pepper)))
}
