package main

import (
	"fmt"
	"html/template"
	"net/http"

	"graphlabsts.core/internal/handlers"
	"graphlabsts.core/internal/repo"

	"github.com/gorilla/mux"
)

func main() {
	templates := template.Must(template.ParseGlob("./templates/*.html"))

	handlers := &handlers.Handler{
		Tmpl: templates,
		Repo: &repo.MySQLRepo{},
	}
	err := handlers.Repo.Connect("root:2808@tcp(localhost:3306)/graphlabs_ts?&charset=utf8&interpolateParams=true")
	if err != nil {
		fmt.Printf("Ошибка при подключении к БД.\n")
		return
	}

	// TODO Сделать перенаправление на страницу входа или профиль с URL "/"
	r := mux.NewRouter()
	r.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/login", handlers.Login).Methods("POST")
	r.HandleFunc("/profile", handlers.ProfilePage).Methods("GET")

	err = http.ListenAndServe(":8080", r)
	panic(err)
}
