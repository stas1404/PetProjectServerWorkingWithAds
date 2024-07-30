package user

import (
	"errors"
)

type User struct {
	ID       int64
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Cookie   string `json:"cookie"`
}

func (u *User) ChangeNicknameEmailAndPassword(nickname, email, password string) {
	u.Nickname = nickname
	u.Email = email
	u.Password = password
}

func (u *User) Check() error {
	if u.Nickname == "" {
		return BadUser("nickname")
	}
	if u.Email == "" {
		return BadUser("email")
	}
	if u.Password == "" {
		return BadUser("password")
	}
	return nil
}

func AreSame(u, m User) bool {
	return u.ID == m.ID && u.Password == m.Password &&
		u.Email == m.Email && u.Nickname == m.Nickname
}

func BadUser(problem string) error {
	return errors.New("Can not create User without " + problem)
}
