package cookie

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"time"
)

const PreffixKey = "ad_server:cookie:"

type CookieRepository interface {
	GetUserIDByCookieValue(ctx context.Context, CookieName string) (int64, error)
	AddCookie(ctx context.Context, cookie http.Cookie, UserID int64)
	//DeleteCookie(ctx context.Context, name string)
}

type CookieRedisRepository struct {
	cr *redis.Client
}

func (m CookieRedisRepository) GetUserIDByCookieValue(ctx context.Context, CookieValue string) (int64, error) {
	u_id, err := m.cr.Get(ctx, PreffixKey+CookieValue).Result()
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(u_id, 10, 64)

}

func (m CookieRedisRepository) AddCookie(ctx context.Context, cookie http.Cookie, UserID int64) {
	m.cr.Set(ctx, PreffixKey+cookie.Value, UserID, time.Until(cookie.Expires))
}

func NewRepository() CookieRepository {
	log.Println("Start Redis")
	return CookieRedisRepository{
		cr: redis.NewClient(&redis.Options{
			Addr:         "redis:6379",
			Password:     "",
			DB:           0,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}),
	}
}
