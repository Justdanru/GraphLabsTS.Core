package models

type ProfilePageContext struct {
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
