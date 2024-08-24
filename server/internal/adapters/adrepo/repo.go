package adrepo

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"os"
	"server/internal/ads"
	"server/internal/repo"
	"server/internal/restriction"
	"server/internal/user"
	"time"
)

// If you want new Repository implementation,
// change this
type PostgresRepository struct {
	pool *pgxpool.Pool
}

const QueryTimeout = time.Second

func ConnectionString() string {
	return "postgres://" + os.Getenv("POSTGRES_USER") + ":" + os.Getenv("POSTGRES_PASSWORD") + "@172.19.0.2:5432/" + os.Getenv("POSTGRES_DB")
}

func SetQueryVariable(q *string, path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	*q = string(b)
}

func SetQueries() {
	SetQueryVariable(&GetAds, "internal/adapters/adrepo/queries/ad_queries/get_ad.sql")
	SetQueryVariable(&GetAdAmount, "internal/adapters/adrepo/queries/ad_queries/get_ad_amount.sql")
	SetQueryVariable(&AddAd, "internal/adapters/adrepo/queries/ad_queries/add_ad.sql")
	SetQueryVariable(&ChangeAd, "internal/adapters/adrepo/queries/ad_queries/change_ad.sql")
	SetQueryVariable(&GetAdsCorresponding, "internal/adapters/adrepo/queries/ad_queries/get_ads_corresponding.sql")

	SetQueryVariable(&GetUser, "internal/adapters/adrepo/queries/user_queries/get_user.sql")
	SetQueryVariable(&GetUserAmount, "internal/adapters/adrepo/queries/user_queries/get_users_amount.sql")
	SetQueryVariable(&AddUser, "internal/adapters/adrepo/queries/user_queries/add_user.sql")
	SetQueryVariable(&CountUser, "internal/adapters/adrepo/queries/user_queries/count_user.sql")
}

func New() repo.Repository {
	fmt.Println(ConnectionString())
	background := context.Background()
	m, err := migrate.New("file://internal/adapters/adrepo/migrations",
		ConnectionString()+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("UP")
		log.Fatal(err)
	}
	SetQueries()
	pool, err := pgxpool.New(background, ConnectionString())
	if err != nil {
		log.Fatal(err)
	}
	return PostgresRepository{pool: pool}
}

func (pr PostgresRepository) GetAd(ctx context.Context, adID int64) (ads.Ad, error) {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	r := pr.pool.QueryRow(timedContext, GetAds, adID)
	var ad ads.Ad
	err := r.Scan(&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID, &ad.Published, &ad.Created, &ad.LastModified)
	return ad, err
}

// Returns slice of ads which corresponds given restricions
func (pr PostgresRepository) GetAdsCorresponding(ctx context.Context, res restriction.Restriction) []ads.Ad {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	r, err := pr.pool.Query(timedContext, GetAdsCorresponding, res.Published, res.Created, res.AuthorIDs, res.Titles)
	var ad_list []ads.Ad
	if err != nil {
		return ad_list
	}
	for r.Next() {
		var ad ads.Ad
		err = r.Scan(&ad.ID, &ad.Title, &ad.Text, &ad.AuthorID, &ad.Published, &ad.Created, &ad.LastModified)
		if err != nil {
			return ad_list
		}
		ad_list = append(ad_list, ad)
	}
	return ad_list
}

// Returns amount of stored ads (taking into account unpublished)
func (pr PostgresRepository) GetAdAmount(ctx context.Context) int64 {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	r := pr.pool.QueryRow(timedContext, GetAdAmount)
	var amount int64
	r.Scan(&amount)
	return amount
}

func (pr PostgresRepository) AddAd(ctx context.Context, ad ads.Ad) error {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := pr.pool.Exec(timedContext, AddAd, ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published, ad.Created, ad.LastModified)
	return err
}

// Returns ErrUnexistingAd if ad with this ID does not exist
func (pr PostgresRepository) ChangeAd(ctx context.Context, ad ads.Ad) error {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	fmt.Println(ad)
	_, err := pr.pool.Exec(timedContext, ChangeAd, ad.ID, ad.Title, ad.Text, ad.AuthorID, ad.Published, ad.Created, ad.LastModified)
	return err
}

// Returns ErrUnexistingUser if User with UserID does not exist
func (pr PostgresRepository) GetUser(ctx context.Context, UserID int64) (user.User, error) {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	r := pr.pool.QueryRow(timedContext, GetUser, UserID)
	var u user.User
	fmt.Println(UserID)
	err := r.Scan(&u.ID, &u.Nickname, &u.Email, &u.Password)
	fmt.Println(err)
	return u, err
}

func (pr PostgresRepository) GetUserAmount(ctx context.Context) int64 {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	r := pr.pool.QueryRow(timedContext, GetUserAmount)
	var amount int64
	r.Scan(&amount)
	return amount
}

func (pr PostgresRepository) AddUser(ctx context.Context, user user.User) error {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := pr.pool.Exec(timedContext, AddUser, user.ID, user.Nickname, user.Email, user.Password)
	return err

}

// Returns ErrUnexistingUser if User with UserID does not exist
func (pr PostgresRepository) ChangeUser(ctx context.Context, user user.User) error {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	_, err := pr.pool.Exec(timedContext, AddUser, user.ID, user.Nickname, user.Email, user.Password)
	return err
}
func (pr PostgresRepository) ExistUserWithID(ctx context.Context, id int64) bool {
	timedContext, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()
	r := pr.pool.QueryRow(timedContext, CountUser, id)
	var amount int
	r.Scan(&amount)
	if amount == 1 {
		return true
	}
	return false
}
