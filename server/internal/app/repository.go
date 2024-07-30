package app

import (
	"net/http"
	"server/internal/ads"
	"server/internal/restriction"
	"server/internal/user"
)

// Abstruct Repository, never change it
type Repository interface {
	GetAd(adID int64) (ads.Ad, error)
	GetAdsCorresponding(res restriction.Restriction) []ads.Ad
	GetAdAmount() int64
	AddAd(ad ads.Ad) error
	ChangeAd(adID int64, ad ads.Ad) error
	GetUser(UserID int64) (user.User, error)
	GetUserAmount() int64
	AddUser(user user.User) error
	ChangeUser(UserID int64, user user.User) error
	ExistUserWithID(id int64) bool
	GetUserByCookieValue(CookieName string) (int64, error)
	AddCookie(cookie http.Cookie, UserID int64)
}
