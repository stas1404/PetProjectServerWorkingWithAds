package ads

import "time"

type Ad struct {
	ID           int64
	Title        string `validate:"min:1, max:100"`
	Text         string `validate:"min:1, max:400"`
	AuthorID     int64
	Published    bool `json:"published"`
	Created      time.Time
	LastModified time.Time
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

func (ad *Ad) ChangeStatus(Published bool) {
	ad.Published = Published
	ad.ChangeModificationTime()
}

func (ad *Ad) ChangeTitleAndText(title, text string) {
	ad.Title = title
	ad.Text = text
	ad.ChangeModificationTime()
}
