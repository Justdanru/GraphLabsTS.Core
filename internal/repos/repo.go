package repos

import "graphlabsts.core/internal/types"

type Repo interface {
	Connect(dsn string) error
	AddUser(user *types.User) error
	AddSubject(title string, creatorId uint64) error
	AddGroup(name string, creatorId uint64) error
}
