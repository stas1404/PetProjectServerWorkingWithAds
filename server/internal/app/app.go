package app

import (
	"context"
	"fmt"
	validator "github.com/stas1404/validator"
	"net/http"
	"server/internal/ads"
	int_errors "server/internal/errors"
	"server/internal/restriction"
	"server/internal/user"
	"time"
)

// Abstruct App, never change it
type App interface {
	GetAd(ctx context.Context, id int64) (ads.Ad, error)
	GetAdsCorresponding(ctx context.Context, res restriction.Restriction) []ads.Ad
	CreateAd(ctx context.Context, Title string, Text string, UserID int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, UserID int64, Published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, UserID int64, Title string, Text string) (ads.Ad, error)
	CreateUser(ctx context.Context, nickname, email, password string, cookie http.Cookie) (user.User, error)
	UpdateUser(ctx context.Context, id int64, nickname, email, password string) (user.User, error)
	GetUserByCookie(ctx context.Context, cookie_name string) (user.User, error)
	GetUserByID(ctx context.Context, id int64) (user.User, error)
	CreateCookie(ctx context.Context, cookie http.Cookie, id int64)
	generateUserID() int64
}

// If you want new App implementation,
// change this
type StructApp struct {
	Ad_amount  int64
	Repository Repository
}

func NewApp(repo Repository) App {
	return StructApp{Ad_amount: 0, Repository: repo}
}

func (app StructApp) GetAd(ctx context.Context, id int64) (ads.Ad, error) {
	if err := ctx.Err(); err != nil {
		return ads.Ad{}, err
	}
	return app.Repository.GetAd(id)
}

func (app StructApp) CreateAd(ctx context.Context, Title string,
	Text string, UserID int64) (ads.Ad, error) {
	if err := ctx.Err(); err != nil {
		return ads.Ad{}, err
	}
	if !app.Repository.ExistUserWithID(UserID) {
		return ads.Ad{}, int_errors.NewErrUnexistingUser(UserID)
	}
	new_ad := ads.Ad{
		ID:           app.Repository.GetAdAmount(),
		Title:        Title,
		Text:         Text,
		AuthorID:     UserID,
		Published:    false,
		Created:      time.Now(),
		LastModified: time.Now(),
	}
	err := validator.Validate(new_ad)
	if err != nil {
		return ads.Ad{}, err
	}
	return new_ad, app.Repository.AddAd(new_ad)
}
func (app StructApp) ChangeAdStatus(ctx context.Context, adID int64,
	UserID int64, Published bool) (ads.Ad, error) {
	if err := ctx.Err(); err != nil {
		return ads.Ad{}, err
	}
	if !app.Repository.ExistUserWithID(UserID) {
		return ads.Ad{}, int_errors.NewErrUnexistingUser(UserID)
	}
	ad, err := app.Repository.GetAd(adID)
	if err != nil {
		return ad, err
	}
	if ad.AuthorID != UserID {
		return ad, int_errors.NewPermissionDenied(UserID)
	}
	ad.ChangeStatus(Published)
	return ad, app.Repository.ChangeAd(adID, ad)
}
func (app StructApp) UpdateAd(ctx context.Context, adID int64,
	UserID int64, Title string, Text string) (ads.Ad, error) {
	if err := ctx.Err(); err != nil {
		return ads.Ad{}, err
	}
	if !app.Repository.ExistUserWithID(UserID) {
		return ads.Ad{}, int_errors.NewErrUnexistingUser(UserID)
	}
	ad, err := app.Repository.GetAd(adID)
	if err != nil {
		return ad, err
	}
	if ad.AuthorID != UserID {
		return ad, int_errors.NewPermissionDenied(UserID)
	}
	ad.ChangeTitleAndText(Title, Text)
	err = validator.Validate(ad)
	if err != nil {
		return ad, err
	}
	return ad, app.Repository.ChangeAd(adID, ad)
}

func (app StructApp) CreateUser(ctx context.Context, nickname, email, password string, cookie http.Cookie) (user.User, error) {
	if ctx.Err() != nil {
		return user.User{}, ctx.Err()
	}
	u := user.User{
		ID:       app.generateUserID(),
		Nickname: nickname,
		Email:    email,
		Password: password,
		Cookie:   cookie.Value,
	}
	err := app.Repository.AddUser(u)
	if err != nil {
		return u, err
	}
	app.Repository.AddCookie(cookie, u.ID)
	return u, err
}

func (app StructApp) UpdateUser(ctx context.Context, id int64, nickname, email, password string) (user.User, error) {
	if ctx.Err() != nil {
		return user.User{}, ctx.Err()
	}
	if !app.Repository.ExistUserWithID(id) {
		return user.User{}, int_errors.NewErrUnexistingUser(id)
	}
	user, err := app.Repository.GetUser(id)
	if err != nil {
		return user, err
	}
	user.ChangeNicknameEmailAndPassword(nickname, email, password)
	err = app.Repository.ChangeUser(id, user)
	return user, err
}

func (app StructApp) GetAdsCorresponding(ctx context.Context, res restriction.Restriction) []ads.Ad {
	if ctx.Err() != nil {
		return []ads.Ad{}
	}
	return app.Repository.GetAdsCorresponding(res)
}

func (app StructApp) GetUserByID(ctx context.Context, id int64) (user.User, error) {
	if ctx.Err() != nil {
		return user.User{}, ctx.Err()
	}
	return app.Repository.GetUser(id)
}

func (app StructApp) GetUserByCookie(ctx context.Context, name string) (u user.User, err error) {
	if err = ctx.Err(); err != nil {
		return
	}
	id, err := app.Repository.GetUserByCookieValue(name)
	fmt.Println("Cookie:", name, ", ID:", id, ", err:", err)
	if err == nil {
		u, err = app.Repository.GetUser(id)
	}
	return
}

func (app StructApp) CreateCookie(ctx context.Context, cookie http.Cookie, id int64) {
	if ctx.Err() != nil {
		return
	}
	app.Repository.AddCookie(cookie, id)
}

func (app StructApp) generateUserID() int64 {
	return app.Repository.GetUserAmount()
}
