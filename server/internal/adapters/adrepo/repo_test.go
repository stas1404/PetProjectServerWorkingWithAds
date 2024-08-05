package adrepo_test

import (
	"github.com/stretchr/testify/assert"
	"server/internal/adapters/adrepo"
	"server/internal/ads"
	repo2 "server/internal/repo"
	"testing"
)

var repo repo2.Repository = adrepo.New()

func TestAddAd(t *testing.T) {
	err := repo.AddAd(ads.New("first title", "first text", 1))
	assert.Equal(t, nil, err, "unexpected error")
}
func TestGetAd(t *testing.T) {
	ad, err := repo.GetAd(0)
	assert.Equal(t, nil, err, "unexpected error")
	assert.Equal(t, "first title", ad.Title, "title was modified by repository")
	assert.Equal(t, "first text", ad.Text, "text was modified by repository")
	assert.Equal(t, int64(1), ad.AuthorID, "author ID was modified by repository")
	ad, err = repo.GetAd(1)
	assert.NotEqual(t, err, nil, "expected an error because try to get unexisting ad")
}

func TestGetAdAmount(t *testing.T) {
	assert.Equal(t, int64(1), repo.GetAdAmount(), "wrong ad amount")
}

func TestChangeAd(t *testing.T) {
	var (
		title          = "second title"
		text           = "second text"
		authorID int64 = 1
	)
	err := repo.ChangeAd(0, ads.New(title, text, authorID))
	assert.NoError(t, err, "unexpected error")
	ad, err := repo.GetAd(0)
	assert.Equal(t, nil, err, "unexpected error")
	assert.Equal(t, title, ad.Title, "title was modified by repository")
	assert.Equal(t, text, ad.Text, "text was modified by repository")
	assert.Equal(t, authorID, ad.AuthorID, "author ID was modified by repository")
	err = repo.ChangeAd(1, ads.New(title, text, authorID))
	assert.NotEqual(t, err, nil, "expected an error because try to modify unexisting ad")
}
