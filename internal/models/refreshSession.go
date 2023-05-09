package models

type RefreshSession struct {
	RefreshToken string
	Fingerprint  string
	UserId       int64
}
