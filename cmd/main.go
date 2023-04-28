package main

import (
	"graphlabsts.core/internal/repos"
	"graphlabsts.core/internal/types"
)

// TODO Вынести в отдельный файл все валидаторы
func main() {
	repo := &repos.MySQLRepo{}
	err := repo.Connect("root:2808@tcp(localhost:3306)/graphlabs_ts?&charset=utf8&interpolateParams=true")
	if err != nil {
		panic(err)
	}
	exampleAdmin := &types.User{
		Name:       "Денис",
		Surname:    "Семёнов",
		LastName:   "Александрович",
		Login:      "admin01",
		Password:   "adminpassword",
		RoleCode:   types.AdminRoleCode,
		TelegramID: "22222222222",
	}
	exampleTeacher := &types.User{
		Name:       "Мария",
		Surname:    "Короткова",
		LastName:   "Александровна",
		Login:      "korotkova01",
		Password:   "123",
		RoleCode:   types.TeacherRoleCode,
		TelegramID: "",
	}
	exampleStudent := &types.User{
		Name:       "Иван",
		Surname:    "Иванов",
		LastName:   "Иванович",
		Login:      "",
		Password:   "",
		RoleCode:   types.StudentRoleCode,
		TelegramID: "11111111111",
	}
	err = repo.AddUser(exampleAdmin)
	if err != nil {
		panic(err)
	}
	err = repo.AddUser(exampleTeacher)
	if err != nil {
		panic(err)
	}
	err = repo.AddUser(exampleStudent)
	if err != nil {
		panic(err)
	}
	err = repo.AddSubject("Дискретная математика", 2)
	if err != nil {
		panic(err)
	}
	err = repo.AddGroup("Б20-555", 1)
	if err != nil {
		panic(err)
	}
}
