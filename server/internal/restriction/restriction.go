package restriction

import (
	"server/internal/ads"
	"slices"
	"time"
)

type Restriction struct {
	Published []bool      `form:"published"`
	Created   []time.Time `form:"created" time_format:"2006-01-02" time_utc:"1"`
	AuthorIDs []int64     `form:"author"`
	Titles    []string    `form:"title"`
}

func CheckCorresponding[E comparable](PRes []E, P E) bool {
	if len(PRes) == 0 {
		return true
	}
	return slices.Contains(PRes, P)
}

func CheckPublishing(Pres []bool, P bool) bool {
	if len(Pres) == 0 {
		return P
	}
	return slices.Contains(Pres, P)
}

func CheckAuthor(IDs []int64, ID int64) bool {
	if len(IDs) == 0 {
		return true
	}
	return slices.Contains(IDs, ID)
}

func Corresponds(ad ads.Ad, res Restriction) bool {
	return CheckPublishing(res.Published, ad.Published) &&
		CheckCorresponding(res.Titles, ad.Title) &&
		CheckCorresponding(res.Created, ad.Created) &&
		CheckAuthor(res.AuthorIDs, ad.AuthorID)

}
