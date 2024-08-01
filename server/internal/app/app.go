package app

import (
	"context"
	"server/internal/adapters/application"
	"server/internal/ads"
	"server/internal/repo"
	"server/internal/restriction"
	"server/internal/user"
)

// Abstruct App, never change it
type App interface {
	GetAd(ctx context.Context, id int64) (ads.Ad, error)
	GetAdsCorresponding(ctx context.Context, res restriction.Restriction) []ads.Ad
	CreateAd(ctx context.Context, Title string, Text string, UserID int64) (ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, UserID int64, Published bool) (ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, UserID int64, Title string, Text string) (ads.Ad, error)
	CreateUser(ctx context.Context, nickname, email, password string) (user.User, error)
	UpdateUser(ctx context.Context, id int64, nickname, email, password string) (user.User, error)
	GetUserByID(ctx context.Context, id int64) (user.User, error)
}

func NewApp(repo repo.Repository) App {
	return application.StructApp{Ad_amount: 0, Repository: repo}
}
