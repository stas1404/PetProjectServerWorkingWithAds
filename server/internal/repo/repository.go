package repo

import (
	"server/internal/ads"
	"server/internal/restriction"
	"server/internal/user"
)

// Abstruct Repository, never change it
type Repository interface {

	// Returns ad if exists. Else returns ErrUnexistingAd
	GetAd(adID int64) (ads.Ad, error)

	// Returns slice of ads which corresponds given restricions
	GetAdsCorresponding(res restriction.Restriction) []ads.Ad

	// Returns amount of stored ads (taking into account unpublished)
	GetAdAmount() int64

	AddAd(ad ads.Ad) error

	// Returns ErrUnexistingAd if ad with this ID does not exist
	ChangeAd(adID int64, ad ads.Ad) error

	// Returns ErrUnexistingUser if User with UserID does not exist
	GetUser(UserID int64) (user.User, error)

	GetUserAmount() int64

	AddUser(user user.User) error

	// Returns ErrUnexistingUser if User with UserID does not exist
	ChangeUser(UserID int64, user user.User) error
	ExistUserWithID(id int64) bool
}
