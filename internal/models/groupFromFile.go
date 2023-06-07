package models

type GroupFromFile struct {
	Name     string
	Subjects []string
	Students []*StudentFromFile
}
