package models

type User struct {
	Id         int64
	RoleCode   int64
	Name       string
	Surname    string
	LastName   string
	TelegramId string
	Groups     []string
	Subjects   []string
	CreatedAt  string
	UpdatedAt  string
}
