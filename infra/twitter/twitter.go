package twitter

import (
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
	ID   string `json:"id_str"`
	User *User  `json:"user"`
	Text string `json:"text"`
}

func (t *Tweet) Adapt() *domain.Tweet {
	return &domain.Tweet{
		ID:   t.ID,
		User: t.User.Adapt(),
		Text: t.Text,
	}
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
