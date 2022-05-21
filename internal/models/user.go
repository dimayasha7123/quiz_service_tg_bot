package models

type User struct {
	ID       int64
	TGID     int64
	QSID     int64
	Username string
}

func NewUser(chatID int64, username string) *User {
	u := &User{
		TGID:     chatID,
		Username: username,
	}
	return u
}
