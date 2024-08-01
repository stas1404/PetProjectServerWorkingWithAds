package user

import (
	"errors"
	errors2 "server/internal/errors"
)

type User struct {
	ID       int64
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) ChangeNicknameEmailAndPassword(nickname, email, password string) {
	u.Nickname = nickname
	u.Email = email
	u.Password = password
}

func AreSame(u, m User) bool {
	return u.ID == m.ID && u.Password == m.Password &&
		u.Email == m.Email && u.Nickname == m.Nickname
}

func BadUser(problem string) errors2.ErrBadUser {
	return errors.New("Can not create User without " + problem)
}
