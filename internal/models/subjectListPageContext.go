package models

type SubjectListPageContext struct {
	UserId   int64
	Subjects []*Subject
	Page     int64
	PrevPage int64
	NextPage int64
}
