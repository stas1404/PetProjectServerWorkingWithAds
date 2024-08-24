package repo

import (
	"context"
	"server/internal/ads"
	"server/internal/restriction"
	"server/internal/user"
)

// Abstruct Repository, never change it
type Repository interface {

	// Returns ad if exists. Else returns ErrUnexistingAd
	GetAd(ctx context.Context, adID int64) (ads.Ad, error)

	// Returns slice of ads which corresponds given restricions
	GetAdsCorresponding(ctx context.Context, res restriction.Restriction) []ads.Ad

	// Returns amount of stored ads (taking into account unpublished)
	GetAdAmount(ctx context.Context) int64

	AddAd(ctx context.Context, ad ads.Ad) error

	// Returns ErrUnexistingAd if ad with this ID does not exist
	ChangeAd(ctx context.Context, ad ads.Ad) error

	// Returns ErrUnexistingUser if User with UserID does not exist
	GetUser(ctx context.Context, UserID int64) (user.User, error)

	GetUserAmount(ctx context.Context) int64

	AddUser(ctx context.Context, user user.User) error

	// Returns ErrUnexistingUser if User with UserID does not exist
	ChangeUser(ctx context.Context, user user.User) error
	ExistUserWithID(ctx context.Context, id int64) bool
}
