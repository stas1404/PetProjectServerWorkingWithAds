package application

import (
	"context"
	"server/internal/repo"
	"time"

	validator "github.com/stas1404/validator"
	"server/internal/ads"

	int_errors "server/internal/errors"
	"server/internal/restriction"
	"server/internal/user"
)

// If you want new App implementation,
// change this
type StructApp struct {
	Ad_amount  int64
	Repository repo.Repository
}

func (app StructApp) GetAd(ctx context.Context, id int64) (ads.Ad, error) {
	return app.Repository.GetAd(ctx, id)
}

func (app StructApp) CreateAd(ctx context.Context, Title string,
	Text string, UserID int64) (ads.Ad, error) {
	if !app.Repository.ExistUserWithID(ctx, UserID) {
		return ads.Ad{}, int_errors.NewErrUnexistingUser(UserID)
	}
	new_ad := ads.Ad{
		ID:           app.Repository.GetAdAmount(ctx),
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
	return new_ad, app.Repository.AddAd(ctx, new_ad)
}
func (app StructApp) ChangeAdStatus(ctx context.Context, adID int64,
	UserID int64, Published bool) (ads.Ad, error) {
	if !app.Repository.ExistUserWithID(ctx, UserID) {
		return ads.Ad{}, int_errors.NewErrUnexistingUser(UserID)
	}
	ad, err := app.Repository.GetAd(ctx, adID)
	if err != nil {
		return ad, err
	}
	if ad.AuthorID != UserID {
		return ad, int_errors.NewPermissionDenied(UserID)
	}
	err = ad.ChangeStatus(Published)
	if err != nil {
		return ad, err
	}
	return ad, app.Repository.ChangeAd(ctx, ad)
}
func (app StructApp) UpdateAd(ctx context.Context, adID int64,
	UserID int64, Title string, Text string) (ads.Ad, error) {
	if !app.Repository.ExistUserWithID(ctx, UserID) {
		return ads.Ad{}, int_errors.NewErrUnexistingUser(UserID)
	}
	ad, err := app.Repository.GetAd(ctx, adID)
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
	return ad, app.Repository.ChangeAd(ctx, ad)
}

func (app StructApp) CreateUser(ctx context.Context, nickname, email, password string) (user.User, error) {
	u := user.User{
		ID:       app.generateUserID(),
		Nickname: nickname,
		Email:    email,
		Password: password,
	}
	err := app.Repository.AddUser(ctx, u)
	return u, err
}

func (app StructApp) UpdateUser(ctx context.Context, id int64, nickname, email, password string) (user.User, error) {
	if !app.Repository.ExistUserWithID(ctx, id) {
		return user.User{}, int_errors.NewErrUnexistingUser(id)
	}
	user, err := app.Repository.GetUser(ctx, id)
	if err != nil {
		return user, err
	}
	user.ChangeNicknameEmailAndPassword(nickname, email, password)
	err = app.Repository.ChangeUser(ctx, user)
	return user, err
}

func (app StructApp) GetAdsCorresponding(ctx context.Context, res restriction.Restriction) []ads.Ad {
	return app.Repository.GetAdsCorresponding(ctx, res)
}

func (app StructApp) GetUserByID(ctx context.Context, id int64) (user.User, error) {
	return app.Repository.GetUser(ctx, id)
}

func (app StructApp) generateUserID() int64 {
	return app.Repository.GetUserAmount(context.Background())
}
