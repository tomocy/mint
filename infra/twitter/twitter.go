package twitter

import (
	"time"

	"github.com/tomocy/mint/domain"
)

type Tweets []*Tweet

func (ts Tweets) Adapt() []*domain.Tweet {
	adapteds := make([]*domain.Tweet, len(ts))
	for i, t := range ts {
		adapteds[i] = t.Adapt()
	}

	return adapteds
}

type Tweet struct {
	ID        string `json:"id_str"`
	User      *User  `json:"user"`
	Text      string `json:"text"`
	FullText  string `json:"full_text"`
	CreatedAt date   `json:"created_at"`
}

func (t *Tweet) Adapt() *domain.Tweet {
	text := t.Text
	if t.FullText != "" {
		text = t.FullText
	}
	return &domain.Tweet{
		ID:        t.ID,
		User:      t.User.Adapt(),
		Text:      text,
		CreatedAt: time.Time(t.CreatedAt),
	}
}

type date time.Time

func (d date) MarshalJSON() ([]byte, error) {
	return []byte((time.Time(d)).Format(time.RubyDate)), nil
}

func (d *date) UnmarshalJSON(data []byte) error {
	withoutQuotes := (string(data))[1 : len(data)-1]
	parsed, err := time.ParseInLocation(time.RubyDate, withoutQuotes, time.UTC)
	if err != nil {
		return err
	}
	*d = date(parsed.In(time.Local))

	return nil
}

type User struct {
	ID         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}

func (u *User) Adapt() *domain.User {
	return &domain.User{
		ID:         u.ID,
		Name:       u.Name,
		ScreenName: u.ScreenName,
	}
}
