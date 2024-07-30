package adrepo

import (
	"fmt"
	"net/http"
	"server/internal/ads"
	"server/internal/app"
	errors2 "server/internal/errors"
	"server/internal/errors/cookie_errors"
	"server/internal/restriction"
	"server/internal/user"
	"sync"
	"time"
)

// If you want new Repository implementation,
// change this
type SliceRepository struct {
	Ads    *[]ads.Ad
	Users  *[]user.User
	Cookie *map[string]UserValid
	Mu     *sync.RWMutex
}

type UserValid struct {
	UserID int64
	Valid  time.Time
}

func (r SliceRepository) ExistAdWithID(id int64) bool {
	if id >= r.GetAdAmount() {
		return false
	}
	return true
}

func (r SliceRepository) ExistUserWithID(id int64) bool {
	if id >= r.GetUserAmount() {
		return false
	}
	return true
}

func (r SliceRepository) GetAdAmount() int64 {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return int64(len((*r.Ads)))
}

func (r SliceRepository) GetUserAmount() int64 {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return int64(len(*r.Users))
}

func (r SliceRepository) GetAd(id int64) (ads.Ad, error) {
	if !r.ExistAdWithID(id) {
		return ads.Ad{}, errors2.NewErrUnexistingAd(id)
	}
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return (*r.Ads)[id], nil
}

func (r SliceRepository) GetAdsCorresponding(res restriction.Restriction) []ads.Ad {
	ads := make([]ads.Ad, 0)
	for _, ad := range *r.Ads {
		if restriction.Corresponds(ad, res) {
			ads = append(ads, ad)
		}
	}
	return ads
}

func (r SliceRepository) GetUser(id int64) (user.User, error) {
	if !r.ExistUserWithID(id) {
		return user.User{}, errors2.NewErrUnexistingUser(id)
	}
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	return (*r.Users)[id], nil
}

func (r SliceRepository) AddAd(ad ads.Ad) error {
	r.Mu.Lock()
	*r.Ads = append(*r.Ads, ad)
	r.Mu.Unlock()
	return nil
}

func (r SliceRepository) AddUser(user user.User) error {
	r.Mu.Lock()
	*r.Users = append(*r.Users, user)
	r.Mu.Unlock()
	return nil
}

func (r SliceRepository) ChangeAd(adID int64, ad ads.Ad) error {
	if !r.ExistAdWithID(adID) {
		return errors2.NewErrUnexistingAd(adID)
	}
	r.Mu.Lock()
	(*r.Ads)[adID] = ad
	r.Mu.Unlock()
	return nil
}

func (r SliceRepository) ChangeUser(UserID int64, user user.User) error {
	if !r.ExistUserWithID(UserID) {
		return errors2.NewErrUnexistingUser(UserID)
	}
	r.Mu.Lock()
	(*r.Users)[UserID] = user
	r.Mu.Unlock()
	return nil
}

func (r SliceRepository) GetUserByCookieValue(CookieValue string) (int64, error) {
	r.Mu.RLock()
	u, ok := (*r.Cookie)[CookieValue]
	r.Mu.RUnlock()
	if !ok {
		return u.UserID, cookie_errors.NewErrUnexistingCookie(CookieValue)
	}
	if u.Valid.Before(time.Now()) {
		r.DeleteCookie(CookieValue)
		return 0, cookie_errors.NewExpiredCookie()
	}
	return u.UserID, nil
}

func (r SliceRepository) DeleteCookie(name string) {
	delete(*r.Cookie, name)
}

func (r SliceRepository) AddCookie(cookie http.Cookie, UserID int64) {
	r.Mu.Lock()
	(*r.Cookie)[cookie.Value] = UserValid{
		UserID: UserID,
		Valid:  cookie.Expires,
	}
	r.Mu.Unlock()
	fmt.Println(r.Cookie, r.Users, r.Ads)
}

func New() app.Repository {
	var ads []ads.Ad = make([]ads.Ad, 0)
	var users []user.User = make([]user.User, 0)
	var cookie map[string]UserValid = make(map[string]UserValid, 0)
	var Mu = sync.RWMutex{}
	return SliceRepository{Ads: &ads, Users: &users, Mu: &Mu, Cookie: &cookie}
}
