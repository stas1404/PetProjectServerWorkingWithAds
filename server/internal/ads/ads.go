package ads

import (
	"errors"
	"time"
)

type Ad struct {
	ID           int64     `json:"id"`
	Title        string    `validate:"min:1, max:100", json:"title"`
	Text         string    `validate:"min:1, max:400", json:"text"`
	AuthorID     int64     `json:"author_id"`
	Published    bool      `json:"published"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

func New(Title, Text string, AuthorID int64) Ad {
	return Ad{
		ID:        0,
		Title:     Title,
		Text:      Text,
		AuthorID:  AuthorID,
		Published: false,
	}
}

func (ad *Ad) ChangeModificationTime() {
	ad.LastModified = time.Now()
}

func (ad *Ad) ChangeStatus(Published bool) error {
	if ad.Published == Published {
		if ad.Published {
			return errors.New("Ad is already published")
		}
		if !ad.Published {
			return errors.New("Ad is already unpublished")
		}
	}
	ad.Published = Published
	ad.ChangeModificationTime()
	return nil
}

func (ad *Ad) ChangeTitleAndText(title, text string) {
	ad.Title = title
	ad.Text = text
	ad.ChangeModificationTime()
}
