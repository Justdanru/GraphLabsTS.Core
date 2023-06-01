package models

type ProfilePageContext struct {
	UserId     int64
	RoleString string
	Name       string
	Surname    string
	LastName   string
	TelegramId string
	Groups     []string
	Subjects   []string
	CreatedAt  string
	UpdatedAt  string
}
