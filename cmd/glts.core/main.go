package main

import (
	"fmt"
	"html/template"
	"net/http"

	"graphlabsts.core/internal/handlers"
	"graphlabsts.core/internal/middleware"
	"graphlabsts.core/internal/repo"

	"github.com/gorilla/mux"
)

func main() {
	templates := template.Must(template.ParseGlob("./templates/*.html"))

	handlers := &handlers.Handler{
		Tmpl: templates,
		Repo: &repo.MySQLRepo{},
	}
	err := handlers.Repo.Connect("root:2808@tcp(mysql:3306)/graphlabs_ts?&charset=utf8&interpolateParams=true")
	if err != nil {
		fmt.Printf("Ошибка при подключении к БД.\n")
		return
	}

	// TODO Сделать перенаправление на страницу входа или профиль с URL "/"
	router := mux.NewRouter()

	router.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	router.HandleFunc("/api/auth", handlers.Authenticate).Methods("POST")
	router.HandleFunc("/profile", handlers.ProfilePage).Methods("GET")

	authMiddleware := middleware.Middleware{
		UncheckPaths: []string{"/login"},
	}

	router.Use(authMiddleware.Authorization)

	err = http.ListenAndServe(":8080", router)
	panic(err)
}
