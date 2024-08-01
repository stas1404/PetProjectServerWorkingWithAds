package cookie

import (
	"net/http"
	"server/internal/errors/cookie_errors"
	"sync"
	"time"
)

type CookieRepository interface {
	GetUserIDByCookieValue(CookieName string) (int64, error)
	AddCookie(cookie http.Cookie, UserID int64)
	DeleteCookie(name string)
}

type CookieMapRepository struct {
	cookies *sync.Map
}

type UserValid struct {
	UserID int64
	Valid  time.Time
}

func (m CookieMapRepository) GetUserIDByCookieValue(CookieValue string) (int64, error) {
	u, ok := m.cookies.Load(CookieValue)
	if !ok {
		return 0, cookie_errors.NewErrUnexistingCookie(CookieValue)
	}
	uv := u.(UserValid)
	if uv.Valid.Before(time.Now()) {
		m.DeleteCookie(CookieValue)
		return 0, cookie_errors.NewExpiredCookie()
	}
	return uv.UserID, nil
}

func (m CookieMapRepository) DeleteCookie(name string) {
	m.cookies.Delete(name)
}

func (m CookieMapRepository) AddCookie(cookie http.Cookie, UserID int64) {
	m.cookies.Store(cookie.Value, UserValid{
		UserID: UserID,
		Valid:  cookie.Expires,
	})
}

func NewRepository() CookieRepository {
	return CookieMapRepository{
		cookies: &sync.Map{},
	}
}
