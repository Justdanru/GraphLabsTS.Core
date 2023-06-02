package models

type GroupListPageContext struct {
	UserId   int64
	Groups   []*Group
	Page     int64
	PrevPage int64
	NextPage int64
}
