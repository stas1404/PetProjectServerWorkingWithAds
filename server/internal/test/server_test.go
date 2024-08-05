package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"server/internal/ports"
	"testing"
)

var users [3]UserTest = [3]UserTest{
	UserTest{
		User: ports.ResponseUser{
			ID:       0,
			Nickname: "first_nickname",
			Email:    "first_email",
			Password: "first_password",
		},
		Cookie: http.Cookie{},
	},
	UserTest{
		User: ports.ResponseUser{
			ID:       0,
			Nickname: "second_nickname",
			Email:    "second_email",
			Password: "second_password",
		},
		Cookie: http.Cookie{},
	},
	UserTest{
		User: ports.ResponseUser{
			ID:       0,
			Nickname: "third_nickname",
			Email:    "third_email",
			Password: "third_password",
		},
		Cookie: http.Cookie{},
	},
}

func TestCreateUser(t *testing.T) {
	var ResponseCode [6]int = [6]int{
		http.StatusCreated,
		http.StatusCreated,
		http.StatusCreated,
		http.StatusBadRequest,
		http.StatusBadRequest,
		http.StatusBadRequest,
	}
	var create_users = append(users[0:3], UserTest{
		User: ports.ResponseUser{
			ID:       0,
			Nickname: "",
			Email:    "wrongnameemail",
			Password: "wrongnamepassword",
		},
		Cookie: http.Cookie{},
	},
		UserTest{
			User: ports.ResponseUser{
				ID:       0,
				Nickname: "wrongemailname",
				Email:    "",
				Password: "wrongemailpassword",
			},
			Cookie: http.Cookie{},
		},
		UserTest{
			User: ports.ResponseUser{
				ID:       0,
				Nickname: "wrongpasswordname",
				Email:    "wrongpasswordemail",
				Password: "",
			},
			Cookie: http.Cookie{},
		})
	cl := getTestClient()
	for i := range create_users {
		//resp, err := cl.CreateUser(create_users[i].User.Nickname, create_users[i].User.Email, create_users[i].User.Password)
		body := map[string]string{
			"nickname": create_users[i].User.Nickname,
			"email":    create_users[i].User.Email,
			"password": create_users[i].User.Password,
		}
		data, err := json.Marshal(body)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, cl.baseURL+"/users/", bytes.NewReader(data))

		require.NoError(t, err)

		req.Header.Add("Content-Type", "application/json")

		resp, err := cl.client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, ResponseCode[i], resp.StatusCode)
		if ResponseCode[i] != http.StatusCreated {
			continue
		}
		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		var u ports.ResponseUser
		err = json.Unmarshal(b, &u)
		require.NoError(t, err)
		assert.Equal(t, create_users[i].User.Email, u.Email)
		assert.Equal(t, create_users[i].User.Nickname, u.Nickname)
		assert.Equal(t, create_users[i].User.Password, u.Password)
	}
}

func TestAuthorization(t *testing.T) {
	cl := getTestClient()
	for i := range users {
		UserWithID, err := cl.CreateUser(users[i].User.Nickname, users[i].User.Email, users[i].User.Password)
		assert.Equal(t, users[i].User.Email, UserWithID.Email)
		assert.Equal(t, users[i].User.Nickname, UserWithID.Nickname)
		assert.Equal(t, users[i].User.Password, UserWithID.Password)
		users[i].User.ID = UserWithID.ID
		require.NoError(t, err)
		users[i].User = UserWithID
		auth_user, err := cl.Authorize(users[i].User)
		require.NoError(t, err)
		assert.NotEmpty(t, auth_user.Cookie)
	}
}
