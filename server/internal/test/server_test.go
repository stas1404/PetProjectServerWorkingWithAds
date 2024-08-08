package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"server/internal/ports"
	"strconv"
	"testing"
	"time"
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
	users := users
	cl := getTestClient()
	for i := range users {
		UserWithID, err := cl.CreateUser(users[i].User.Nickname, users[i].User.Email, users[i].User.Password)
		require.NoError(t, err)
		assert.Equal(t, users[i].User.Email, UserWithID.Email)
		assert.Equal(t, users[i].User.Nickname, UserWithID.Nickname)
		assert.Equal(t, users[i].User.Password, UserWithID.Password)
		users[i].User.ID = UserWithID.ID
	}
	for i := range users {
		auth_user, err := cl.Authorize(users[i].User)
		require.NoError(t, err)
		assert.Equal(t, users[i].User.Email, auth_user.User.Email)
		assert.Equal(t, users[i].User.Nickname, auth_user.User.Nickname)
		assert.Equal(t, users[i].User.Password, auth_user.User.Password)
		assert.NotEmpty(t, auth_user.Cookie)
	}
}

var titles [3]string = [3]string{"first title", "second title", "third title"}
var text [3]string = [3]string{"first text", "second text", "third text"}

func TestCreateAd(t *testing.T) {
	cl := getTestClient()
	users := users[0:3]
	for i := range users {
		UserWithID, err := cl.CreateUser(users[i].User.Nickname, users[i].User.Email, users[i].User.Password)
		require.NoError(t, err)
		auth_user, err := cl.Authorize(UserWithID)
		require.NoError(t, err)
		users[i] = auth_user
	}
	users = append(users, UserTest{
		User:   ports.ResponseUser{},
		Cookie: http.Cookie{},
	})
	for i := range titles {
		for _, u := range users {
			start := time.Now()
			body := map[string]string{
				"title": strconv.FormatInt(u.User.ID, 10) + titles[i],
				"text":  strconv.FormatInt(u.User.ID, 10) + text[i],
			}
			data, err := json.Marshal(body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, cl.baseURL+"/ads/", bytes.NewReader(data))
			require.NoError(t, err)
			req.AddCookie(&u.Cookie)
			req.Header.Add("Content-Type", "application/json")
			resp, err := cl.client.Do(req)
			if u.Cookie.Value == "" {
				assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
				continue
			}
			var ad ports.ResponseAd

			b, err := io.ReadAll(resp.Body)

			require.NoError(t, err)

			err = json.Unmarshal(b, &ad)
			require.NoError(t, err)
			assert.False(t, ad.Published)
			assert.Equal(t, strconv.FormatInt(u.User.ID, 10)+titles[i], ad.Title)
			assert.Equal(t, strconv.FormatInt(u.User.ID, 10)+text[i], ad.Text)
			assert.WithinRange(t, ad.Created, start, time.Now())
		}
	}
}

var modified_titles [3]string = [3]string{"Modified first title", "Modified second title", "Modified third title"}
var modified_text [3]string = [3]string{"Modified first text", "Modified second text", "Modified third text"}

func FuzzModifyAd(f *testing.F) {
	users := users[0:3]
	cl := getTestClient()
	for i := range users {
		UserWithID, err := cl.CreateUser(users[i].User.Nickname, users[i].User.Email, users[i].User.Password)
		if err != nil {
			log.Fatal(users[i], "create user", err)
		}
		auth_user, err := cl.Authorize(UserWithID)
		if err != nil {
			log.Fatal(users[i], "create user", err)
		}
		users[i] = auth_user
	}
	users = append(users, UserTest{
		User:   ports.ResponseUser{},
		Cookie: http.Cookie{},
	})
	author := make([]int64, len(titles)*len(users))
	for i := range titles {
		for _, u := range users {
			ad, _ := cl.CreateAd(u, strconv.FormatInt(u.User.ID, 10)+titles[i], strconv.FormatInt(u.User.ID, 10)+text[i])
			author[ad.ID] = u.User.ID
		}
	}
	f.Fuzz(func(t *testing.T, random_data uint8) {
		user_number := (random_data / 10) % 4
		modif_number := (random_data % 10) % 9
		start := time.Now()
		body := map[string]string{
			"title": strconv.FormatInt(users[user_number].User.ID, 10) + modified_titles[user_number%3],
			"text":  strconv.FormatInt(users[user_number].User.ID, 10) + modified_text[user_number%3],
		}
		data, err := json.Marshal(body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, cl.baseURL+"/ads/"+strconv.FormatUint(uint64(modif_number), 10)+"/edit/", bytes.NewReader(data))
		require.NoError(t, err)
		req.AddCookie(&users[user_number].Cookie)
		req.Header.Add("Content-Type", "application/json")
		resp, err := cl.client.Do(req)
		if user_number == uint8(len(users)-1) {
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			return
		}

		if author[modif_number] != users[user_number].User.ID {
			assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			return
		}
		var ad ports.ResponseAd

		b, err := io.ReadAll(resp.Body)

		require.NoError(t, err)

		err = json.Unmarshal(b, &ad)
		require.NoError(t, err)
		assert.Equal(t, strconv.FormatInt(users[user_number].User.ID, 10)+modified_titles[user_number], ad.Title)
		assert.Equal(t, strconv.FormatInt(users[user_number].User.ID, 10)+modified_text[user_number], ad.Text)
		assert.WithinRange(t, ad.LastModified, start, time.Now())

	})
}
