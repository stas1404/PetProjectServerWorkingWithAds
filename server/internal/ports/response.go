package ports

import (
	"server/internal/user"
	"time"
)

type ResponseAd struct {
	ID           int64
	Title        string `validate:"min:1, max:100"`
	Text         string `validate:"min:1, max:400"`
	Published    bool   `json:"published"`
	Created      time.Time
	LastModified time.Time
}

type ResponseUser struct {
	ID       int64
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u ResponseUser) Check() error {
	if u.Nickname == "" {
		return user.BadUser("nickname")
	}
	if u.Email == "" {
		return user.BadUser("email")
	}
	if u.Password == "" {
		return user.BadUser("password")
	}
	return nil
}
