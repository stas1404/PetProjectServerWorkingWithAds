package adrepo

import (
	"server/internal/ads"
	errors2 "server/internal/errors"
	"server/internal/repo"
	"server/internal/restriction"
	"server/internal/user"
	"sync"
)

// If you want new Repository implementation,
// change this
type SliceRepository struct {
	Ads   *[]ads.Ad
	Users *[]user.User
	Mu    *sync.RWMutex
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

func New() repo.Repository {
	var ads []ads.Ad = make([]ads.Ad, 0)
	var users []user.User = make([]user.User, 0)
	var Mu = sync.RWMutex{}
	return SliceRepository{Ads: &ads, Users: &users, Mu: &Mu}
}
