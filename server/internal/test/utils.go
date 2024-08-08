package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"server/internal/adapters/adrepo"
	"server/internal/app"
	"server/internal/ports"
	"server/internal/ports/httpgin"
)

type UserTest struct {
	User   ports.ResponseUser
	Cookie http.Cookie
}

type testClient struct {
	client  *http.Client
	baseURL string
}

type adData struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type adsResponse struct {
	Data []adData `json:"data"`
}

func getTestClient() *testClient {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New()))
	testServer := httptest.NewServer(server.Handler())

	return &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (c *testClient) CreateUser(nickname, email, password string) (ports.ResponseUser, error) {
	body := map[string]string{
		"nickname": nickname,
		"email":    email,
		"password": password,
	}
	data, err := json.Marshal(body)

	if err != nil {
		return ports.ResponseUser{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/users/", bytes.NewReader(data))

	if err != nil {
		return ports.ResponseUser{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return ports.ResponseUser{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return ports.ResponseUser{}, errors.New("invalid status")
	}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return ports.ResponseUser{}, err
	}

	var u ports.ResponseUser
	err = json.Unmarshal(b, &u)
	return u, err
}

func (c *testClient) Authorize(user ports.ResponseUser) (UserTest, error) {
	data, err := json.Marshal(user)

	if err != nil {
		return UserTest{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/users/authorization/", bytes.NewReader(data))

	if err != nil {
		return UserTest{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return UserTest{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return UserTest{}, errors.New("invalid status")
	}

	var u UserTest = UserTest{
		User:   user,
		Cookie: *resp.Cookies()[0],
	}
	return u, err
}

func (c *testClient) CreateAd(user UserTest, title, text string) (ports.ResponseAd, error) {
	body := map[string]string{
		"title": title,
		"text":  text,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return ports.ResponseAd{}, err
	}
	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/ads/", bytes.NewReader(data))
	if err != nil {
		return ports.ResponseAd{}, err
	}
	req.AddCookie(&user.Cookie)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)

	if err != nil {
		return ports.ResponseAd{}, err
	}

	if resp.StatusCode != http.StatusCreated {
		return ports.ResponseAd{}, errors.New("invalid status")
	}

	var ans ports.ResponseAd

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return ports.ResponseAd{}, err
	}

	err = json.Unmarshal(b, &ans)

	return ans, err
}
