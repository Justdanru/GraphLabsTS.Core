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
		Tmpl:                       templates,
		Repo:                       &repo.MySQLRepo{},
		UncheckAuthMiddlewarePaths: []string{"/", "/login", "/api/auth"},
	}
	err := handlers.Repo.Connect("root:2808@tcp(mysql:3306)/graphlabs_ts?&charset=utf8&interpolateParams=true")
	if err != nil {
		fmt.Printf("Ошибка при подключении к БД.\n")
		return
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	})

	router.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	router.HandleFunc("/api/auth", handlers.Authenticate).Methods("POST")
	router.HandleFunc("/users/{user_id}", handlers.ProfilePage).Methods("GET")

	router.Use(handlers.Authorize)

	err = http.ListenAndServe(":8080", router)
	panic(err)
}
